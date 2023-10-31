package licensing

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/clients"
)

type LicensingClientInterface interface {
	GetBillingPolicies(ctx context.Context) ([]BillingPolicyDto, error)
	GetBillingPolicy(ctx context.Context, id string) (*BillingPolicyDto, error)
	CreateBillingPolicy(ctx context.Context, policyToCreate BillingPolicyDto) (*BillingPolicyDto, error)
}

type LicensingClient struct {
	ppapi api.ApiClientInterface
	token string
}

const (
	API_VERSION = "2022-03-01-preview"
)

func NewLicensingClient(ppapi *clients.PowerPlatoformApiClient) LicensingClientInterface {
	return &LicensingClient{
		ppapi: ppapi,
	}
}

func (client *BillingPolicyClient) GetBillingPolicy(ctx context.Context, id string) (*BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   "api.powerplatform.com",
		Path:   fmt.Sprintf("/licensing/billingPolicies/%s", id),
	}

	policy := BillingPolicyDto{}
	_, err := client.ppapi.ExecuteBase(ctx, client.token, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &policy)

	return &policy, err
}

func (client *BillingPolicyClient) GetBillingPolicies(ctx context.Context) ([]BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   "api.powerplatform.com",
		Path:   "/licensing/billingPolicies",
	}

	policies := []BillingPolicyDto{}
	_, err := client.ppapi.ExecuteBase(ctx, client.token, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &policies)

	return policies, err
}

func (client *BillingPolicyClient) CreateBillingPolicy(ctx context.Context, policyToCreate BillingPolicyDto) (*BillingPolicyDto, error) {
	return nil, nil
}

func (client *BillingPolicyClient) UpdateBillingPolicy(ctx context.Context, policyToUpdate BillingPolicyDto) (*BillingPolicyDto, error) {
	return nil, nil
}

func (client *BillingPolicyClient) DeleteBillingPolicy(ctx context.Context, id string) error {
	return nil
}

// func (client *PowerPlatformClientApi) GetBillingPolicies(ctx context.Context) ([]models.BillingPolicyDto, error) {
// 	apiUrl := &url.URL{
// 		Scheme: "https",
// 		Host:   client.BaseApi.GetConfig().Urls.PowerPlatformUrl,
// 		Path:   "/licensing/billingPolicies",
// 	}
// 	values := url.Values{}
// 	values.Add("api-version", "2022-03-01-preview")
// 	apiUrl.RawQuery = values.Encode()

// 	billingPolicies := models.BillingPolicyDtoArray{}
// 	_, err := client.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &billingPolicies)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return billingPolicies.Value, nil
// }

// func (client *PowerPlatformClientApi) GetBillingPolicy(ctx context.Context, id string) (*models.BillingPolicyDto, error) {
// 	apiUrl := &url.URL{
// 		Scheme:   "https",
// 		Host:     "api.powerplatform.com",
// 		Path:     fmt.Sprintf("/licensing/billingPolicies/%s", id),
// 		RawQuery: url.Values{"api-version": {"2022-03-01-preview"}}.Encode(),
// 	}

// 	policy := models.BillingPolicyDto{}

// 	_, err := client.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &policy)
// 	return &policy, err

// }

// func (client *PowerPlatformClientApi) CreateBillingPolicy(ctx context.Context, policyToCreate BillingPolicyDto) (*BillingPolicyDto, error) {
// 	return nil, nil
// }
