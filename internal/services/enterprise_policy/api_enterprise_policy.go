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

// https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/1bd6bd98-e60e-e571-b6d6-2d249dbf3148/enterprisePolicies/Encryption/link?api-version=2019-10-01
// Request Method:
// POST
// Status Code:
// 202 Accepted
// {"systemId":"/regions/europe/providers/Microsoft.PowerPlatform/enterprisePolicies/aee7e0b0-c23a-4820-ae0e-ebe2a7162b4f"}
// ----
// Properties/SystemId:  "/regions/europe/providers/Microsoft.PowerPlatform/enterprisePolicies/aee7e0b0-c23a-4820-ae0e-ebe2a7162b4f",

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
	apiResponse, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, linkEnterprosePolicyDto, []int{http.StatusAccepted}, nil)
	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Policy Linking Operation HTTP Status: '"+apiResponse.HttpResponse.Status+"'")
	tflog.Debug(ctx, "Waiting for operation to complete")

	_, err = client.Api.DoWaitForLifecycleOperationStatus(ctx, apiResponse)
	if err != nil {
		return err
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
	apiResponse, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, linkEnterprosePolicyDto, []int{http.StatusAccepted}, nil)
	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Policy Unlinking Operation HTTP Status: '"+apiResponse.HttpResponse.Status+"'")
	tflog.Debug(ctx, "Waiting for operation to complete")

	_, err = client.Api.DoWaitForLifecycleOperationStatus(ctx, apiResponse)
	if err != nil {
		return err
	}
	return nil
}
