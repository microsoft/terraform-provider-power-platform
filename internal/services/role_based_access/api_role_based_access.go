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
func (client *client) CreateRoleAssignment(ctx context.Context, scope assignmentScope, request roleAssignmentRequestDto) (*roleAssignmentDto, error) {
	tenantScope, err := client.tenantScope(ctx)
	if err != nil {
		return nil, err
	}
	request.Scope = scope.qualify(tenantScope)

	apiUrl := client.url(scope.collectionPath())
	response := roleAssignmentDto{}
	for {
		_, err := client.Api.Execute(ctx, nil, "POST", apiUrl, nil, request, []int{http.StatusOK, http.StatusCreated}, &response)
		if err == nil {
			return &response, nil
		}
		if !isScopeNotYetPropagated(err) {
			return nil, err
		}

		waitFor := api.DefaultRetryAfter()
		tflog.Debug(ctx, fmt.Sprintf("Scope %s is not yet available for role assignments, retrying after %s", request.Scope, waitFor))
		if sleepErr := client.Api.SleepWithContext(ctx, waitFor); sleepErr != nil {
			return nil, err
		}
	}
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
	_, err := client.Api.Execute(ctx, nil, "GET", client.url(scope.collectionPath()), nil, nil, []int{http.StatusOK}, &response)
	if err != nil {
		return nil, err
	}
	return response.Value, nil
}

// DeleteRoleAssignment deletes a role assignment at the given scope.
func (client *client) DeleteRoleAssignment(ctx context.Context, scope assignmentScope, roleAssignmentId string) error {
	_, err := client.Api.Execute(ctx, nil, "DELETE", client.url(scope.assignmentPath(roleAssignmentId)), nil, nil, []int{http.StatusOK, http.StatusNoContent}, nil)
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
