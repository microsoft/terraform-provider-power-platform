package powerplatform

import (
	"context"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

// this line is here to make sure the interface is implemented
var _ PowerPlatformClientApiInterface = &PowerPlatformClientApi{}

type PowerPlatformClientApiInterface interface {
	GetBase() api.ApiClientInterface
	Execute(ctx context.Context, method string, url string, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*api.ApiHttpResponse, error)
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
