// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

const apiVersion = "2024-10-01"

func NewRoleBasedAccessClient(apiClient *api.Client) RoleBasedAccessClient {
	return RoleBasedAccessClient{
		Api: apiClient,
	}
}

type RoleBasedAccessClient struct {
	Api *api.Client
}

// CreateRoleAssignment creates a role assignment at tenant level.
func (client *RoleBasedAccessClient) CreateRoleAssignment(ctx context.Context, request roleAssignmentRequestDto) (*roleAssignmentDto, error) {
	return client.createRoleAssignmentAtPath(ctx, "/authorization/roleAssignments", request)
}

// CreateEnvironmentGroupRoleAssignment creates a role assignment for an environment group.
func (client *RoleBasedAccessClient) CreateEnvironmentGroupRoleAssignment(ctx context.Context, environmentGroupId string, request roleAssignmentRequestDto) (*roleAssignmentDto, error) {
	path := fmt.Sprintf("/authorization/environmentGroups/%s/roleAssignments", environmentGroupId)
	return client.createRoleAssignmentAtPath(ctx, path, request)
}

// CreateEnvironmentRoleAssignment creates a role assignment for an environment.
func (client *RoleBasedAccessClient) CreateEnvironmentRoleAssignment(ctx context.Context, environmentId string, request roleAssignmentRequestDto) (*roleAssignmentDto, error) {
	path := fmt.Sprintf("/authorization/environments/%s/roleAssignments", environmentId)
	return client.createRoleAssignmentAtPath(ctx, path, request)
}

func (client *RoleBasedAccessClient) createRoleAssignmentAtPath(ctx context.Context, path string, request roleAssignmentRequestDto) (*roleAssignmentDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   path,
	}
	values := url.Values{}
	values.Add("api-version", apiVersion)
	apiUrl.RawQuery = values.Encode()

	response := roleAssignmentDto{}
	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, request, []int{http.StatusCreated}, &response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// ListRoleAssignments lists role assignments at tenant level.
func (client *RoleBasedAccessClient) ListRoleAssignments(ctx context.Context) ([]roleAssignmentDto, error) {
	return client.listRoleAssignmentsAtPath(ctx, "/authorization/roleAssignments")
}

// ListEnvironmentGroupRoleAssignments lists role assignments for an environment group.
func (client *RoleBasedAccessClient) ListEnvironmentGroupRoleAssignments(ctx context.Context, environmentGroupId string) ([]roleAssignmentDto, error) {
	path := fmt.Sprintf("/authorization/environmentGroups/%s/roleAssignments", environmentGroupId)
	return client.listRoleAssignmentsAtPath(ctx, path)
}

// ListEnvironmentRoleAssignments lists role assignments for an environment.
func (client *RoleBasedAccessClient) ListEnvironmentRoleAssignments(ctx context.Context, environmentId string) ([]roleAssignmentDto, error) {
	path := fmt.Sprintf("/authorization/environments/%s/roleAssignments", environmentId)
	return client.listRoleAssignmentsAtPath(ctx, path)
}

func (client *RoleBasedAccessClient) listRoleAssignmentsAtPath(ctx context.Context, path string) ([]roleAssignmentDto, error) {
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
func (client *RoleBasedAccessClient) DeleteRoleAssignment(ctx context.Context, roleAssignmentId string) error {
	path := fmt.Sprintf("/authorization/roleAssignments/%s", roleAssignmentId)
	return client.deleteRoleAssignmentAtPath(ctx, path)
}

// DeleteEnvironmentGroupRoleAssignment deletes a role assignment for an environment group.
func (client *RoleBasedAccessClient) DeleteEnvironmentGroupRoleAssignment(ctx context.Context, environmentGroupId, roleAssignmentId string) error {
	path := fmt.Sprintf("/authorization/environmentGroups/%s/roleAssignments/%s", environmentGroupId, roleAssignmentId)
	return client.deleteRoleAssignmentAtPath(ctx, path)
}

// DeleteEnvironmentRoleAssignment deletes a role assignment for an environment.
func (client *RoleBasedAccessClient) DeleteEnvironmentRoleAssignment(ctx context.Context, environmentId, roleAssignmentId string) error {
	path := fmt.Sprintf("/authorization/environments/%s/roleAssignments/%s", environmentId, roleAssignmentId)
	return client.deleteRoleAssignmentAtPath(ctx, path)
}

func (client *RoleBasedAccessClient) deleteRoleAssignmentAtPath(ctx context.Context, path string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   path,
	}
	values := url.Values{}
	values.Add("api-version", apiVersion)
	apiUrl.RawQuery = values.Encode()

	_, err := client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusNoContent}, nil)
	return err
}

// ListRoleDefinitions lists available role definitions.
func (client *RoleBasedAccessClient) ListRoleDefinitions(ctx context.Context) ([]roleDefinitionDto, error) {
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
