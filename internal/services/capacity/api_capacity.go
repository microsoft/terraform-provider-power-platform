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

func (client *client) GetTenantCapacity(ctx context.Context) (*tenantCapacityDetailsDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   "/licensing/tenantCapacity",
	}
	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.APPLICATION_API_VERSION)
	apiUrl.RawQuery = values.Encode()

	var dto tenantCapacityDetailsDto

	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &dto)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant capacity details: %w", err)
	}

	return &dto, nil
}
