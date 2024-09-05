// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package capacity

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

type CapacityClient struct {
	Api *api.ApiClient
}

func NewCapacityClient(api *api.ApiClient) CapacityClient {
	return CapacityClient{
		Api: api,
	}
}


func (client *CapacityClient) GetTenantCapacity(ctx context.Context, tenantId string) (*TenantCapacityDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.LicensingUrl,
		Path:   fmt.Sprintf("/v0.1-alpha/tenants/%s/TenantCapacity", tenantId),
	}

	var dto TenantCapacityDto

	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &dto)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}
