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
	if env.Properties.RuntimeEndpoints.PowerVirtualAgents == "" {
		return "", fmt.Errorf("Power Virtual Agents runtime endpoint is not available in the environment")
	}

	u, err := url.Parse(env.Properties.RuntimeEndpoints.PowerVirtualAgents)
	if err != nil {
		return "", err
	}

	return u.Host, nil
}

func (client *client) getEnvironment(ctx context.Context, environmentId string) (*EnvironmentIdDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	env := EnvironmentIdDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}

func (client *client) getCopilotStudioAppInsightsConfiguration(ctx context.Context, environmentId, botId string) (*CopilotStudioAppInsightsDto, error) {
	copilotStudioEndpoint, err := client.getCopilotStudioEndpoint(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	env, err := client.getEnvironment(ctx, environmentId)
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

	copilotStudioAppInsights := CopilotStudioAppInsightsDto{}
	_, err = client.Api.Execute(ctx, []string{constants.COPILOT_SCOPE}, "GET", apiUrl.String(), http.Header{"x-cci-tenantid": {env.Properties.TenantId}}, nil, []int{http.StatusOK}, &copilotStudioAppInsights)
	if err != nil {
		return nil, err
	}
	copilotStudioAppInsights.EnvironmentId = environmentId
	copilotStudioAppInsights.BotId = botId
	return &copilotStudioAppInsights, nil
}

func (client *client) updateCopilotStudioAppInsightsConfiguration(ctx context.Context, copilotStudioAppInsightsConfig CopilotStudioAppInsightsDto, botId string) (*CopilotStudioAppInsightsDto, error) {
	copilotStudioEndpoint, err := client.getCopilotStudioEndpoint(ctx, copilotStudioAppInsightsConfig.EnvironmentId)
	if err != nil {
		return nil, err
	}

	env, err := client.getEnvironment(ctx, copilotStudioAppInsightsConfig.EnvironmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   copilotStudioEndpoint,
		Path:   fmt.Sprintf("/api/botmanagement/2022-01-15/environments/%s/bots/%s/applicationinsightsconfiguration", copilotStudioAppInsightsConfig.EnvironmentId, botId),
	}
	values := url.Values{}
	apiUrl.RawQuery = values.Encode()

	updatedCopilotStudioAppInsightsConfiguration := CopilotStudioAppInsightsDto{}

	_, err = client.Api.Execute(ctx, []string{constants.COPILOT_SCOPE}, "PUT", apiUrl.String(), http.Header{"x-cci-tenantid": {env.Properties.TenantId}}, copilotStudioAppInsightsConfig, []int{http.StatusOK}, &updatedCopilotStudioAppInsightsConfiguration)
	if err != nil {
		return nil, err
	}
	if len(updatedCopilotStudioAppInsightsConfiguration.Errors) > 0 {
		return nil, fmt.Errorf("Error updating Application Insights configuration: %s", updatedCopilotStudioAppInsightsConfiguration.Errors)
	}
	updatedCopilotStudioAppInsightsConfiguration.EnvironmentId = copilotStudioAppInsightsConfig.EnvironmentId
	updatedCopilotStudioAppInsightsConfiguration.BotId = botId
	return &updatedCopilotStudioAppInsightsConfiguration, nil
}
