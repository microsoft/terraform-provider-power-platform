// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"net/http"
	"net/url"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewTenantSettingsClient(api *api.ApiClient) TenantSettingsClient {
	return TenantSettingsClient{
		Api: api,
	}
}

type TenantSettingsClient struct {
	Api *api.ApiClient
}

func (client *TenantSettingsClient) GetTenant(ctx context.Context) (*TenantDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/tenant",
	}

	values := url.Values{}
	values.Add("api-version", "2020-08-01")
	apiUrl.RawQuery = values.Encode()

	tenant := TenantDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &tenant)
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

func (client *TenantSettingsClient) GetTenantSettings(ctx context.Context) (*TenantSettingsDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/listTenantSettings",
	}

	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	tenantSettings := TenantSettingsDto{}
	_, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, nil, []int{http.StatusOK}, &tenantSettings)
	if err != nil {
		return nil, err
	}
	return &tenantSettings, nil
}

func (client *TenantSettingsClient) UpdateTenantSettings(ctx context.Context, tenantSettings TenantSettingsDto) (*TenantSettingsDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/scopes/admin/updateTenantSettings",
	}

	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	var backendSettings TenantSettingsDto
	_, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, tenantSettings, []int{http.StatusOK}, &backendSettings)
	if err != nil {
		return nil, err
	}
	return &backendSettings, nil
}
