// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewLocationsClient(api *api.ApiClient) LocationsClient {
	return LocationsClient{
		Api: api,
	}
}

type LocationsClient struct {
	Api *api.ApiClient
}

func (client *LocationsClient) GetLocations(ctx context.Context) (LocationsDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/locations",
	}
	values := url.Values{
		"api-version": []string{"2023-06-01"},
	}
	apiUrl.RawQuery = values.Encode()

	locations := LocationsDto{}

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
