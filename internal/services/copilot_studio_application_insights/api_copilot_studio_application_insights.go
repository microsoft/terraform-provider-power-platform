// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package copilot_studio_application_insights

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

func newCopilotStudioClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}

type client struct {
	Api *api.Client
}

func (client *client) getCopilotStudioEndpoint(ctx context.Context, environmentId string) (string, error) {
	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return "", err
	}
	return env.Properties.RuntimeEndpoints.PowerVirtualAgents, nil
}

func (client *client) getEnvironment(ctx context.Context, environmentId string) (*environmentIdDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	env := environmentIdDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}

func (client *client) getCopilotStudioAppInsightsConfiguration(ctx context.Context, environmentId, botId string) (*copilotStudioAppInsightsDto, error) {
	copilotStudioEndpoint, err := client.getCopilotStudioEndpoint(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   copilotStudioEndpoint,
		Path:   fmt.Sprintf("/api/botmanagement/2022-01-15/environments/%s/bots/%s/applicationinsightsconfiguration", environmentId, botId),
	}
	values := url.Values{}
	apiUrl.RawQuery = values.Encode()

	copilotStudioAppInsights := copilotStudioAppInsightsDto{}

	_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &copilotStudioAppInsights)
	if err != nil {
		return nil, err
	}
	return &copilotStudioAppInsights, nil
}

func (client *client) updateCopilotStudioAppInsightsConfiguration(ctx context.Context, environmentId, botId string, copilotStudioAppInsightsConfig copilotStudioAppInsightsDto) (*copilotStudioAppInsightsDto, error) {
	copilotStudioEndpoint, err := client.getCopilotStudioEndpoint(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   copilotStudioEndpoint,
		Path:   fmt.Sprintf("/api/botmanagement/2022-01-15/environments/%s/bots/%s/applicationinsightsconfiguration", environmentId, botId),
	}
	values := url.Values{}
	apiUrl.RawQuery = values.Encode()

	updatedCopilotStudioAppInsightsConfiguration := copilotStudioAppInsightsDto{}

	_, err = client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, copilotStudioAppInsightsConfig, []int{http.StatusOK}, &updatedCopilotStudioAppInsightsConfiguration)
	if err != nil {
		return nil, err
	}

	return &updatedCopilotStudioAppInsightsConfiguration, nil
}
