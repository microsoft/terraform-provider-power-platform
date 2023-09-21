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

	powerplatform_bapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi/models"
)

var _ BapiClientInterface = &BapiClientImplementation{}

type BapiClientInterface interface {
	Initialize(context.Context) (string, error)

	GetEnvironments(ctx context.Context) ([]models.EnvironmentDto, error)
	GetEnvironment(ctx context.Context, environmentId string) (*models.EnvironmentDto, error)
}

type BapiClientImplementation struct {
	Config ProviderConfig
	Auth   BapiAuthInterface
}

func (client *BapiClientImplementation) doRequest(ctx context.Context, request *http.Request) (*powerplatform_bapi.ApiHttpResponse, error) {
	token, err := client.Initialize(ctx)
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

func (client *BapiClientImplementation) Initialize(ctx context.Context) (string, error) {

	if client.Auth.IsTokenExpiredOrEmpty() {
		if client.Config.Credentials.IsClientSecretCredentialsProvided() {
			token, err := client.Auth.AuthenticateClientSecret(ctx, client.Config.Credentials.TenantId, client.Config.Credentials.ClientId, client.Config.Credentials.Secret)
			if err != nil {
				return "", err
			}
			return token, nil
		} else if client.Config.Credentials.IsUserPassCredentialsProvided() {
			token, err := client.Auth.AuthenticateUserPass(ctx, client.Config.Credentials.TenantId, client.Config.Credentials.Username, client.Config.Credentials.Password)
			if err != nil {
				return "", err
			}
			return token, nil
		} else {
			return "", errors.New("no credentials provided")
		}
	} else {
		//todo this is not implemented yet
		token, err := client.Auth.RefreshToken()
		if err != nil {
			return "", err
		}
		return token, nil

	}
}

func (client *BapiClientImplementation) GetEnvironment(ctx context.Context, environmentId string) (*models.EnvironmentDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Config.Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("$expand", "permissions,properties.capacity")
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()
	request, err := http.NewRequestWithContext(ctx, "GET", apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	apiResponse, err := client.doRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	env := models.EnvironmentDto{}
	err = apiResponse.MarshallTo(&env)
	if err != nil {
		return nil, err
	}

	if env.Properties.LinkedEnvironmentMetadata.SecurityGroupId == "" {
		env.Properties.LinkedEnvironmentMetadata.SecurityGroupId = "00000000-0000-0000-0000-000000000000"
	}

	return &env, nil
}

func (client *BapiClientImplementation) GetEnvironments(ctx context.Context) ([]models.EnvironmentDto, error) {

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Config.Urls.BapiUrl,
		Path:   "/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments",
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()
	request, err := http.NewRequestWithContext(ctx, "GET", apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	apiResponse, err := client.doRequest(ctx, request)
	if err != nil {
		return nil, err
	}

	envArray := models.EnvironmentDtoArray{}
	err = apiResponse.MarshallTo(&envArray)
	if err != nil {
		return nil, err
	}

	return envArray.Value, nil
}
