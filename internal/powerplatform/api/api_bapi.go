package powerplatform

import (
	"context"

	common "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/common"
)

type BapiClientApi struct {
	baseApi         *ApiClientBase
	auth            *BapiAuth
	dataverseClient *DataverseClientApi
}

func NewBapiClientApi(baseApi *ApiClientBase, auth *BapiAuth, dataverseClient *DataverseClientApi) *BapiClientApi {
	return &BapiClientApi{
		baseApi:         baseApi,
		auth:            auth,
		dataverseClient: dataverseClient,
	}
}

func (client *BapiClientApi) SetDataverseClient(dataverseClient *DataverseClientApi) {
	client.dataverseClient = dataverseClient
}

func (client *BapiClientApi) GetConfig() *common.ProviderConfig {
	return client.baseApi.Config
}

func (client *BapiClientApi) Execute(ctx context.Context, method string, url string, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*ApiHttpResponse, error) {
	token, err := client.baseApi.InitializeBase(ctx, client.auth)
	if err != nil {
		return nil, err
	}
	return client.baseApi.ExecuteBase(ctx, token, method, url, body, acceptableStatusCodes, responseObj)
}
