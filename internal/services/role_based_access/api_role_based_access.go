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

// CreateRoleAssignment creates a role assignment at the given scope, with explicit relationship
// semantics. The configuration identifies the relationship (scope, principal, type, role) while the
// assignment id is computed, and the API happily stores duplicates of the same relationship, so:
//
//   - if exactly one matching assignment already exists it is adopted without a POST, and
//     destroying the resource will remove it;
//   - if several already exist the create fails, because adopting any one of them would leave
//     Terraform managing an arbitrary duplicate: deduplicate or import instead;
//   - otherwise the relationship is created with a single POST. The POST is never replayed, since
//     a retry after an ambiguous failure could grant twice. The one exception is the specific
//     scope-propagation 400, which proves nothing was committed. Any other failure that leaves
//     the outcome unknown is returned as an explicit unknown-outcome error: the API assigns the
//     id server-side, permits duplicates, and issues no idempotency or correlation token, so no
//     later listing can prove which assignment, if any, this request created. Guessing could
//     seize a concurrent caller's assignment, so nothing is adopted and nothing reaches state.
func (client *client) CreateRoleAssignment(ctx context.Context, scope assignmentScope, request roleAssignmentRequestDto) (*roleAssignmentDto, error) {
	tenantScope, err := client.tenantScope(ctx)
	if err != nil {
		return nil, err
	}
	request.Scope = scope.qualify(tenantScope)

	// Preflight, authoritatively: adopting or refusing duplicates is only sound against a real
	// listing, so a failed list stops the create rather than degrading to an empty one. A missing
	// child scope is the one retriable case, since a freshly created environment or group can lag.
	var existing []roleAssignmentDto
	for {
		existing, err = client.ListRoleAssignments(ctx, scope)
		if err == nil {
			break
		}
		if scope.kind != scopeTenant && errors.Is(err, customerrors.ErrObjectNotFound) {
			waitFor := api.DefaultRetryAfter()
			tflog.Debug(ctx, fmt.Sprintf("Scope %s is not visible for role assignments yet, retrying preflight after %s", request.Scope, waitFor))
			if sleepErr := client.Api.SleepWithContext(ctx, waitFor); sleepErr != nil {
				return nil, err
			}
			continue
		}
		return nil, fmt.Errorf("could not list the existing role assignments before creating, so it is unknown whether one for this relationship already exists, and creating blindly could duplicate or mismanage it: %w", err)
	}

	matches := matchingAssignments(existing, request)
	if len(matches) == 1 {
		tflog.Debug(ctx, fmt.Sprintf("Adopting existing role assignment %s at scope %s", matches[0].RoleAssignmentId, request.Scope))
		return &matches[0], nil
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("found %d existing role assignments for this principal, role and scope; the API permits duplicates, so deduplicate them or import one before managing it with Terraform", len(matches))
	}

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
		if isAmbiguousCreateFailure(err) {
			return nil, unknownOutcomeError(err)
		}
		// A definitive rejection: the request never committed.
		return nil, err
	}
}

// unknownOutcomeError wraps a create failure whose outcome the API gives no way to resolve.
func unknownOutcomeError(err error) error {
	return fmt.Errorf("the create outcome is unknown. The API may have committed the assignment, but it provides no idempotency key or correlation identifier that would let the provider identify it safely. Inspect the role assignments at this scope: if one was created, import it using the scope-shaped import id, and if duplicates exist, deduplicate them first. Terraform has not recorded an assignment in state. Original failure: %w", err)
}

// isAmbiguousCreateFailure reports whether the request might have committed despite the error. A
// definitive HTTP rejection (400, 401, 403, 404, 409) means the server refused it; a transport
// failure, timeout, throttle or server error leaves the outcome unknown.
func isAmbiguousCreateFailure(err error) bool {
	var httpErr customerrors.UnexpectedHttpStatusCodeError
	if !errors.As(err, &httpErr) {
		// No HTTP status at all: the transport failed with the outcome unknown.
		return true
	}
	return httpErr.StatusCode == http.StatusRequestTimeout ||
		httpErr.StatusCode == http.StatusTooManyRequests ||
		httpErr.StatusCode >= http.StatusInternalServerError
}

// matchingAssignments returns the assignments with the request's principal, role and type.
// Identifiers are compared case-insensitively, since the service renders guids in its own casing.
// Expiring assignments never match: this resource has no expiry input, so it can neither have
// created one nor faithfully manage one.
func matchingAssignments(assignments []roleAssignmentDto, request roleAssignmentRequestDto) []roleAssignmentDto {
	var matches []roleAssignmentDto
	for i := range assignments {
		if assignments[i].ExpiresOn != nil && *assignments[i].ExpiresOn != "" {
			continue
		}
		if strings.EqualFold(assignments[i].PrincipalObjectId, request.PrincipalObjectId) &&
			strings.EqualFold(assignments[i].RoleDefinitionId, request.RoleDefinitionId) &&
			strings.EqualFold(assignments[i].PrincipalType, request.PrincipalType) {
			matches = append(matches, assignments[i])
		}
	}
	return matches
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
	// The API returns 404 when the scope itself no longer exists. That only means something for
	// environment and environment group scopes, whose parents can be deleted; the tenant cannot
	// disappear, so a tenant-scope 404 is a service fault and must stay an error. Otherwise a
	// transient 404 would silently untrack an active tenant-wide grant.
	if resp.HttpResponse.StatusCode == http.StatusNotFound {
		if scope.kind == scopeTenant {
			return nil, errors.New("the tenant role assignment collection returned 404, which cannot mean a deleted scope; treating it as a service error")
		}
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
