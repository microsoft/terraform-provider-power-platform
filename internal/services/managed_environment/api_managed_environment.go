// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managed_environment

import (
	"context"
	"errors"
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

func (client *client) EnableManagedEnvironment(ctx context.Context, managedEnvSettings environment.GovernanceConfigurationDto, environmentId string) error {
	return client.enableManagedEnvironmentWithRetry(ctx, managedEnvSettings, environmentId, 0)
}

func (client *client) enableManagedEnvironmentWithRetry(ctx context.Context, managedEnvSettings environment.GovernanceConfigurationDto, environmentId string, retryCount int) error {
	const maxRetries = 10
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/environments/%s/governanceConfiguration", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	apiResponse, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, managedEnvSettings, []int{http.StatusNoContent, http.StatusAccepted, http.StatusConflict}, nil)
	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Managed Environment Enablement Operation HTTP Status: '"+apiResponse.HttpResponse.Status+"'")

	tflog.Debug(ctx, "Waiting for Managed Environment Enablement Operation to complete")
	lifecycleResponse, err := client.Api.DoWaitForLifecycleOperationStatus(ctx, apiResponse)
	if err != nil {
		return err
	}
	if lifecycleResponse != nil && lifecycleResponse.State.Id == "Failed" {
		if retryCount >= maxRetries {
			return fmt.Errorf("maximum retries (%d) reached for EnableManagedEnvironment on lifecycle failure", maxRetries)
		}
		if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return err
		}
		tflog.Info(ctx, "Managed Environment Enablement Operation failed. Retrying...")
		return client.enableManagedEnvironmentWithRetry(ctx, managedEnvSettings, environmentId, retryCount+1)
	}
	return nil
}

func (client *client) DisableManagedEnvironment(ctx context.Context, environmentId string) error {
	return client.disableManagedEnvironmentWithRetry(ctx, environmentId, 0)
}

func (client *client) disableManagedEnvironmentWithRetry(ctx context.Context, environmentId string, retryCount int) error {
	const maxRetries = 10
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

	apiResponse, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, managedEnv, []int{http.StatusAccepted, http.StatusConflict}, nil)
	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Managed Environment Disablement Operation HTTP Status: '"+apiResponse.HttpResponse.Status+"'")
	tflog.Debug(ctx, "Waiting for Managed Environment Disablement Operation to complete")

	lifecycleResponse, err := client.Api.DoWaitForLifecycleOperationStatus(ctx, apiResponse)
	if err != nil {
		return err
	}
	if lifecycleResponse != nil && lifecycleResponse.State.Id == "Failed" {
		if retryCount >= maxRetries {
			return fmt.Errorf("maximum retries (%d) reached for DisableManagedEnvironment on lifecycle failure", maxRetries)
		}
		if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return err
		}
		tflog.Info(ctx, "Managed Environment Disablement Operation failed. Retrying...")
		return client.disableManagedEnvironmentWithRetry(ctx, environmentId, retryCount+1)
	}
	return nil
}

type SolutionCheckerRule struct {
	Code string `json:"code"`
}

func (client *client) FetchSolutionCheckerRules(ctx context.Context, environmentId string) ([]string, error) {
	if client.environmentClient == (environment.Client{}) {
		return nil, errors.New("environmentClient is not initialized")
	}

	env, err := client.environmentClient.GetEnvironment(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	if env.Properties.RuntimeEndpoints.PowerAppsAdvisor == "" {
		return nil, errors.New("PowerAppsAdvisor URL is empty")
	}

	powerAppsAdvisorUrl, err := url.Parse(env.Properties.RuntimeEndpoints.PowerAppsAdvisor)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PowerAppsAdvisor URL: %w", err)
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   powerAppsAdvisorUrl.Host,
		Path:   "/api/rule",
	}
	values := url.Values{}
	values.Add("api-version", "2.0")
	// Currently, the ruleset is always the same for all regions
	values.Add("ruleset", constants.SOLUTION_CHECKER_RULESET_ID)
	apiUrl.RawQuery = values.Encode()

	tflog.Debug(ctx, fmt.Sprintf("Constructed API URL: %s", apiUrl.String()))

	solutionCheckerRulesArrayDto := []SolutionCheckerRule{}
	_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &solutionCheckerRulesArrayDto)
	if err != nil {
		return nil, err
	}

	var codes []string
	for _, rule := range solutionCheckerRulesArrayDto {
		codes = append(codes, rule.Code)
	}

	return codes, nil
}
