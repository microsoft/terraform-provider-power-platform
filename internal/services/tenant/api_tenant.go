// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant

import (
	"context"
	"net/http"
	"net/url"

	api "github.com/microsoft/terraform-provider-power-platform/internal/api"
)

func NewTenantClient(api *api.ApiClient) TenantClient {
	return TenantClient{
		Api: api,
	}
}

type TenantClient struct {
	Api *api.ApiClient
}

func (client *TenantClient) GetTenant(ctx context.Context) (*TenantDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/tenant",
	}

	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	var dto TenantDto

	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &dto)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}
