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

func NewLocationsClient(apiClient *api.Client) LocationsClient {
	return LocationsClient{
		Api: apiClient,
	}
}

type LocationsClient struct {
	Api *api.Client
}

func (client *LocationsClient) GetLocations(ctx context.Context) (Dto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/locations",
	}
	values := url.Values{
		"api-version": []string{"2023-06-01"},
	}
	apiUrl.RawQuery = values.Encode()

	locations := Dto{}

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
