// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managed_environment

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment"
)

func newManagedEnvironmentClient(apiClient *api.Client) client {
	return client{
		Api:               apiClient,
		environmentClient: environment.NewEnvironmentClient(apiClient),
	}
}

type client struct {
	Api               *api.Client
	environmentClient environment.Client
}

func (client *client) GetManagedEnvironmentSettings(ctx context.Context, environmentId string) (*environment.GovernanceConfigurationDto, error) {
	managedEnvSettings, err := client.environmentClient.GetEnvironment(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	return &managedEnvSettings.Properties.GovernanceConfiguration, nil
}

func (client *client) EnableManagedEnvironment(ctx context.Context, managedEnvSettings environment.GovernanceConfigurationDto, environmentId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/environments/%s/governanceConfiguration", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	apiResponse, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, managedEnvSettings, []int{http.StatusNoContent, http.StatusAccepted}, nil)
	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Managed Environment Enablement Operation HTTP Status: '"+apiResponse.Response.Status+"'")

	tflog.Debug(ctx, "Waiting for Managed Environment Enablement Operation to complete")
	_, err = client.Api.DoWaitForLifecycleOperationStatus(ctx, apiResponse)
	if err != nil {
		return err
	}
	return nil
}

func (client *client) DisableManagedEnvironment(ctx context.Context, environmentId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/environments/%s/governanceConfiguration", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	managedEnv := environment.GovernanceConfigurationDto{
		ProtectionLevel: "Basic",
	}

	apiResponse, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, managedEnv, []int{http.StatusAccepted}, nil)
	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Managed Environment Disablement Operation HTTP Status: '"+apiResponse.Response.Status+"'")
	tflog.Debug(ctx, "Waiting for Managed Environment Disablement Operation to complete")

	_, err = client.Api.DoWaitForLifecycleOperationStatus(ctx, apiResponse)
	if err != nil {
		return err
	}
	return nil
}
