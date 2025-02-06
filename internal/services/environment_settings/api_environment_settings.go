// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_settings

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
)

func newEnvironmentSettingsClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}

type client struct {
	Api *api.Client
}

func (client *client) DataverseExists(ctx context.Context, environmentId string) (bool, error) {
	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return false, err
	}
	return env.Properties.LinkedEnvironmentMetadata.InstanceURL != "", nil
}

func (client *client) GetEnvironmentSettings(ctx context.Context, environmentId string) (*environmentSettingsDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.0/organizations",
	}

	environmentSettings := environmentSettingsValueDto{}
	_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &environmentSettings)
	if err != nil {
		return nil, err
	}
	return &environmentSettings.Value[0], nil
}

func (client *client) UpdateEnvironmentSettings(ctx context.Context, environmentId string, environmentSettings environmentSettingsDto) (*environmentSettingsDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	settings, err := client.GetEnvironmentSettings(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/v9.0/organizations(%s)", *settings.OrganizationId),
	}

	resp, err := client.Api.Execute(ctx, nil, "PATCH", apiUrl.String(), nil, environmentSettings, []int{http.StatusNoContent, http.StatusInternalServerError}, nil)
	if resp.HttpResponse.StatusCode == http.StatusInternalServerError {
		return nil, customerrors.WrapIntoProviderError(nil, customerrors.ERROR_ENVIRONMENT_SETTINGS_FAILED, string(resp.BodyAsBytes))
	}
	if err != nil {
		return nil, err
	}

	return client.GetEnvironmentSettings(ctx, environmentId)
}

func (client *client) GetEnvironmentHostById(ctx context.Context, environmentId string) (string, error) {
	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return "", err
	}
	environmentUrl := strings.TrimSuffix(env.Properties.LinkedEnvironmentMetadata.InstanceURL, "/")
	if environmentUrl == "" {
		return "", customerrors.WrapIntoProviderError(nil, customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND, "environment url not found, please check if the environment has dataverse linked")
	}

	envUrl, err := url.Parse(environmentUrl)
	if err != nil {
		return "", err
	}
	return envUrl.Host, nil
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
