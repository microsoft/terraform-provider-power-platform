// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_groups

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment_group_rule_set"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/tenant"
)

func newEnvironmentGroupClient(apiClient *api.Client, tenantClient tenant.Client, ruleSetClient environment_group_rule_set.Client) client {
	return client{
		Api:        apiClient,
		TenantApi:  tenantClient,
		RuleSetApi: ruleSetClient,
	}
}

type client struct {
	Api        *api.Client
	TenantApi  tenant.Client
	RuleSetApi environment_group_rule_set.Client
}

func (client *client) CreateEnvironmentGroup(ctx context.Context, environmentGroup environmentGroupDto) (*environmentGroupDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/environmentGroups",
	}

	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.BAP_2021_API_VERSION)
	apiUrl.RawQuery = values.Encode()

	newEnvironmentGroup := environmentGroupDto{}
	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, environmentGroup, []int{http.StatusCreated}, &newEnvironmentGroup)
	if err != nil {
		return nil, fmt.Errorf("failed to create environment group: %w", err)
	}

	return &newEnvironmentGroup, nil
}

// DeleteEnvironmentGroup deletes an environment group.
func (client *client) DeleteEnvironmentGroup(ctx context.Context, environmentGroupId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/environmentGroups/" + environmentGroupId,
	}

	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.BAP_2021_API_VERSION)
	apiUrl.RawQuery = values.Encode()

	resp, err := client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusConflict}, nil)
	if err != nil {
		if resp.HttpResponse.StatusCode == http.StatusConflict {
			if len(resp.BodyAsBytes) == 0 {
				return errors.New("failed to delete environment group")
			}

			body := string(resp.BodyAsBytes[:])
			if strings.Contains(body, "EnvironmentsInEnvironmentGroup") {
				return customerrors.WrapIntoProviderError(err, customerrors.ErrorCode(constants.ERROR_ENVIRONMENTS_IN_ENV_GROUP), "Failed to delete environment group because it contains environments")
			} else if strings.Contains(body, "PolicyAssignedToEnvironmentGroup") {
				return customerrors.WrapIntoProviderError(err, customerrors.ErrorCode(constants.ERROR_POLICY_ASSIGNED_TO_ENV_GROUP), "Failed to delete environment group because it has a policy assigned")
			}
			return errors.New(body)
		}
		return err
	}
	return nil
}

// updateEnvironmentGroup updates an environment group.
func (client *client) UpdateEnvironmentGroup(ctx context.Context, environmentGroupId string, environmentGroup environmentGroupDto) (*environmentGroupDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/environmentGroups/" + environmentGroupId,
	}

	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.BAP_2021_API_VERSION)
	apiUrl.RawQuery = values.Encode()

	updatedEnvironmentGroup := environmentGroupDto{}
	_, err := client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, environmentGroup, []int{http.StatusOK}, &updatedEnvironmentGroup)
	if err != nil {
		return nil, fmt.Errorf("failed to update environment group: %w", err)
	}

	return &updatedEnvironmentGroup, nil
}

// GetEnvironmentGroup gets an environment group.
func (client *client) GetEnvironmentGroup(ctx context.Context, environmentGroupId string) (*environmentGroupDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/environmentGroups/" + environmentGroupId,
	}

	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.BAP_2021_API_VERSION)
	apiUrl.RawQuery = values.Encode()

	environmentGroup := environmentGroupDto{}
	httpResponse, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNotFound}, &environmentGroup)
	if httpResponse.HttpResponse.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get environment group: %w", err)
	}

	return &environmentGroup, nil
}

func (client *client) GetEnvironmentsInEnvironmentGroup(ctx context.Context, environmentGroupId string) ([]environmentDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "providers/Microsoft.BusinessAppPlatform/scopes/admin/environments",
	}

	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.BAP_2021_API_VERSION)
	values.Add("$filter", fmt.Sprintf("properties/parentEnvironmentGroup/id eq %s", environmentGroupId))
	apiUrl.RawQuery = values.Encode()

	environments := environmentArrayDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &environments)
	if err != nil {
		return nil, fmt.Errorf("failed to get environments in environment group: %w", err)
	}

	return environments.Value, nil
}

func (client *client) RemoveEnvironmentFromEnvironmentGroup(ctx context.Context, environmentGroupId, environmentId string) error {
	tenantDto, err := client.TenantApi.GetTenant(ctx)
	if err != nil {
		return fmt.Errorf("failed to get tenant: %w", err)
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   helpers.BuildTenantHostUri(tenantDto.TenantId, client.Api.GetConfig().Urls.PowerPlatformUrl),
		Path:   fmt.Sprintf("/environmentmanagement/environmentGroups/%s/removeEnvironment/%s", environmentGroupId, environmentId),
	}

	values := url.Values{}
	values.Add("api-version", "1")
	apiUrl.RawQuery = values.Encode()

	_, err = client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, nil, []int{http.StatusAccepted}, nil)
	if err != nil {
		return fmt.Errorf("failed to remove environment from environment group: %w", err)
	}
	return nil
}
