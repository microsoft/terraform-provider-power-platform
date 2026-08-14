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

// tenantScope returns the fully qualified tenant scope string required by the RBAC API.
func (client *client) tenantScope(ctx context.Context) (string, error) {
	tenantDto, err := client.TenantApi.GetTenant(ctx)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("/tenants/%s", tenantDto.TenantId), nil
}

// CreateRoleAssignment creates a role assignment at tenant level.
func (client *client) CreateRoleAssignment(ctx context.Context, request roleAssignmentRequestDto) (*roleAssignmentDto, error) {
	scope, err := client.tenantScope(ctx)
	if err != nil {
		return nil, err
	}
	request.Scope = scope
	return client.createRoleAssignmentAtPath(ctx, "/authorization/roleAssignments", request)
}

// CreateEnvironmentGroupRoleAssignment creates a role assignment for an environment group.
func (client *client) CreateEnvironmentGroupRoleAssignment(ctx context.Context, environmentGroupId string, request roleAssignmentRequestDto) (*roleAssignmentDto, error) {
	scope, err := client.tenantScope(ctx)
	if err != nil {
		return nil, err
	}
	request.Scope = fmt.Sprintf("%s/environmentGroups/%s", scope, environmentGroupId)
	path := fmt.Sprintf("/authorization/environmentGroups/%s/roleAssignments", environmentGroupId)
	return client.createRoleAssignmentAtPath(ctx, path, request)
}

// CreateEnvironmentRoleAssignment creates a role assignment for an environment.
func (client *client) CreateEnvironmentRoleAssignment(ctx context.Context, environmentId string, request roleAssignmentRequestDto) (*roleAssignmentDto, error) {
	scope, err := client.tenantScope(ctx)
	if err != nil {
		return nil, err
	}
	request.Scope = fmt.Sprintf("%s/environments/%s", scope, environmentId)
	path := fmt.Sprintf("/authorization/environments/%s/roleAssignments", environmentId)
	return client.createRoleAssignmentAtPath(ctx, path, request)
}

func (client *client) createRoleAssignmentAtPath(ctx context.Context, path string, request roleAssignmentRequestDto) (*roleAssignmentDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   path,
	}
	values := url.Values{}
	values.Add("api-version", apiVersion)
	apiUrl.RawQuery = values.Encode()

	response := roleAssignmentDto{}
	for {
		_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, request, []int{http.StatusOK, http.StatusCreated}, &response)
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

// ListRoleAssignments lists role assignments at tenant level.
func (client *client) ListRoleAssignments(ctx context.Context) ([]roleAssignmentDto, error) {
	return client.listRoleAssignmentsAtPath(ctx, "/authorization/roleAssignments")
}

// ListEnvironmentGroupRoleAssignments lists role assignments for an environment group.
func (client *client) ListEnvironmentGroupRoleAssignments(ctx context.Context, environmentGroupId string) ([]roleAssignmentDto, error) {
	path := fmt.Sprintf("/authorization/environmentGroups/%s/roleAssignments", environmentGroupId)
	return client.listRoleAssignmentsAtPath(ctx, path)
}

// ListEnvironmentRoleAssignments lists role assignments for an environment.
func (client *client) ListEnvironmentRoleAssignments(ctx context.Context, environmentId string) ([]roleAssignmentDto, error) {
	path := fmt.Sprintf("/authorization/environments/%s/roleAssignments", environmentId)
	return client.listRoleAssignmentsAtPath(ctx, path)
}

func (client *client) listRoleAssignmentsAtPath(ctx context.Context, path string) ([]roleAssignmentDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   path,
	}
	values := url.Values{}
	values.Add("api-version", apiVersion)
	apiUrl.RawQuery = values.Encode()

	response := roleAssignmentsListDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &response)
	if err != nil {
		return nil, err
	}
	return response.Value, nil
}

// DeleteRoleAssignment deletes a role assignment at tenant level.
func (client *client) DeleteRoleAssignment(ctx context.Context, roleAssignmentId string) error {
	path := fmt.Sprintf("/authorization/roleAssignments/%s", roleAssignmentId)
	return client.deleteRoleAssignmentAtPath(ctx, path)
}

// DeleteEnvironmentGroupRoleAssignment deletes a role assignment for an environment group.
func (client *client) DeleteEnvironmentGroupRoleAssignment(ctx context.Context, environmentGroupId, roleAssignmentId string) error {
	path := fmt.Sprintf("/authorization/environmentGroups/%s/roleAssignments/%s", environmentGroupId, roleAssignmentId)
	return client.deleteRoleAssignmentAtPath(ctx, path)
}

// DeleteEnvironmentRoleAssignment deletes a role assignment for an environment.
func (client *client) DeleteEnvironmentRoleAssignment(ctx context.Context, environmentId, roleAssignmentId string) error {
	path := fmt.Sprintf("/authorization/environments/%s/roleAssignments/%s", environmentId, roleAssignmentId)
	return client.deleteRoleAssignmentAtPath(ctx, path)
}

func (client *client) deleteRoleAssignmentAtPath(ctx context.Context, path string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   path,
	}
	values := url.Values{}
	values.Add("api-version", apiVersion)
	apiUrl.RawQuery = values.Encode()

	_, err := client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNoContent}, nil)
	return err
}

// ListRoleDefinitions lists available role definitions.
func (client *client) ListRoleDefinitions(ctx context.Context) ([]roleDefinitionDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   "/authorization/roleDefinitions",
	}
	values := url.Values{}
	values.Add("api-version", apiVersion)
	apiUrl.RawQuery = values.Encode()

	response := roleDefinitionsListDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &response)
	if err != nil {
		return nil, err
	}
	return response.Value, nil
}
