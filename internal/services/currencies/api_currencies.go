// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package currencies

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

func NewCurrenciesClient(apiClient *api.Client) Client {
	return Client{
		Api: apiClient,
	}
}

type Client struct {
	Api *api.Client
}

func (client *Client) GetCurrenciesByLocation(ctx context.Context, location string) (Dto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/locations/%s/environmentCurrencies", location),
	}
	apiUrl.RawQuery = url.Values{
		"api-version": []string{"2023-06-01"},
	}.Encode()

	currencies := Dto{}

	response, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return currencies, err
	}

	defer response.Response.Body.Close()

	err = json.Unmarshal(response.BodyAsBytes, &currencies)

	if err != nil {
		return currencies, err
	}

	return currencies, nil
}
