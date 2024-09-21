// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_templates

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

func NewEnvironmentTemplatesClient(apiClient *api.Client) EnvironmentTemplatesClient {
	return EnvironmentTemplatesClient{
		Api: apiClient,
	}
}

type EnvironmentTemplatesClient struct {
	Api *api.Client
}

func (client *EnvironmentTemplatesClient) GetEnvironmentTemplatesByLocation(ctx context.Context, location string) (Dto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/locations/%s/templates", location),
	}
	apiUrl.RawQuery = url.Values{
		"api-version": []string{"2023-06-01"},
	}.Encode()

	templates := Dto{}

	response, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
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
