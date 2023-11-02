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
		ppLicesingApi: ppLicensingApi,
	}
}

func (client *LicensingClient) GetBillingPolicy(ctx context.Context, billingId string) (*BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.ppLicesingApi.GetConfig().Urls.PowerPlatformLicensingUrl,
		Path:   fmt.Sprintf("/licensing/billingPolicies/%s", billingId),
	}

	policy := BillingPolicyDto{}
	_, err := client.ppLicesingApi.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &policy)

	return &policy, err
}

func (client *LicensingClient) GetBillingPolicies(ctx context.Context) ([]BillingPolicyDto, error) {
	return nil, nil
}

func (client *LicensingClient) CreateBillingPolicy(ctx context.Context, policyToCreate BillingPolicyCreateDto) (*BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.ppLicesingApi.GetConfig().Urls.PowerPlatformLicensingUrl,
		Path:   fmt.Sprintf("/v1.0/tenants/%s/BillingPolicies", client.ppLicesingApi.GetConfig().Credentials.TenantId),
	}

	policy := BillingPolicyDto{}
	_, err := client.ppLicesingApi.Execute(ctx, "POST", apiUrl.String(), policyToCreate, []int{http.StatusOK}, nil)

	return &policy, err
}

func (client *LicensingClient) UpdateBillingPolicy(ctx context.Context, policyToUpdate BillingPolicyDto) (*BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.ppLicesingApi.GetConfig().Urls.PowerPlatformLicensingUrl,
		Path:   fmt.Sprintf("/v1.0/tenants/%s/BillingPolicies", client.ppLicesingApi.GetConfig().Credentials.TenantId),
	}

	policy := BillingPolicyDto{}
	_, err := client.ppLicesingApi.Execute(ctx, "PUT", apiUrl.String(), policyToUpdate, []int{http.StatusOK}, &policy)

	return &policy, err
}

func (client *LicensingClient) DeleteBillingPolicy(ctx context.Context, id string) error {
	return nil
}
