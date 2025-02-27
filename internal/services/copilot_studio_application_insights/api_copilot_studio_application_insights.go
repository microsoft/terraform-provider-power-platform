// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package copilot_studio_application_insights

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment"
)

func newCopilotStudioClient(apiClient *api.Client) client {
	return client{
		Api:               apiClient,
		EnvironmentClient: environment.NewEnvironmentClient(apiClient),
	}
}

type client struct {
	Api               *api.Client
	EnvironmentClient environment.Client
}

func (client *client) getCopilotStudioEndpoint(ctx context.Context, environmentId string) (string, error) {
	env, err := client.EnvironmentClient.GetEnvironment(ctx, environmentId)
	if err != nil {
		return "", err
	}
	if env == nil || env.Properties == nil || env.Properties.RuntimeEndpoints == nil || env.Properties.RuntimeEndpoints.PowerVirtualAgents == "" {
		return "", fmt.Errorf("Power Virtual Agents runtime endpoint is not available in the environment")
	}

	u, err := url.Parse(env.Properties.RuntimeEndpoints.PowerVirtualAgents)
	if err != nil {
		return "", err
	}

	return u.Host, nil
}

func parseImportId(importId string) (envId string, botId string, err error) {
	parts := strings.Split(importId, "_")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid import id format")
	}
	return parts[0], parts[1], nil
}

func (client *client) getCopilotStudioAppInsightsConfiguration(ctx context.Context, importId string) (*CopilotStudioAppInsightsDto, error) {
	environmentId, botId, err := parseImportId(importId)
	if err != nil {
		return nil, err
	}
	copilotStudioEndpoint, err := client.getCopilotStudioEndpoint(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	env, err := client.EnvironmentClient.GetEnvironment(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	if env == nil || env.Properties == nil || env.Properties.TenantId == "" {
		return nil, fmt.Errorf("TenantId value is not available in the environment properties")
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

	env, err := client.EnvironmentClient.GetEnvironment(ctx, copilotStudioAppInsightsConfig.EnvironmentId)
	if err != nil {
		return nil, err
	}

	if env == nil || env.Properties == nil || env.Properties.TenantId == "" {
		return nil, fmt.Errorf("TenantId value is not available in the environment properties")
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   copilotStudioEndpoint,
		Path:   fmt.Sprintf("/api/botmanagement/2022-01-15/environments/%s/bots/%s/applicationinsightsconfiguration", copilotStudioAppInsightsConfig.EnvironmentId, botId),
	}
	values := url.Values{}
	apiUrl.RawQuery = values.Encode()

	updatedCopilotStudioAppInsightsConfiguration := CopilotStudioAppInsightsDto{}

	resp, err := client.Api.Execute(ctx, []string{constants.COPILOT_SCOPE}, "PUT", apiUrl.String(), http.Header{"x-cci-tenantid": {env.Properties.TenantId}}, copilotStudioAppInsightsConfig, []int{http.StatusOK, http.StatusInternalServerError}, &updatedCopilotStudioAppInsightsConfiguration)
	if err != nil {
		return nil, err
	}
	if resp.HttpResponse.StatusCode == http.StatusInternalServerError {
		return nil, fmt.Errorf("Error updating Application Insights configuration: %s", string(resp.BodyAsBytes))
	}

	if len(updatedCopilotStudioAppInsightsConfiguration.Errors) > 0 {
		return nil, fmt.Errorf("Error updating Application Insights configuration: %s", updatedCopilotStudioAppInsightsConfiguration.Errors)
	}
	updatedCopilotStudioAppInsightsConfiguration.EnvironmentId = copilotStudioAppInsightsConfig.EnvironmentId
	updatedCopilotStudioAppInsightsConfiguration.BotId = botId
	return &updatedCopilotStudioAppInsightsConfiguration, nil
}
