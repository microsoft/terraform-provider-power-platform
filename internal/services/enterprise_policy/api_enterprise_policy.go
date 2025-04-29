// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package enterprise_policy

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

func newEnterprisePolicyClient(apiClient *api.Client) Client {
	return Client{
		Api:               apiClient,
		EnvironmentClient: environment.NewEnvironmentClient(apiClient),
	}
}

type Client struct {
	Api               *api.Client
	EnvironmentClient environment.Client
}

func (client *Client) LinkEnterprisePolicy(ctx context.Context, environmentId, environmentType, systemId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/environments/%s/enterprisePolicies/%s/link", environmentId, environmentType),
	}

	values := url.Values{}
	values.Add("api-version", "2019-10-01")
	apiUrl.RawQuery = values.Encode()

	linkEnterprosePolicyDto := linkEnterprosePolicyDto{
		SystemId: systemId,
	}

	apiResponse, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, linkEnterprosePolicyDto, []int{http.StatusAccepted, http.StatusConflict}, nil)
	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Policy Linking Operation HTTP Status: '"+apiResponse.HttpResponse.Status+"'")
	tflog.Debug(ctx, "Waiting for operation to complete")

	lifecycleResponse, err := client.Api.DoWaitForLifecycleOperationStatus(ctx, apiResponse)
	if err != nil {
		return err
	}
	if lifecycleResponse != nil && lifecycleResponse.State.Id == "Failed" {
		if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return err
		}
		tflog.Info(ctx, "Policy Linking Operation failed. Retrying...")
		return client.LinkEnterprisePolicy(ctx, environmentId, environmentType, systemId)
	}
	return nil
}

func (client *Client) UnLinkEnterprisePolicy(ctx context.Context, environmentId, environmentType, systemId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/environments/%s/enterprisePolicies/%s/unlink", environmentId, environmentType),
	}

	values := url.Values{}
	values.Add("api-version", "2019-10-01")
	apiUrl.RawQuery = values.Encode()

	linkEnterprosePolicyDto := linkEnterprosePolicyDto{
		SystemId: systemId,
	}
	apiResponse, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, linkEnterprosePolicyDto, []int{http.StatusAccepted, http.StatusConflict}, nil)
	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Policy Unlinking Operation HTTP Status: '"+apiResponse.HttpResponse.Status+"'")
	tflog.Debug(ctx, "Waiting for operation to complete")

	lifecycleResponse, err := client.Api.DoWaitForLifecycleOperationStatus(ctx, apiResponse)
	if err != nil {
		return err
	}
	if lifecycleResponse != nil && lifecycleResponse.State.Id == "Failed" {
		if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return err
		}
		tflog.Info(ctx, "Policy Unlinking Operation failed. Retrying...")
		return client.UnLinkEnterprisePolicy(ctx, environmentId, environmentType, systemId)
	}
	return nil
}
