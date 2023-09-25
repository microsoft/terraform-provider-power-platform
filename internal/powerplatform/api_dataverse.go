package powerplatform

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	powerplatform_bapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi/models"
)

var _ DataverseClientInterface = &DataverseClientImplementation{}

type DataverseClientInterface interface {
	Initialize(ctx context.Context, environmentUrl string) (string, error)

	GetSolutions(ctx context.Context, environmentId string) ([]models.SolutionDto, error)
}

type DataverseClientImplementation struct {
	BaseApi    ApiClientInterface
	Auth       DataverseAuthInterface
	BapiClient BapiClientInterface
}

func (client *DataverseClientImplementation) doRequest(ctx context.Context, environmentUrl string, request *http.Request) (*powerplatform_bapi.ApiHttpResponse, error) {
	token, err := client.Initialize(ctx, environmentUrl)
	if err != nil {
		return nil, err
	}
	return client.BaseApi.DoRequest(token, request)
}

func (client *DataverseClientImplementation) Initialize(ctx context.Context, environmentUrl string) (string, error) {

	token, err := client.Auth.GetToken(environmentUrl)

	if _, ok := err.(*TokeExpiredError); ok {
		tflog.Debug(ctx, "Token expired. authenticating...")

		if client.BaseApi.GetConfig().Credentials.IsClientSecretCredentialsProvided() {
			token, err := client.Auth.AuthenticateClientSecret(ctx, environmentUrl, client.BaseApi.GetConfig().Credentials.TenantId, client.BaseApi.GetConfig().Credentials.ClientId, client.BaseApi.GetConfig().Credentials.Secret)
			if err != nil {
				return "", err
			}
			return token, nil
		} else if client.BaseApi.GetConfig().Credentials.IsUserPassCredentialsProvided() {
			token, err := client.Auth.AuthenticateUserPass(ctx, environmentUrl, client.BaseApi.GetConfig().Credentials.TenantId, client.BaseApi.GetConfig().Credentials.Username, client.BaseApi.GetConfig().Credentials.Password)
			if err != nil {
				return "", err
			}
			return token, nil
		} else {
			return "", errors.New("no credentials provided")
		}

	} else if err != nil {
		return "", err
	} else {
		return token, nil
	}
}

func (client *DataverseClientImplementation) getEnvironmentUrlById(ctx context.Context, environmentId string) (string, error) {
	env, err := client.BapiClient.GetEnvironment(ctx, environmentId)
	if err != nil {
		return "", err
	}
	environmentUrl := strings.TrimSuffix(env.Properties.LinkedEnvironmentMetadata.InstanceURL, "/")
	return environmentUrl, nil
}

func (client *DataverseClientImplementation) GetSolutions(ctx context.Context, environmentId string) ([]models.SolutionDto, error) {
	environmentUrl, err := client.getEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   "/api/data/v9.2/solutions",
	}
	values := url.Values{}
	values.Add("$expand", "publisherid")
	values.Add("$filter", "(isvisible eq true)")
	values.Add("$orderby", "createdon desc")
	apiUrl.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, "GET", apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	apiResponse, err := client.doRequest(ctx, environmentUrl, request)
	if err != nil {
		return nil, err
	}

	solutionArray := models.SolutionDtoArray{}
	err = apiResponse.MarshallTo(&solutionArray)
	if err != nil {
		return nil, err
	}

	for inx := range solutionArray.Value {
		solutionArray.Value[inx].EnvironmentName = environmentId
	}

	solutions := make([]models.SolutionDto, 0)
	solutions = append(solutions, solutionArray.Value...)

	return solutions, nil

}
