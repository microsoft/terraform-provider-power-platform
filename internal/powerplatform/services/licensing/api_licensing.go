// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

type LicensingClient struct {
	Api *api.ApiClient
}

func NewLicensingClient(api *api.ApiClient) LicensingClient {
	return LicensingClient{
		Api: api,
	}
}

const (
	API_VERSION = "2022-03-01-preview"
)

func (client *LicensingClient) GetBillingPolicies(ctx context.Context) ([]BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   "licensing/billingPolicies",
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION)
	apiUrl.RawQuery = values.Encode()

	policies := BillingPolicyArrayDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &policies)

	return policies.Value, err
}

func (client *LicensingClient) GetBillingPolicy(ctx context.Context, billingId string) (*BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/licensing/billingPolicies/%s", billingId),
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION)
	apiUrl.RawQuery = values.Encode()

	policy := BillingPolicyDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &policy)

	return &policy, err
}

func (client *LicensingClient) CreateBillingPolicy(ctx context.Context, policyToCreate BillingPolicyCreateDto) (*BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   "/licensing/BillingPolicies",
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION)
	apiUrl.RawQuery = values.Encode()

	policy := &BillingPolicyDto{}
	_, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, policyToCreate, []int{http.StatusCreated}, policy)

	// If billing policy status is not Enabled or Disabled, wait for it to reach a terminal state
	if policy.Status != "Enabled" && policy.Status != "Disabled" {
		policy, err = client.DoWaitForFinalStatus(ctx, policy)

		if err != nil {
			return nil, err
		}
	}

	return policy, err
}

func (client *LicensingClient) UpdateBillingPolicy(ctx context.Context, billingId string, policyToUpdate BillingPolicyUpdateDto) (*BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/licensing/billingPolicies/%s", billingId),
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION)
	apiUrl.RawQuery = values.Encode()

	policy := &BillingPolicyDto{}
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

func (client *LicensingClient) DeleteBillingPolicy(ctx context.Context, billingId string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/licensing/BillingPolicies/%s", billingId),
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION)
	apiUrl.RawQuery = values.Encode()

	_, err := client.Api.Execute(ctx, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusNoContent}, nil)

	return err
}

func (client *LicensingClient) GetEnvironmentsForBillingPolicy(ctx context.Context, billingId string) ([]string, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("licensing/billingPolicies/%s/environments", billingId),
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION)
	apiUrl.RawQuery = values.Encode()

	billingPolicyEnvironments := BillingPolicyEnvironmentsArrayResponseDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &billingPolicyEnvironments)
	if err != nil {
		return nil, err
	}

	environments := []string{}
	for _, billingPolicyEnvironment := range billingPolicyEnvironments.Value {
		environments = append(environments, billingPolicyEnvironment.EnvironmentId)
	}
	return environments, err
}

func (client *LicensingClient) AddEnvironmentsToBillingPolicy(ctx context.Context, billingId string, environmentIds []string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("licensing/billingPolicies/%s/environments/add", billingId),
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION)
	apiUrl.RawQuery = values.Encode()

	environments := BillingPolicyEnvironmentsArrayDto{
		EnvironmentIds: environmentIds,
	}
	_, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, environments, []int{http.StatusOK}, nil)

	return err
}

func (client *LicensingClient) RemoveEnvironmentsToBillingPolicy(ctx context.Context, billingId string, environmentIds []string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("licensing/billingPolicies/%s/environments/remove", billingId),
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION)
	apiUrl.RawQuery = values.Encode()

	environments := BillingPolicyEnvironmentsArrayDto{
		EnvironmentIds: environmentIds,
	}
	_, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, environments, []int{http.StatusOK}, nil)

	return err
}

func (client *LicensingClient) DoWaitForFinalStatus(ctx context.Context, billingPolicyDto *BillingPolicyDto) (*BillingPolicyDto, error) {
	billingId := billingPolicyDto.Id

	retryAfter := time.Duration(5) * time.Second

	timeout := time.Duration(3) * time.Minute

	startTime := time.Now()

	for {
		billingPolicy, err := client.GetBillingPolicy(ctx, billingId)

		if err != nil {
			return nil, err
		}

		if billingPolicy.Status == "Enabled" || billingPolicy.Status == "Disabled" {
			return billingPolicy, nil
		}

		if time.Since(startTime) >= timeout {
			tflog.Debug(ctx, "Timeout reached while waiting for billing policy to reach a terminal state (Enabled or Disabled)")
			err := fmt.Errorf("timeout reached while waiting for billing policy to reach a terminal state (Enabled or Disabled)")
			return nil, err
		}

		client.Api.Sleep(retryAfter)

		tflog.Debug(ctx, fmt.Sprintf("Billing Policy Operation State: '%s'", billingPolicy.Status))
	}
}
