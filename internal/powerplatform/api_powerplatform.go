package powerplatform

import (
	"context"
	"net/http"
	"net/url"

	common "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/common"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/models"
)

var _ PowerPlatformClientInterface = &PowerPlatformClientImplementation{}

type PowerPlatformClientInterface interface {
	GetBase() common.ApiClientInterface

	GetBillingPolicies(ctx context.Context) ([]models.BillingPolicyDto, error)
}

type PowerPlatformClientImplementation struct {
	BaseApi common.ApiClientInterface
	Auth    PowerPlatformAuthInterface
}

func (client *PowerPlatformClientImplementation) GetBase() common.ApiClientInterface {
	return client.BaseApi
}

func (client *PowerPlatformClientImplementation) doRequest(ctx context.Context, request *http.Request) (*common.ApiHttpResponse, error) {
	token, err := client.BaseApi.Initialize(ctx)
	if err != nil {
		return nil, err
	}

	return client.BaseApi.DoRequest(token, request)
}

func (client *PowerPlatformClientImplementation) GetBillingPolicies(ctx context.Context) ([]models.BillingPolicyDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BaseApi.GetConfig().Urls.PowerPlatformUrl,
		Path:   "/licensing/billingPolicies",
	}
	values := url.Values{}
	values.Add("api-version", "2022-03-01-preview")
	apiUrl.RawQuery = values.Encode()
	request, err := http.NewRequestWithContext(ctx, "GET", apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	apiResponse, err := client.doRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	billingPolicies := models.BillingPolicyDtoArray{}
	err = apiResponse.MarshallTo(&billingPolicies)
	if err != nil {
		return nil, err
	}

	return billingPolicies.Value, nil
}
