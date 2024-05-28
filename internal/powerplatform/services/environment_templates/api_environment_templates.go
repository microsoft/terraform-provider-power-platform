// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewEnvironmentTemplatesClient(api *api.ApiClient) EnvironmentTemplatesClient {
	return EnvironmentTemplatesClient{
		Api: api,
	}
}

type EnvironmentTemplatesClient struct {
	Api *api.ApiClient
}

func (client *EnvironmentTemplatesClient) GetEnvironmentTemplatesByLocation(ctx context.Context, location string) (EnvironmentTemplatesDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/locations/%s/templates", location),
	}
	apiUrl.RawQuery = url.Values{
		"api-version": []string{"2023-06-01"},
	}.Encode()

	templates := EnvironmentTemplatesDto{}

	response, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return templates, err
	}

	defer response.Response.Body.Close()

	err = json.Unmarshal(response.BodyAsBytes, &templates)

	if err != nil {
		return templates, err
	}

	return templates, nil
}
