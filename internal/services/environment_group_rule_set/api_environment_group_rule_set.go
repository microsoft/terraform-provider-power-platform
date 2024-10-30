// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_set

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/tenant"
)

func NewEnvironmentGroupRuleSetClient(apiClient *api.Client, tenantClient tenant.Client) Client {
	return Client{
		Api:       apiClient,
		TenantApi: tenantClient,
	}
}

type Client struct {
	Api       *api.Client
	TenantApi tenant.Client
}

func (client *Client) GetEnvironmentGroupRuleSet(ctx context.Context, environmentGroupId string) (*EnvironmentGroupRuleSetValueSetDto, error) {
	tenantDto, err := client.TenantApi.GetTenant(ctx)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   helpers.BuildTenantHostUri(tenantDto.TenantId, client.Api.GetConfig().Urls.PowerPlatformUrl),
		Path:   fmt.Sprintf("/governance/environmentGroups/%s/ruleSets", environmentGroupId),
	}

	values := url.Values{}
	values.Add("api-version", "2021-10-01-preview")
	apiUrl.RawQuery = values.Encode()

	environmentGroupRuleSet := environmentGroupRuleSetDto{}
	resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNoContent}, &environmentGroupRuleSet)
	if err != nil {
		return nil, err
	}

	if resp.HttpResponse.StatusCode == http.StatusNoContent {
		return nil, customerrors.WrapIntoProviderError(err, customerrors.ERROR_OBJECT_NOT_FOUND, "rule set '%s' not found")
	}

	if len(environmentGroupRuleSet.Value) == 0 {
		return nil, fmt.Errorf("no environment group ruleset found for environment group id %s", environmentGroupId)
	}

	return &environmentGroupRuleSet.Value[0], nil
}

func (client *Client) CreateEnvironmentGroupRuleSet(ctx context.Context, environmentGroupId string, newEnvironmentGroupRuleSet EnvironmentGroupRuleSetValueSetDto) (*EnvironmentGroupRuleSetValueSetDto, error) {
	tenantDto, err := client.TenantApi.GetTenant(ctx)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   helpers.BuildTenantHostUri(tenantDto.TenantId, client.Api.GetConfig().Urls.PowerPlatformUrl),
		Path:   fmt.Sprintf("/governance/environmentGroups/%s/ruleSets", environmentGroupId),
	}

	values := url.Values{}
	values.Add("api-version", "2021-10-01-preview")
	apiUrl.RawQuery = values.Encode()

	environmentGroupRuleSet := EnvironmentGroupRuleSetValueSetDto{}
	_, err = client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, newEnvironmentGroupRuleSet, []int{http.StatusCreated}, &environmentGroupRuleSet)

	if err != nil {
		return nil, err
	}

	if len(environmentGroupRuleSet.Parameters) == 0 {
		return nil, fmt.Errorf("no environment group ruleset found for environment group id %s", environmentGroupId)
	}

	return &environmentGroupRuleSet, nil
}

func (client *Client) UpdateEnvironmentGroupRuleSet(ctx context.Context, environmentGroupId string, newEnvironmentGroupRuleSet EnvironmentGroupRuleSetValueSetDto) (*EnvironmentGroupRuleSetValueSetDto, error) {
	tenantDto, err := client.TenantApi.GetTenant(ctx)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   helpers.BuildTenantHostUri(tenantDto.TenantId, client.Api.GetConfig().Urls.PowerPlatformUrl),
		Path:   fmt.Sprintf("/governance/ruleSets/%s", environmentGroupId),
	}

	values := url.Values{}
	values.Add("api-version", "2021-10-01-preview")
	apiUrl.RawQuery = values.Encode()

	environmentGroupRuleSet := EnvironmentGroupRuleSetValueSetDto{}
	_, err = client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, newEnvironmentGroupRuleSet, []int{http.StatusOK}, &environmentGroupRuleSet)

	if err != nil {
		return nil, err
	}

	return &environmentGroupRuleSet, nil
}

func (client *Client) DeleteEnvironmentGroupRuleSet(ctx context.Context, ruleSetId string) error {
	tenantDto, err := client.TenantApi.GetTenant(ctx)
	if err != nil {
		return err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   helpers.BuildTenantHostUri(tenantDto.TenantId, client.Api.GetConfig().Urls.PowerPlatformUrl),
		Path:   fmt.Sprintf("/governance/ruleSets/%s", ruleSetId),
	}

	values := url.Values{}
	values.Add("api-version", "2021-10-01-preview")
	apiUrl.RawQuery = values.Encode()

	_, err = client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)

	return err
}
