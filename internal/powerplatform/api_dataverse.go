package powerplatform

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	powerplatform_bapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi/models"
)

var _ DataverseClientInterface = &DataverseClientImplementation{}

type DataverseClientInterface interface {
	Initialize(ctx context.Context, environmentUrl string) (string, error)

	GetSolutions(ctx context.Context, environmentId string) ([]models.SolutionDto, error)
}

type DataverseClientImplementation struct {
	Config     ProviderConfig
	Auth       DataverseAuthInterface
	BapiClient BapiClientInterface
}

// TODO remove duplicate method that is the same for all clients
func (client *DataverseClientImplementation) doRequest(ctx context.Context, environmentUrl string, request *http.Request) (*powerplatform_bapi.ApiHttpResponse, error) {
	token, err := client.Initialize(ctx, environmentUrl)
	if err != nil {
		return nil, err
	}

	apiHttpResponse := &powerplatform_bapi.ApiHttpResponse{}

	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}

	//todo validate that initializing the http client everytime is ok from performance perspective
	httpClient := http.DefaultClient

	if request.Header["Authorization"] == nil {
		request.Header.Set("Authorization", "Bearer "+token)
	}

	request.Header.Set("User-Agent", "terraform-provider-power-platform")

	response, err := httpClient.Do(request)
	apiHttpResponse.Response = response
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(response.Body)
	apiHttpResponse.BodyAsBytes = body
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		if len(body) != 0 {
			errorResponse := make(map[string]interface{}, 0)
			err = json.NewDecoder(bytes.NewBuffer(body)).Decode(&errorResponse)
			if err != nil {
				return nil, err
			}

			return apiHttpResponse, fmt.Errorf("status: %d, body: %s", response.StatusCode, errorResponse)
		} else {
			return nil, fmt.Errorf("status: %d", response.StatusCode)
		}
	}
	return apiHttpResponse, nil
}

func (client *DataverseClientImplementation) Initialize(ctx context.Context, environmentUrl string) (string, error) {

	if client.Auth.IsTokenExpiredOrEmpty(environmentUrl) {
		if client.Config.Credentials.IsClientSecretCredentialsProvided() {
			token, err := client.Auth.AuthenticateClientSecret(ctx, environmentUrl, client.Config.Credentials.TenantId, client.Config.Credentials.ClientId, client.Config.Credentials.Secret)
			if err != nil {
				return "", err
			}
			return token, nil
		} else if client.Config.Credentials.IsUserPassCredentialsProvided() {
			token, err := client.Auth.AuthenticateUserPass(ctx, environmentUrl, client.Config.Credentials.TenantId, client.Config.Credentials.Username, client.Config.Credentials.Password)
			if err != nil {
				return "", err
			}
			return token, nil
		} else {
			return "", errors.New("no credentials provided")
		}
	} else {
		//todo this is not implemented yet
		token, err := client.Auth.RefreshToken(environmentUrl)
		if err != nil {
			return "", err
		}
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
