// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package locations

import (
	"context"
	"net/http"
	"net/url"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

func newLocationsClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}

type client struct {
	Api *api.Client
}

func (client *client) GetLocations(ctx context.Context) (locationDto, error) {
	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.BAP_API_VERSION)
	
	apiUrl := &url.URL{
		Scheme:   constants.HTTPS,
		Host:     client.Api.GetConfig().Urls.BapiUrl,
		Path:     "/providers/Microsoft.BusinessAppPlatform/locations",
		RawQuery: values.Encode(),
	}

	var locations locationDto
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &locations)
	return locations, err
}
