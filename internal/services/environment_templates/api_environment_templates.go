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

func newEnvironmentTemplatesClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}

type client struct {
	Api *api.Client
}

func (client *client) GetEnvironmentTemplatesByLocation(ctx context.Context, location string) (environmentTemplateDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/locations/%s/templates", location),
	}
	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.BAP_API_VERSION)
	apiUrl.RawQuery = values.Encode()

	templates := environmentTemplateDto{}

	response, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return templates, fmt.Errorf("failed to execute API request for environment templates: %w", err)
	}

	defer response.HttpResponse.Body.Close()

	err = json.Unmarshal(response.BodyAsBytes, &templates)
	if err != nil {
		return templates, fmt.Errorf("failed to unmarshal environment templates response: %w", err)
	}

	return templates, nil
}
