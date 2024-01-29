package powerplatform

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	config "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/config"
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

func (client *DataverseClientApi) GetConfig() *config.ProviderConfig {
	return client.baseApi.Config
}

func (client *DataverseClientApi) SetBapiClient(bapiClient *BapiClientApi) {
	client.bapiClient = bapiClient
}

func (client *DataverseClientApi) Initialize(ctx context.Context, environmentUrl string) (string, error) {
	environmentUrl = strings.TrimSuffix(environmentUrl, "/")
	scopes := []string{environmentUrl + "//.default"}

	token, err := client.auth.GetToken(scopes)

	if _, ok := err.(*TokeExpiredError); ok {
		tflog.Debug(ctx, "Token expired. authenticating...")

		if client.baseApi.GetConfig().Credentials.IsClientSecretCredentialsProvided() {
			token, err := client.auth.AuthenticateClientSecret(ctx, scopes, client.baseApi.GetConfig().Credentials)
			if err != nil {
				return "", err
			}
			tflog.Info(ctx, fmt.Sprintln("Dataverse token aquired: ", "********"))
			return token, nil
		} else if client.baseApi.GetConfig().Credentials.IsUserPassCredentialsProvided() {
			token, err := client.auth.AuthenticateUserPass(ctx, scopes, client.baseApi.GetConfig().Credentials)
			if err != nil {
				return "", err
			}
			tflog.Info(ctx, fmt.Sprintln("Dataverse token aquired: ", "********"))
			return token, nil
		} else if client.baseApi.GetConfig().Credentials.UseCli {
			token, err := client.auth.AuthUsingCli(ctx, scopes, client.baseApi.GetConfig().Credentials)
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

func (client *DataverseClientApi) Execute(ctx context.Context, environmentUrl, method string, url string, headers http.Header, body interface{}, acceptableStatusCodes []int, responseObj interface{}) (*ApiHttpResponse, error) {
	token, err := client.Initialize(ctx, environmentUrl)
	if err != nil {
		return nil, err
	}
	return client.baseApi.ExecuteBase(ctx, token, method, url, headers, body, acceptableStatusCodes, responseObj)
}
