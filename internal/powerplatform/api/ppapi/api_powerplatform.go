package powerplatform

import (
	"context"
	"net/http"
	"net/url"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/models"
)

var _ api.PowerPlatformClientApiInterface = &PowerPlatformClientApi{}

type PowerPlatformClientApi struct {
	BaseApi api.ApiClientInterface
	Auth    PowerPlatformAuthInterface
}

func (client *PowerPlatformClientApi) GetBase() api.ApiClientInterface {
	return client.BaseApi
}

func (client *PowerPlatformClientApi) Execute(ctx context.Context, method string, url string, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*api.ApiHttpResponse, error) {
	token, err := client.GetBase().InitializeBase(ctx)
	if err != nil {
		return nil, err
	}
	return client.GetBase().ExecuteBase(ctx, token, method, url, body, acceptableStatusCodes, responseObj)
}

func (client *PowerPlatformClientApi) GetBillingPolicies(ctx context.Context) ([]models.BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.GetBase().GetConfig().Urls.PowerPlatformUrl,
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
