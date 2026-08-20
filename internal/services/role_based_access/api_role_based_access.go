// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/tenant"
)

const apiVersion = "2024-10-01"

// The RBAC API rejects a freshly created environment or environment group until its id propagates.
var scopePropagationErrorCodes = []string{"EnvironmentIdInvalid", "EnvironmentGroupIdInvalid"}

func newRoleBasedAccessClient(apiClient *api.Client) client {
	return client{
		Api:       apiClient,
		TenantApi: tenant.NewTenantClient(apiClient),
	}
}

type client struct {
	Api       *api.Client
	TenantApi tenant.Client
}

// url builds an absolute RBAC API url for the given path.
func (client *client) url(path string) string {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   path,
	}
	values := url.Values{}
	values.Add("api-version", apiVersion)
	apiUrl.RawQuery = values.Encode()
	return apiUrl.String()
}

// tenantScope returns the fully qualified tenant scope string required by the RBAC API.
func (client *client) tenantScope(ctx context.Context) (string, error) {
	tenantDto, err := client.TenantApi.GetTenant(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/tenants/%s", tenantDto.TenantId), nil
}

// CreateRoleAssignment creates a role assignment at the given scope.
//
// The POST is not idempotent, so it is never replayed: a retried request that actually committed the
// first time would grant access twice or fail as a duplicate while the real assignment goes
// untracked. The one deliberate retry is the scope-propagation 400, which is safe because a 400
// means nothing was created. Any other failure is reconciled by listing the scope: if an assignment
// matching the principal, role and type is there, the request committed and it is adopted.
func (client *client) CreateRoleAssignment(ctx context.Context, scope assignmentScope, request roleAssignmentRequestDto) (*roleAssignmentDto, error) {
	tenantScope, err := client.tenantScope(ctx)
	if err != nil {
		return nil, err
	}
	request.Scope = scope.qualify(tenantScope)

	apiUrl := client.url(scope.collectionPath())
	response := roleAssignmentDto{}
	for {
		_, err := client.Api.ExecuteWithoutRetry(ctx, nil, "POST", apiUrl, nil, request, []int{http.StatusOK, http.StatusCreated}, &response)
		if err == nil {
			return &response, nil
		}
		if isScopeNotYetPropagated(err) {
			waitFor := api.DefaultRetryAfter()
			tflog.Debug(ctx, fmt.Sprintf("Scope %s is not yet available for role assignments, retrying after %s", request.Scope, waitFor))
			if sleepErr := client.Api.SleepWithContext(ctx, waitFor); sleepErr != nil {
				return nil, err
			}
			continue
		}

		if adopted := client.findExistingAssignment(ctx, scope, request); adopted != nil {
			tflog.Debug(ctx, fmt.Sprintf("Create at scope %s failed ambiguously but the assignment exists; adopting %s", request.Scope, adopted.RoleAssignmentId))
			return adopted, nil
		}
		return nil, err
	}
}

// findExistingAssignment looks for an assignment matching the request's principal, role and type at
// the given scope. It is used to reconcile an ambiguous create failure, so its own errors are
// swallowed: the caller returns the original failure.
func (client *client) findExistingAssignment(ctx context.Context, scope assignmentScope, request roleAssignmentRequestDto) *roleAssignmentDto {
	assignments, err := client.ListRoleAssignments(ctx, scope)
	if err != nil {
		return nil
	}
	for i := range assignments {
		if assignments[i].PrincipalObjectId == request.PrincipalObjectId &&
			assignments[i].RoleDefinitionId == request.RoleDefinitionId &&
			assignments[i].PrincipalType == request.PrincipalType {
			return &assignments[i]
		}
	}
	return nil
}

func isScopeNotYetPropagated(err error) bool {
	var httpErr customerrors.UnexpectedHttpStatusCodeError
	if !errors.As(err, &httpErr) || httpErr.StatusCode != http.StatusBadRequest {
		return false
	}
	for _, code := range scopePropagationErrorCodes {
		if strings.Contains(string(httpErr.Body), code) {
			return true
		}
	}
	return false
}

// ListRoleAssignments lists the role assignments at the given scope.
func (client *client) ListRoleAssignments(ctx context.Context, scope assignmentScope) ([]roleAssignmentDto, error) {
	response := roleAssignmentsListDto{}
	resp, err := client.Api.Execute(ctx, nil, "GET", client.url(scope.collectionPath()), nil, nil, []int{http.StatusOK, http.StatusNotFound}, &response)
	if err != nil {
		return nil, err
	}
	// The API returns 404 when the scope itself (the environment or environment group) no longer
	// exists, which the caller distinguishes from an empty list so it can drop child state.
	if resp.HttpResponse.StatusCode == http.StatusNotFound {
		return nil, customerrors.WrapIntoProviderError(nil, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("scope not found: %s", scope))
	}
	return response.Value, nil
}

// DeleteRoleAssignment deletes a role assignment at the given scope.
func (client *client) DeleteRoleAssignment(ctx context.Context, scope assignmentScope, roleAssignmentId string) error {
	// 404 is accepted: an assignment that is already absent, or whose scope is gone, is destroyed.
	_, err := client.Api.Execute(ctx, nil, "DELETE", client.url(scope.assignmentPath(roleAssignmentId)), nil, nil, []int{http.StatusOK, http.StatusNoContent, http.StatusNotFound}, nil)
	return err
}

// ListRoleDefinitions lists available role definitions.
func (client *client) ListRoleDefinitions(ctx context.Context) ([]roleDefinitionDto, error) {
	response := roleDefinitionsListDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", client.url("/authorization/roleDefinitions"), nil, nil, []int{http.StatusOK}, &response)
	if err != nil {
		return nil, err
	}
	return response.Value, nil
}
