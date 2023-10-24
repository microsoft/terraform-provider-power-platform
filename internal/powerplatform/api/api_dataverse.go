package powerplatform

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	common "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/common"
)

type DataverseClientApi struct {
	baseApi    *ApiClientBase
	auth       *DataverseAuth
	bapiClient *BapiClientApi
}

func NewDataverseClientApi(baseApi *ApiClientBase, auth *DataverseAuth) *DataverseClientApi {
	return &DataverseClientApi{
		baseApi: baseApi,
		auth:    auth,
	}
}

func (client *DataverseClientApi) GetConfig() *common.ProviderConfig {
	return client.baseApi.Config
}

func (client *DataverseClientApi) SetBapiClient(bapiClient *BapiClientApi) {
	client.bapiClient = bapiClient
}

func (client *DataverseClientApi) Initialize(ctx context.Context, environmentUrl string) (string, error) {

	token, err := client.auth.GetToken(environmentUrl)

	if _, ok := err.(*TokeExpiredError); ok {
		tflog.Debug(ctx, "Token expired. authenticating...")

		if client.baseApi.GetConfig().Credentials.IsClientSecretCredentialsProvided() {
			token, err := client.auth.AuthenticateClientSecret(ctx, environmentUrl, client.baseApi.GetConfig().Credentials.TenantId, client.baseApi.GetConfig().Credentials.ClientId, client.baseApi.GetConfig().Credentials.Secret)
			if err != nil {
				return "", err
			}
			tflog.Info(ctx, fmt.Sprintln("Dataverse token aquired: ", "********"))
			return token, nil
		} else if client.baseApi.GetConfig().Credentials.IsUserPassCredentialsProvided() {
			token, err := client.auth.AuthenticateUserPass(ctx, environmentUrl, client.baseApi.GetConfig().Credentials.TenantId, client.baseApi.GetConfig().Credentials.Username, client.baseApi.GetConfig().Credentials.Password)
			if err != nil {
				return "", err
			}
			tflog.Info(ctx, fmt.Sprintln("Dataverse token aquired: ", "********"))
			return token, nil
		} else {
			return "", errors.New("no credentials provided")
		}

	} else if err != nil {
		return "", err
	} else {
		tflog.Info(ctx, fmt.Sprintln("Dataverse token aquired: ", "********"))
		return token, nil
	}
}

func (client *DataverseClientApi) Execute(ctx context.Context, environmentUrl, method string, url string, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*ApiHttpResponse, error) {
	token, err := client.Initialize(ctx, environmentUrl)
	if err != nil {
		return nil, err
	}
	return client.baseApi.ExecuteBase(ctx, token, method, url, body, acceptableStatusCodes, responseObj)
}
