// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant

import (
	"context"
	"net/http"
	"net/url"

	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewTenantClient(apiClient *api.ApiClient) ClientTenant {
	return ClientTenant{
		Api: apiClient,
	}
}

type ClientTenant struct {
	Api *api.ApiClient
}

func (client *ClientTenant) GetTenant(ctx context.Context) (*DtoTenant, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/tenant",
	}

	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	var dto DtoTenant

	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &dto)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}
