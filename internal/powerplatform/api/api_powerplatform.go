package powerplatform

import (
	"context"
	"net/http"
	"net/url"

	common "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/common"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/models"
)

type PowerPlatformClientApi struct {
	baseApi *ApiClientBase
	Auth    *PowerPlatformAuth
}

func NewPowerPlatformClientApi(baseApi *ApiClientBase, auth *PowerPlatformAuth) *PowerPlatformClientApi {
	return &PowerPlatformClientApi{
		baseApi: baseApi,
		Auth:    auth,
	}
}

func (client *PowerPlatformClientApi) GetConfig() *common.ProviderConfig {
	return client.baseApi.Config
}

func (client *PowerPlatformClientApi) Execute(ctx context.Context, method string, url string, headers http.Header, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*ApiHttpResponse, error) {
	token, err := client.baseApi.InitializeBaseWithUserNamePassword(ctx, client.Auth)
	if err != nil {
		return nil, err
	}
	return client.baseApi.ExecuteBase(ctx, token, method, url, headers, body, acceptableStatusCodes, responseObj)
}

func (client *PowerPlatformClientApi) GetBillingPolicies(ctx context.Context) ([]models.BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.baseApi.GetConfig().Urls.PowerPlatformUrl,
		Path:   "/licensing/billingPolicies",
	}
	values := url.Values{}
	values.Add("api-version", "2022-03-01-preview")
	apiUrl.RawQuery = values.Encode()

	billingPolicies := models.BillingPolicyDtoArray{}
	_, err := client.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &billingPolicies)
	if err != nil {
		return nil, err
	}

	return billingPolicies.Value, nil
}
