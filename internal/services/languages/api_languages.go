// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package languages

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

func NewLanguagesClient(apiClient *api.Client) Client {
	return Client{
		Api: apiClient,
	}
}

type Client struct {
	Api *api.Client
}

func (client *Client) GetLanguagesByLocation(ctx context.Context, location string) (ArrayDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/locations/%s/environmentLanguages", location),
	}
	apiUrl.RawQuery = url.Values{
		"api-version": []string{"2023-06-01"},
	}.Encode()

	languages := ArrayDto{}

	response, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return languages, err
	}

	defer response.Response.Body.Close()

	err = json.Unmarshal(response.BodyAsBytes, &languages)

	if err != nil {
		return languages, err
	}

	return languages, nil
}
