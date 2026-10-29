// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package enterprise_policy

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

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

// buildEnterprisePolicyURL builds the URL for enterprise policy operations.
func (client *Client) buildEnterprisePolicyURL(environmentId, environmentType, action string) string {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/environments/%s/enterprisePolicies/%s/%s", environmentId, environmentType, action),
	}

	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.ENTERPRISE_POLICY_API_VERSION)
	apiUrl.RawQuery = values.Encode()

	return apiUrl.String()
}

// executePolicyOperation executes a policy operation (link/unlink) with common retry logic.
func (client *Client) executePolicyOperation(ctx context.Context, environmentId, environmentType, systemId, action string) error {
	return client.executePolicyOperationWithRetry(ctx, environmentId, environmentType, systemId, action, 0)
}

func (client *Client) executePolicyOperationWithRetry(ctx context.Context, environmentId, environmentType, systemId, action string, retryCount int) error {
	apiUrl := client.buildEnterprisePolicyURL(environmentId, environmentType, action)

	linkEnterprosePolicyDto := linkEnterprosePolicyDto{
		SystemId: systemId,
	}

	apiResponse, err := client.Api.Execute(ctx, nil, "POST", apiUrl, nil, linkEnterprosePolicyDto, []int{http.StatusAccepted, http.StatusConflict}, nil)
	if err != nil {
		return err
	}

	tflog.Debug(ctx, fmt.Sprintf("Policy %s Operation HTTP Status: '%s'", action, apiResponse.HttpResponse.Status))

	if apiResponse.HttpResponse.StatusCode == http.StatusConflict {
		env, envErr := client.EnvironmentClient.GetEnvironment(ctx, environmentId)
		if envErr != nil {
			return envErr
		}
		if enterprisePolicyOperationDesiredStateExists(env, environmentType, systemId, action) {
			tflog.Debug(ctx, fmt.Sprintf("Policy %s desired state already exists, nothing to do", action))
			return nil
		}
		if retryCount >= constants.MAX_RETRY_COUNT {
			return fmt.Errorf("maximum retries (%d) reached for enterprise policy %s: the operation kept returning 409 and the requested state was not established", constants.MAX_RETRY_COUNT, action)
		}
		if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return err
		}
		tflog.Info(ctx, fmt.Sprintf("Policy %s operation was rejected with 409 and the requested state is not established. Retrying...", action))
		return client.executePolicyOperationWithRetry(ctx, environmentId, environmentType, systemId, action, retryCount+1)
	}

	tflog.Debug(ctx, "Waiting for operation to complete")

	lifecycleResponse, err := client.Api.DoWaitForLifecycleOperationStatus(ctx, apiResponse)
	if err != nil {
		return err
	}
	if lifecycleResponse != nil && lifecycleResponse.State.Id == "Failed" {
		if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return err
		}
		tflog.Info(ctx, fmt.Sprintf("Policy %s Operation failed. Retrying...", action))
		return client.executePolicyOperation(ctx, environmentId, environmentType, systemId, action)
	}
	return nil
}

func enterprisePolicyOperationDesiredStateExists(env *environment.EnvironmentDto, environmentType, systemId, action string) bool {
	var policy *environment.EnterprisePolicyDto
	switch environmentType {
	case NETWORK_INJECTION_POLICY_TYPE:
		if env.Properties.EnterprisePolicies != nil {
			policy = env.Properties.EnterprisePolicies.Vnets
		}
	case ENCRYPTION_POLICY_TYPE:
		if env.Properties.EnterprisePolicies != nil {
			policy = env.Properties.EnterprisePolicies.CustomerManagedKeys
		}
	case IDENTITY_POLICY_TYPE:
		if env.Properties.EnterprisePolicies != nil {
			policy = env.Properties.EnterprisePolicies.Identity
		}
	default:
		return false
	}

	requestedPolicyIsPresent := policy != nil && strings.EqualFold(policy.SystemId, systemId)
	switch action {
	case "link":
		return requestedPolicyIsPresent && strings.EqualFold(policy.LinkStatus, "Linked")
	case "unlink":
		return !requestedPolicyIsPresent || strings.EqualFold(policy.LinkStatus, "Unlinked")
	default:
		return false
	}
}

func (client *Client) LinkEnterprisePolicy(ctx context.Context, environmentId, environmentType, systemId string) error {
	return client.executePolicyOperation(ctx, environmentId, environmentType, systemId, "link")
}

func (client *Client) UnLinkEnterprisePolicy(ctx context.Context, environmentId, environmentType, systemId string) error {
	return client.executePolicyOperation(ctx, environmentId, environmentType, systemId, "unlink")
}
