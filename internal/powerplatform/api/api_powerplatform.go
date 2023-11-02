package powerplatform

import (
	"context"

	common "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/common"
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

func (client *PowerPlatformClientApi) Execute(ctx context.Context, method string, url string, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*ApiHttpResponse, error) {
	token, err := client.baseApi.InitializeBase(ctx, client.Auth)
	if err != nil {
		return nil, err
	}
	return client.baseApi.ExecuteBase(ctx, token, method, url, body, acceptableStatusCodes, responseObj)
}
