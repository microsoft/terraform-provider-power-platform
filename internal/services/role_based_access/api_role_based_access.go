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
//   - if a matching assignment already exists the create fails and hands over its import id,
//     like any idiomatic resource: importing is the way to manage an existing grant. Automatic
//     adoption is deliberately absent, because a create cannot tell a fresh configuration from
//     the create leg of a replacement, and adopting during a create_before_destroy replacement
//     would let the deposed instance destroy the very grant just adopted;
//   - if several already exist the create fails too: deduplicate or import one;
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
				return nil, fmt.Errorf("interrupted while waiting for scope %s to become visible for role assignments: %w (last list failure: %v)", request.Scope, sleepErr, err)
			}
			continue
		}
		return nil, fmt.Errorf("could not list the existing role assignments before creating, so it is unknown whether one for this relationship already exists, and creating blindly could duplicate or mismanage it: %w", err)
	}

	matches := matchingAssignments(existing, request)
	if len(matches) == 1 {
		return nil, fmt.Errorf("a role assignment for this principal, role and scope already exists; import it instead of creating it: terraform import <address> %q. If this create is the forced recreate of the assignment Terraform already manages, remove create_before_destroy or untaint the resource, since the existing assignment is that one", scope.importId(matches[0].RoleAssignmentId))
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
				// The propagation 400 proved nothing was committed, so this stays a plain error.
				return nil, fmt.Errorf("interrupted while waiting for scope %s to accept role assignments: %w (last propagation rejection: %v)", request.Scope, sleepErr, err)
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
// definitive HTTP rejection (400, 401, 403, 404, 409) means the server refused it, and an error
// raised before the request was sent (scope, url, token or request construction) proves nothing
// was committed; a failure on or after the wire without a status, which the api client marks with
// RequestSentError, leaves the outcome unknown.
func isAmbiguousCreateFailure(err error) bool {
	var httpErr customerrors.UnexpectedHttpStatusCodeError
	if errors.As(err, &httpErr) {
		return httpErr.StatusCode == http.StatusRequestTimeout ||
			httpErr.StatusCode == http.StatusTooManyRequests ||
			httpErr.StatusCode >= http.StatusInternalServerError
	}
	var sent api.RequestSentError
	return errors.As(err, &sent)
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

// ResolveRoleDefinitionByName resolves a role definition name to its definition. Names are matched
// case-insensitively, since Microsoft has recased display names before. Exactly one match is
// required: none is a configuration error, and several cannot be told apart by name, so both fail
// with instructions rather than guessing.
func (client *client) ResolveRoleDefinitionByName(ctx context.Context, name string) (*roleDefinitionDto, error) {
	definitions, err := client.ListRoleDefinitions(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not list the role definitions to resolve the name %q: %w", name, err)
	}
	var matches []roleDefinitionDto
	for i := range definitions {
		if strings.EqualFold(definitions[i].RoleDefinitionName, name) {
			matches = append(matches, definitions[i])
		}
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf("no role definition is named %q; the catalogue holds %d definitions, which the powerplatform_role_definitions data source lists", name, len(definitions))
	}
	if len(matches) > 1 {
		ids := make([]string, len(matches))
		for i := range matches {
			ids[i] = matches[i].RoleDefinitionId
		}
		return nil, fmt.Errorf("%d role definitions are named %q (%s); use role_definition_id to pick one", len(matches), name, strings.Join(ids, ", "))
	}
	return &matches[0], nil
}

// RoleDefinitionNameById returns the display name for a role definition id, or "" when the
// catalogue cannot be listed or does not hold the id. The name is cosmetic alongside the id
// identity, so this lookup cannot be allowed to starve a mutation: it sends exactly one request,
// because a retried request against a persistently failing catalogue would consume the operation's
// whole context. It runs only from Read, where a single stalled request can still occupy the read
// until its timeout, and that is the accepted cost of the courtesy.
func (client *client) RoleDefinitionNameById(ctx context.Context, roleDefinitionId string) string {
	response := roleDefinitionsListDto{}
	_, err := client.Api.ExecuteWithoutRetry(ctx, nil, "GET", client.url("/authorization/roleDefinitions"), nil, nil, []int{http.StatusOK}, &response)
	if err != nil {
		tflog.Debug(ctx, fmt.Sprintf("Could not list role definitions to name %s: %s", roleDefinitionId, err))
		return ""
	}
	for i := range response.Value {
		if strings.EqualFold(response.Value[i].RoleDefinitionId, roleDefinitionId) {
			return response.Value[i].RoleDefinitionName
		}
	}
	return ""
}
