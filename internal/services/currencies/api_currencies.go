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
)

func NewCurrenciesClient(api *api.ApiClient) CurrenciesClient {
	return CurrenciesClient{
		Api: api,
	}
}

type CurrenciesClient struct {
	Api *api.ApiClient
}

func (client *CurrenciesClient) GetCurrenciesByLocation(ctx context.Context, location string) (CurrenciesDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/locations/%s/environmentCurrencies", location),
	}
	apiUrl.RawQuery = url.Values{
		"api-version": []string{"2023-06-01"},
	}.Encode()

	currencies := CurrenciesDto{}

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
