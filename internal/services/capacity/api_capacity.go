// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package capacity

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

type Client struct {
	Api *api.Client
}

func NewCapacityClient(clientApi *api.Client) Client {
	return Client{
		Api: clientApi,
	}
}

func (client *Client) GetTenantCapacity(ctx context.Context, tenantId string) (*Dto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.LicensingUrl,
		Path:   fmt.Sprintf("/v0.1-alpha/tenants/%s/TenantCapacity", tenantId),
	}

	var dto Dto

	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &dto)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}
