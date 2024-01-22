package powerplatform

import (
	"context"
	"net/http"

	config "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/config"
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

func (client *PowerPlatformClientApi) GetConfig() *config.ProviderConfig {
	return client.baseApi.Config
}

func (client *PowerPlatformClientApi) Execute(ctx context.Context, method string, url string, headers http.Header, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*ApiHttpResponse, error) {
	token, err := client.baseApi.InitializeBase(ctx, []string{"https://api.powerplatform.com/.default"}, client.Auth)
	if err != nil {
		return nil, err
	}
	return client.baseApi.ExecuteBase(ctx, token, method, url, headers, body, acceptableStatusCodes, responseObj)
}
