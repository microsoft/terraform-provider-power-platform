package powerplatform

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

type LicensingClient struct {
	ppApi *api.PowerPlatformClientApi
}

func NewLicensingClient(ppApi *api.PowerPlatformClientApi) LicensingClient {
	return LicensingClient{
		ppApi: ppApi,
	}
}

const (
	API_VERSION = "2022-03-01-preview"
)

func (client *LicensingClient) GetBillingPolicies(ctx context.Context) ([]BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.ppApi.GetConfig().Urls.PowerPlatformUrl,
		Path:   "licensing/billingPolicies",
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION)
	apiUrl.RawQuery = values.Encode()

	policies := BillingPolicyArrayDto{}
	_, err := client.ppApi.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &policies)

	return policies.Value, err
}

func (client *LicensingClient) GetBillingPolicy(ctx context.Context, billingId string) (*BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.ppApi.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/licensing/billingPolicies/%s", billingId),
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION)
	apiUrl.RawQuery = values.Encode()

	policy := BillingPolicyDto{}
	_, err := client.ppApi.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &policy)

	return &policy, err
}

func (client *LicensingClient) CreateBillingPolicy(ctx context.Context, policyToCreate BillingPolicyCreateDto) (*BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.ppApi.GetConfig().Urls.PowerPlatformUrl,
		Path:   "/licensing/BillingPolicies",
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION)
	apiUrl.RawQuery = values.Encode()

	policy := BillingPolicyDto{}
	_, err := client.ppApi.Execute(ctx, "POST", apiUrl.String(), nil, policyToCreate, []int{http.StatusCreated}, &policy)

	return &policy, err
}

func (client *LicensingClient) UpdateBillingPolicy(ctx context.Context, billingId string, policyToUpdate BillingPolicyUpdateDto) (*BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.ppApi.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/licensing/billingPolicies/%s", billingId),
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION)
	apiUrl.RawQuery = values.Encode()

	policy := BillingPolicyDto{}
	_, err := client.ppApi.Execute(ctx, "PUT", apiUrl.String(), nil, policyToUpdate, []int{http.StatusOK}, &policy)

	return &policy, err
}

func (client *LicensingClient) DeleteBillingPolicy(ctx context.Context, billingId string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.ppApi.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("/licensing/BillingPolicies/%s", billingId),
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION)
	apiUrl.RawQuery = values.Encode()

	_, err := client.ppApi.Execute(ctx, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusNoContent}, nil)

	return err
}

func (client *LicensingClient) AddEnvironmentsToBillingPolicy(ctx context.Context, billingId string, environmentIds []string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.ppApi.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("licensing/billingPolicies/%s/environments/add", billingId),
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION)
	apiUrl.RawQuery = values.Encode()

	environments := BillingPolicyEnvironmentsArrayDto{
		EnvironmentIds: environmentIds,
	}
	_, err := client.ppApi.Execute(ctx, "POST", apiUrl.String(), nil, environments, []int{http.StatusOK}, nil)

	return err
}

func (client *LicensingClient) RemoveEnvironmentsToBillingPolicy(ctx context.Context, billingId string, environmentIds []string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.ppApi.GetConfig().Urls.PowerPlatformUrl,
		Path:   fmt.Sprintf("licensing/billingPolicies/%s/environments/add", billingId),
	}

	values := url.Values{}
	values.Add("api-version", API_VERSION)
	apiUrl.RawQuery = values.Encode()

	environments := BillingPolicyEnvironmentsArrayDto{
		EnvironmentIds: environmentIds,
	}
	_, err := client.ppApi.Execute(ctx, "POST", apiUrl.String(), nil, environments, []int{http.StatusOK}, nil)

	return err
}
