package powerplatform

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/models"
)

// this line is here to make sure the interface is implemented
var _ PowerPlatformClientApiInterface = &PowerPlatformClientApi{}

type PowerPlatformClientApiInterface interface {
	GetBase() api.ApiClientInterface
	Execute(ctx context.Context, method string, url string, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*api.ApiHttpResponse, error)

	GetBillingPolicies(ctx context.Context) ([]models.BillingPolicyDto, error)
	GetBillingPolicy(ctx context.Context, id string) (*models.BillingPolicyDto, error)
	CreateBillingPolicy(ctx context.Context, policyToCreate models.BillingPolicyDto) (*models.BillingPolicyDto, error)
}

type PowerPlatformClientApi struct {
	BaseApi api.ApiClientInterface
	Auth    PowerPlatformAuthInterface
}

func (client *PowerPlatformClientApi) GetBase() api.ApiClientInterface {
	return client.BaseApi
}

func (client *PowerPlatformClientApi) Execute(ctx context.Context, method string, url string, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*api.ApiHttpResponse, error) {
	token, err := client.BaseApi.InitializeBase(ctx)
	if err != nil {
		return nil, err
	}
	return client.BaseApi.ExecuteBase(ctx, token, method, url, body, acceptableStatusCodes, responseObj)
}

func (client *PowerPlatformClientApi) GetBillingPolicies(ctx context.Context) ([]models.BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BaseApi.GetConfig().Urls.PowerPlatformUrl,
		Path:   "/licensing/billingPolicies",
	}
	values := url.Values{}
	values.Add("api-version", "2022-03-01-preview")
	apiUrl.RawQuery = values.Encode()

	billingPolicies := models.BillingPolicyDtoArray{}
	_, err := client.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &billingPolicies)
	if err != nil {
		return nil, err
	}

	return billingPolicies.Value, nil
}

func (client *PowerPlatformClientApi) GetBillingPolicy(ctx context.Context, id string) (*models.BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme:   "https",
		Host:     "api.powerplatform.com",
		Path:     fmt.Sprintf("/licensing/billingPolicies/%s", id),
		RawQuery: url.Values{"api-version": {"2022-03-01-preview"}}.Encode(),
	}

	policy := models.BillingPolicyDto{}

	_, err := client.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &policy)
	return &policy, err

}

func (client *PowerPlatformClientApi) CreateBillingPolicy(ctx context.Context, policyToCreate BillingPolicyDto) (*BillingPolicyDto, error) {
	return nil, nil
}
