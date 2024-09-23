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

func newCapacityClient(clientApi *api.Client) client {
	return client{
		Api: clientApi,
	}
}

type client struct {
	Api *api.Client
}

func (client *client) GetTenantCapacity(ctx context.Context, tenantId string) (*capacityDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.LicensingUrl,
		Path:   fmt.Sprintf("/v0.1-alpha/tenants/%s/TenantCapacity", tenantId),
	}

	var dto capacityDto

	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &dto)
	if err != nil {
		return nil, err
	}

	return &dto, nil
}
