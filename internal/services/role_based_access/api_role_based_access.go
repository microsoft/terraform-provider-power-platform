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
//   - otherwise the relationship is created. The POST is never replayed, since a retry after an
//     ambiguous failure could grant twice. Instead the scope is re-listed and an assignment
//     matching the relationship that was NOT in the pre-create baseline is adopted: that one is
//     ours. Adopting the first tuple match would risk adopting a pre-existing duplicate and later
//     destroying it while the one this create made lives on untracked.
func (client *client) CreateRoleAssignment(ctx context.Context, scope assignmentScope, request roleAssignmentRequestDto) (*roleAssignmentDto, error) {
	tenantScope, err := client.tenantScope(ctx)
	if err != nil {
		return nil, err
	}
	request.Scope = scope.qualify(tenantScope)

	// Preflight, authoritatively: the ownership rules below are only sound against a real baseline,
	// so a failed list stops the create rather than degrading to an empty one. A missing child
	// scope is the one retriable case, since a freshly created environment or group can lag here.
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
		return nil, fmt.Errorf("could not list the existing role assignments before creating, and creating without that baseline could duplicate or mismanage an existing assignment: %w", err)
	}

	matches := matchingAssignments(existing, request)
	if len(matches) == 1 {
		tflog.Debug(ctx, fmt.Sprintf("Adopting existing role assignment %s at scope %s", matches[0].RoleAssignmentId, request.Scope))
		return &matches[0], nil
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("found %d existing role assignments for this principal, role and scope; the API permits duplicates, so deduplicate them or import one before managing it with Terraform", len(matches))
	}
	baseline := map[string]struct{}{}
	for i := range existing {
		baseline[strings.ToLower(existing[i].RoleAssignmentId)] = struct{}{}
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
		if !isAmbiguousCreateFailure(err) {
			// A definitive rejection: the request never committed, so there is nothing to reconcile.
			return nil, err
		}
		return client.reconcileAmbiguousCreate(ctx, scope, request, baseline, err)
	}
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

// reconcileAttempts bounds how long an ambiguous create failure is reconciled for. The committed
// assignment can lag the failed response, so one immediate list is not enough, but polling forever
// would hide a genuine failure.
const reconcileAttempts = 4

// reconcileAmbiguousCreate decides ownership after a POST whose outcome is unknown. The API issues
// no correlation token, so ownership can only be inferred, and the inference must run the whole
// window: an early adoption could seize a concurrent caller's assignment while ours surfaces a poll
// later. Every fresh id observed anywhere in the window is accumulated, and if that set ever holds
// more than one id the outcome is declared ambiguous, because at most one of them can be ours.
// Only a single id, observed as the sole fresh candidate across the entire window and still present
// at its end, is adopted. Anything less provable returns the original failure with instructions,
// never a guess.
func (client *client) reconcileAmbiguousCreate(ctx context.Context, scope assignmentScope, request roleAssignmentRequestDto, baseline map[string]struct{}, createErr error) (*roleAssignmentDto, error) {
	seen := map[string]struct{}{}
	var last *roleAssignmentDto
	for attempt := 0; attempt < reconcileAttempts; attempt++ {
		if attempt > 0 {
			if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
				break
			}
		}
		assignments, err := client.ListRoleAssignments(ctx, scope)
		if err != nil {
			continue
		}

		last = nil
		for _, match := range matchingAssignments(assignments, request) {
			if _, preexisting := baseline[strings.ToLower(match.RoleAssignmentId)]; preexisting {
				continue
			}
			seen[strings.ToLower(match.RoleAssignmentId)] = struct{}{}
			adoptable := match
			last = &adoptable
		}
		if len(seen) > 1 {
			return nil, fmt.Errorf("the create failed ambiguously and %d distinct new assignments for this principal, role and scope were observed; at most one can be the one this create made, so none can be adopted safely. Inspect them, remove any duplicate, and import the survivor. Original failure: %w", len(seen), createErr)
		}
	}
	if len(seen) == 1 && last != nil {
		tflog.Debug(ctx, fmt.Sprintf("Create at scope %s failed ambiguously; %s was the only new assignment observed across the whole reconcile window, adopting it", request.Scope, last.RoleAssignmentId))
		return last, nil
	}
	return nil, fmt.Errorf("the create failed and it could not be proven whether the assignment was committed; list the scope's role assignments and import the assignment if it exists. Original failure: %w", createErr)
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
