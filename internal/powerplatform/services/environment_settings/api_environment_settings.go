// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewEnvironmentSettingsClient(api *api.ApiClient) EnvironmentSettingsClient {
	return EnvironmentSettingsClient{
		Api: api,
	}
}

type EnvironmentSettingsClient struct {
	Api *api.ApiClient
}

func (client *EnvironmentSettingsClient) DataverseExists(ctx context.Context, environmentId string) (bool, error) {

	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return false, err
	}
	return env.Properties.LinkedEnvironmentMetadata.InstanceURL != "", nil
}

func (client *EnvironmentSettingsClient) GetEnvironmentSettings(ctx context.Context, environmentId string) (*EnvironmentSettingsDto, error) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   "/api/data/v9.0/organizations",
	}

	environmentSettings := EnvironmentSettingsValueDto{}
	_, err = client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &environmentSettings)
	if err != nil {
		return nil, err
	}
	return &environmentSettings.Value[0], nil
}

func (client *EnvironmentSettingsClient) UpdateEnvironmentSettings(ctx context.Context, environmentId string, environmentSettings EnvironmentSettingsDto) (*EnvironmentSettingsDto, error) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	settings, err := client.GetEnvironmentSettings(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   fmt.Sprintf("/api/data/v9.0/organizations(%s)", *settings.OrganizationId),
	}

	_, err = client.Api.Execute(ctx, "PATCH", apiUrl.String(), nil, environmentSettings, []int{http.StatusNoContent}, nil)
	if err != nil {
		return nil, err
	}

	return client.GetEnvironmentSettings(ctx, environmentId)
}

func (client *EnvironmentSettingsClient) GetEnvironmentUrlById(ctx context.Context, environmentId string) (string, error) {
	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return "", err
	}
	environmentUrl := strings.TrimSuffix(env.Properties.LinkedEnvironmentMetadata.InstanceURL, "/")
	return environmentUrl, nil
}

func (client *EnvironmentSettingsClient) getEnvironment(ctx context.Context, environmentId string) (*EnvironmentIdDto, error) {

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	env := EnvironmentIdDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}
