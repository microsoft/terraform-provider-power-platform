// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package licensing

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
)

func NewLicensingClient(apiClient *api.Client) Client {
	return Client{
		Api: apiClient,
	}
}

type Client struct {
	Api *api.Client
}

func (client *Client) GetBillingPolicies(ctx context.Context) ([]BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   "licensing/billingPolicies",
	}

	values := url.Values{}
	values.Add("api-version", "2022-03-01-preview")
	apiUrl.RawQuery = values.Encode()

	policies := BillingPolicyArrayDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &policies)

	return policies.Value, err
}

func (client *Client) GetBillingPolicy(ctx context.Context, billingId string) (*BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/licensing/billingPolicies/%s", billingId),
	}

	values := url.Values{}
	values.Add("api-version", "2022-03-01-preview")
	apiUrl.RawQuery = values.Encode()

	policy := BillingPolicyDto{}
	resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &policy)

	if resp != nil && resp.HttpResponse.StatusCode == http.StatusNotFound {
		return nil, customerrors.WrapIntoProviderError(err, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("Billing Policy with ID '%s' not found", billingId))
	}
	return &policy, err
}

func (client *Client) CreateBillingPolicy(ctx context.Context, policyToCreate billingPolicyCreateDto) (*BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   "/licensing/BillingPolicies",
	}

	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, "2022-03-01-preview")
	apiUrl.RawQuery = values.Encode()

	policy := &BillingPolicyDto{}
	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, policyToCreate, []int{http.StatusCreated}, policy)
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

func (client *Client) UpdateBillingPolicy(ctx context.Context, billingId string, policyToUpdate BillingPolicyUpdateDto) (*BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/licensing/billingPolicies/%s", billingId),
	}

	values := url.Values{}
	values.Add("api-version", "2022-03-01-preview")
	apiUrl.RawQuery = values.Encode()

	policy := &BillingPolicyDto{}
	_, err := client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, policyToUpdate, []int{http.StatusOK}, policy)

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
	values.Add("api-version", "2022-03-01-preview")
	apiUrl.RawQuery = values.Encode()

	_, err := client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusNoContent}, nil)

	return err
}

func (client *Client) GetEnvironmentsForBillingPolicy(ctx context.Context, billingId string) ([]string, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("licensing/billingPolicies/%s/environments", billingId),
	}

	values := url.Values{}
	values.Add("api-version", "2022-03-01-preview")
	apiUrl.RawQuery = values.Encode()

	billingPolicyEnvironments := BillingPolicyEnvironmentsArrayResponseDto{}
	resp, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &billingPolicyEnvironments)
	if err != nil {
		if resp != nil && resp.HttpResponse.StatusCode == http.StatusNotFound {
			return nil, customerrors.WrapIntoProviderError(err, customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND), fmt.Sprintf("Billing Policy with ID '%s' not found", billingId))
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
	values.Add("api-version", "2022-03-01-preview")
	apiUrl.RawQuery = values.Encode()

	environments := BillingPolicyEnvironmentsArrayDto{
		EnvironmentIds: environmentIds,
	}
	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, environments, []int{http.StatusOK}, nil)

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
	values.Add("api-version", "2022-03-01-preview")
	apiUrl.RawQuery = values.Encode()

	environments := BillingPolicyEnvironmentsArrayDto{
		EnvironmentIds: environmentIds,
	}
	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, environments, []int{http.StatusOK}, nil)

	return err
}

func (client *Client) DoWaitForFinalStatus(ctx context.Context, billingPolicyDto *BillingPolicyDto) (*BillingPolicyDto, error) {
	billingId := billingPolicyDto.Id

	for {
		billingPolicy, err := client.GetBillingPolicy(ctx, billingId)

		if err != nil {
			return nil, err
		}

		if billingPolicy.Status == "Enabled" || billingPolicy.Status == "Disabled" {
			return billingPolicy, nil
		}

		if err := client.Api.SleepWithContext(ctx, api.DefaultRetryAfter()); err != nil {
			return nil, err
		}

		tflog.Debug(ctx, fmt.Sprintf("Billing Policy Operation State: '%s'", billingPolicy.Status))
	}
}
