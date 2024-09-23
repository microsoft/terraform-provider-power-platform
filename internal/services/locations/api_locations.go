// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package locations

import (
	"context"
	"encoding/json"
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
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/locations",
	}
	values := url.Values{
		"api-version": []string{"2023-06-01"},
	}
	apiUrl.RawQuery = values.Encode()

	locations := locationDto{}

	response, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return locations, err
	}

	defer response.Response.Body.Close()

	err = json.Unmarshal(response.BodyAsBytes, &locations)

	if err != nil {
		return locations, err
	}

	return locations, nil
}
