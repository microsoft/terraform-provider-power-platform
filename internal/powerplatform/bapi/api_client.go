package powerplatform_bapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi/models"
)

type ApiClient struct {
	HttpClient      *http.Client
	Token           string
	BapiUrl         string
	PowerAppsApiUrl string

	Provider         *Provider
	DataverseAuthMap map[string]*AuthResponse
}

type Provider struct {
	TenantId     string
	ClientId     string
	ClientSecret string

	Username string
	Password string
}

var _ ApiClientInterface = &ApiClient{}

//go:generate mockgen -destination=../../mocks/client_mocks_bapi.go -package=powerplatform_mocks "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi" ApiClientInterface
type ApiClientInterface interface {
	DoAuthClientSecret(ctx context.Context, tenantId, applicationId, clientSecret string) (*AuthResponse, error)
	DoAuthUsernamePassword(ctx context.Context, tenantId, username, password string) (*AuthResponse, error)

	DoAuthClientSecretForDataverse(ctx context.Context, environmentUrl string) (*AuthResponse, error)

	GetEnvironments(ctx context.Context) ([]models.EnvironmentDto, error)
	GetEnvironment(ctx context.Context, environmentId string) (*models.EnvironmentDto, error)
	CreateEnvironment(ctx context.Context, environment models.EnvironmentCreateDto) (*models.EnvironmentDto, error)
	UpdateEnvironment(ctx context.Context, environmentId string, environment models.EnvironmentDto) (*models.EnvironmentDto, error)
	DeleteEnvironment(ctx context.Context, environmentId string) error

	CreateSolution(ctx context.Context, environmentId string, solutionToCreate models.ImportSolutionDto, content []byte, settings []byte) (*models.SolutionDto, error)
	GetSolutions(ctx context.Context, environmentId string) ([]models.SolutionDto, error)
	GetSolution(ctx context.Context, environmentId string, solutionName string) (*models.SolutionDto, error)
	DeleteSolution(ctx context.Context, environmentId string, solutionName string) error

	GetPowerApps(ctx context.Context, environmentId string) ([]models.PowerAppBapi, error)

	GetConnectors(ctx context.Context) ([]models.ConnectorDto, error)
	GetPolicies(ctx context.Context) ([]models.DlpPolicyModel, error)
	GetPolicy(ctx context.Context, name string) (*models.DlpPolicyModel, error)
	DeletePolicy(ctx context.Context, name string) error
	UpdatePolicy(ctx context.Context, name string, policyToUpdate models.DlpPolicyModel) (*models.DlpPolicyModel, error)
	CreatePolicy(ctx context.Context, policyToCreate models.DlpPolicyModel) (*models.DlpPolicyModel, error)
}

type ApiHttpResponse struct {
	Response    *http.Response
	BodyAsBytes []byte
}

func (apiResponse *ApiHttpResponse) MarshallTo(obj interface{}) error {
	err := json.NewDecoder(bytes.NewReader(apiResponse.BodyAsBytes)).Decode(&obj)
	if err != nil {
		return err
	}
	return nil
}

func (apiResponse *ApiHttpResponse) GetHeader(name string) string {
	return apiResponse.Response.Header.Get(name)
}

func (ApiHttpResponse *ApiHttpResponse) ValidateStatusCode(expectedStatusCode int) error {
	if ApiHttpResponse.Response.StatusCode != expectedStatusCode {
		return fmt.Errorf("expected status code %d, got %d", expectedStatusCode, ApiHttpResponse.Response.StatusCode)
	}
	return nil
}

func (client *ApiClient) doRequest(request *http.Request) (*ApiHttpResponse, error) {

	apiHttpResponse := &ApiHttpResponse{}

	if request.Header.Get("Content-Type") == "" {
		request.Header.Set("Content-Type", "application/json")
	}

	client.HttpClient = http.DefaultClient

	if request.Header["Authorization"] == nil {
		request.Header.Set("Authorization", "Bearer "+client.Token)
	}

	request.Header.Set("User-Agent", "terraform-provider-power-platform")

	response, err := client.HttpClient.Do(request)
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
