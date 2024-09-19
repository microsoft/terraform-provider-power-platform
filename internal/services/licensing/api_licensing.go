// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package licensing

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type Client struct {
	Api *api.Client
}

func NewLicensingClient(apiClient *api.Client) Client {
	return Client{
		Api: apiClient,
	}
}

const (
	API_VERSION_2022_03_01_preview = "2022-03-01-preview"
)

func (client *Client) GetBillingPolicies(ctx context.Context) ([]billingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   "licensing/billingPolicies",
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION_2022_03_01_preview)
	apiUrl.RawQuery = values.Encode()

	policies := billingPolicyArrayDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &policies)

	return policies.Value, err
}

func (client *Client) GetBillingPolicy(ctx context.Context, billingId string) (*billingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/licensing/billingPolicies/%s", billingId),
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION_2022_03_01_preview)
	apiUrl.RawQuery = values.Encode()

	policy := billingPolicyDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &policy)

	if err != nil && strings.ContainsAny(err.Error(), "404") {
		return nil, helpers.WrapIntoProviderError(err, helpers.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("Billing Policy with ID '%s' not found", billingId))
	}
	return &policy, err
}

func (client *Client) CreateBillingPolicy(ctx context.Context, policyToCreate billingPolicyCreateDto) (*billingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   "/licensing/BillingPolicies",
	}

	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, API_VERSION_2022_03_01_preview)
	apiUrl.RawQuery = values.Encode()

	policy := &billingPolicyDto{}
	_, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, policyToCreate, []int{http.StatusCreated}, policy)
	if err != nil {
		return nil, err
	}

	// If billing policy status is not Enabled or Disabled, wait for it to reach a terminal state
	if policy.Status != "Enabled" && policy.Status != "Disabled" {
		policy, err = client.DoWaitForFinalStatus(ctx, policy)

		if err != nil {
			return nil, err
		}
	}

	return policy, err
}

func (client *Client) UpdateBillingPolicy(ctx context.Context, billingId string, policyToUpdate billingPolicyUpdateDto) (*billingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/licensing/billingPolicies/%s", billingId),
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION_2022_03_01_preview)
	apiUrl.RawQuery = values.Encode()

	policy := &billingPolicyDto{}
	_, err := client.Api.Execute(ctx, "PUT", apiUrl.String(), nil, policyToUpdate, []int{http.StatusOK}, policy)

	// If billing policy status is not Enabled or Disabled, wait for it to reach a terminal state
	if policy.Status != "Enabled" && policy.Status != "Disabled" {
		policy, err = client.DoWaitForFinalStatus(ctx, policy)

		if err != nil {
			return nil, err
		}
	}

	return policy, err
}

func (client *Client) DeleteBillingPolicy(ctx context.Context, billingId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/licensing/BillingPolicies/%s", billingId),
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION_2022_03_01_preview)
	apiUrl.RawQuery = values.Encode()

	_, err := client.Api.Execute(ctx, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusNoContent}, nil)

	return err
}

func (client *Client) GetEnvironmentsForBillingPolicy(ctx context.Context, billingId string) ([]string, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("licensing/billingPolicies/%s/environments", billingId),
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION_2022_03_01_preview)
	apiUrl.RawQuery = values.Encode()

	billingPolicyEnvironments := billingPolicyEnvironmentsArrayResponseDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &billingPolicyEnvironments)
	if err != nil {
		if strings.ContainsAny(err.Error(), "404") {
			return nil, helpers.WrapIntoProviderError(err, helpers.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("Billing Policy with ID '%s' not found", billingId))
		}
		return nil, err
	}

	environments := []string{}
	for _, billingPolicyEnvironment := range billingPolicyEnvironments.Value {
		environments = append(environments, billingPolicyEnvironment.EnvironmentId)
	}
	return environments, err
}

func (client *Client) AddEnvironmentsToBillingPolicy(ctx context.Context, billingId string, environmentIds []string) error {
	if len(environmentIds) == 0 {
		return nil
	}
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("licensing/billingPolicies/%s/environments/add", billingId),
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION_2022_03_01_preview)
	apiUrl.RawQuery = values.Encode()

	environments := billingPolicyEnvironmentsArrayDto{
		EnvironmentIds: environmentIds,
	}
	_, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, environments, []int{http.StatusOK}, nil)

	return err
}

func (client *Client) RemoveEnvironmentsToBillingPolicy(ctx context.Context, billingId string, environmentIds []string) error {
	if len(environmentIds) == 0 {
		return nil
	}
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("licensing/billingPolicies/%s/environments/remove", billingId),
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION_2022_03_01_preview)
	apiUrl.RawQuery = values.Encode()

	environments := billingPolicyEnvironmentsArrayDto{
		EnvironmentIds: environmentIds,
	}
	_, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, environments, []int{http.StatusOK}, nil)

	return err
}

func (client *Client) DoWaitForFinalStatus(ctx context.Context, billingPolicyDto *billingPolicyDto) (*billingPolicyDto, error) {
	billingId := billingPolicyDto.Id

	for {
		billingPolicy, err := client.GetBillingPolicy(ctx, billingId)

		if err != nil {
			return nil, err
		}

		if billingPolicy.Status == "Enabled" || billingPolicy.Status == "Disabled" {
			return billingPolicy, nil
		}

		err = client.Api.SleepWithContext(ctx, client.Api.RetryAfterDefault())
		if err != nil {
			return nil, err
		}

		tflog.Debug(ctx, fmt.Sprintf("Billing Policy Operation State: '%s'", billingPolicy.Status))
	}
}
