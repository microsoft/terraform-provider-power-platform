// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package disaster_recovery

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

func newDisasterRecoveryClient(apiClient *api.Client) client {
	return client{
		Api:               apiClient,
		environmentClient: environment.NewEnvironmentClient(apiClient),
	}
}

type client struct {
	Api               *api.Client
	environmentClient environment.Client
}

func (client *client) EnableDisasterRecovery(ctx context.Context, environmentId string) error {
	return client.setDisasterRecoveryWithRetry(ctx, environmentId, "Enabled", 0)
}

func (client *client) DisableDisasterRecovery(ctx context.Context, environmentId string) error {
	return client.setDisasterRecoveryWithRetry(ctx, environmentId, "Disabled", 0)
}

func (client *client) setDisasterRecoveryWithRetry(ctx context.Context, environmentId, state string, retryCount int) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.BAP_2021_API_VERSION)
	apiUrl.RawQuery = values.Encode()

	body := disasterRecoveryPatchDto{
		Properties: disasterRecoveryPatchPropertiesDto{
			States: disasterRecoveryPatchStatesDto{
				DisasterRecovery: environment.DisasterRecoveryStateDto{
					Id: state,
				},
			},
		},
	}

	tflog.Debug(ctx, fmt.Sprintf("Setting Disaster Recovery to '%s' for environment '%s'", state, environmentId))

	apiResponse, err := client.Api.Execute(ctx, nil, "PATCH", apiUrl.String(), nil, body, []int{http.StatusAccepted, http.StatusConflict}, nil)
	if err != nil {
		return err
	}

	tflog.Debug(ctx, fmt.Sprintf("Disaster Recovery PATCH HTTP Status: '%s'", apiResponse.HttpResponse.Status))

	tflog.Debug(ctx, "Waiting for Disaster Recovery lifecycle operation to complete")
	lifecycleResponse, err := client.Api.DoWaitForLifecycleOperationStatus(ctx, apiResponse)
	if err != nil {
		return err
	}

	if lifecycleResponse != nil && lifecycleResponse.State.Id == "Failed" {
		if retryCount >= constants.MAX_RETRY_COUNT {
			return fmt.Errorf("maximum retries (%d) reached for setting disaster recovery to '%s'", constants.MAX_RETRY_COUNT, state)
		}
		if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return err
		}
		tflog.Info(ctx, fmt.Sprintf("Disaster Recovery operation failed. Retrying (%d/%d)...", retryCount+1, constants.MAX_RETRY_COUNT))
		return client.setDisasterRecoveryWithRetry(ctx, environmentId, state, retryCount+1)
	}

	return nil
}

func (client *client) GetDisasterRecoveryState(ctx context.Context, environmentId string) (*environment.EnvironmentDto, error) {
	return client.environmentClient.GetEnvironment(ctx, environmentId)
}
