package powerplatform

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	bapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api/bapi"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/models"
)

var _ DataverseClientInterface = &DataverseClientImplementation{}

type DataverseClientInterface interface {
	Initialize(ctx context.Context, environmentUrl string) (string, error)

	GetSolutions(ctx context.Context, environmentId string) ([]models.SolutionDto, error)
	CreateSolution(ctx context.Context, environmentId string, solutionToCreate models.ImportSolutionDto, content []byte, settings []byte) (*models.SolutionDto, error)
	GetSolution(ctx context.Context, environmentId string, solutionName string) (*models.SolutionDto, error)
	DeleteSolution(ctx context.Context, environmentId string, solutionName string) error
}

type DataverseClientImplementation struct {
	BaseApi    api.ApiClientInterface
	Auth       DataverseAuthInterface
	BapiClient bapi.BapiClientInterface
}

func (client *DataverseClientImplementation) doRequest(ctx context.Context, environmentUrl string, request *http.Request) (*api.ApiHttpResponse, error) {
	token, err := client.Initialize(ctx, environmentUrl)
	if err != nil {
		return nil, err
	}
	return client.BaseApi.DoRequest(token, request)
}

func (client *DataverseClientImplementation) getEnvironmentUrlById(ctx context.Context, environmentId string) (string, error) {
	env, err := client.BapiClient.GetEnvironment(ctx, environmentId)
	if err != nil {
		return "", err
	}
	environmentUrl := strings.TrimSuffix(env.Properties.LinkedEnvironmentMetadata.InstanceURL, "/")
	return environmentUrl, nil
}

func (client *DataverseClientImplementation) Initialize(ctx context.Context, environmentUrl string) (string, error) {

	token, err := client.Auth.GetToken(environmentUrl)

	if _, ok := err.(*api.TokeExpiredError); ok {
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

func (client *DataverseClientImplementation) CreateSolution(ctx context.Context, environmentId string, solutionToCreate models.ImportSolutionDto, content []byte, settings []byte) (*models.SolutionDto, error) {
	environmentUrl, err := client.getEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	if content == nil {
		err = fmt.Errorf("solution content is nil")
		return nil, err
	}

	//stage solution
	stageSolutionRequestBody, err := json.Marshal(models.StageSolutionImportDto{
		CustomizationFile: base64.StdEncoding.EncodeToString(content),
	})
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   "/api/data/v9.2/StageSolution",
	}

	stageSolutionRequest, err := http.NewRequestWithContext(ctx, "POST", apiUrl.String(), bytes.NewReader(stageSolutionRequestBody))
	if err != nil {
		return nil, err
	}
	apiResponse, err := client.doRequest(ctx, environmentUrl, stageSolutionRequest)
	if err != nil {
		return nil, err
	}
	stageSolutionResponse := models.StageSolutionImportResponseDto{}
	err = apiResponse.MarshallTo(&stageSolutionResponse)
	if err != nil {
		return nil, err
	}
	if stageSolutionResponse.StageSolutionResults.StageSolutionStatus != "Passed" {
		return nil, fmt.Errorf("stage solution failed: %s", stageSolutionResponse.StageSolutionResults.StageSolutionStatus)
	}

	//import solution
	solutionComponents, err := client.createSolutionComponentParameters(ctx, settings)
	if err != nil {
		return nil, err
	}

	importSolutionRequestBody, err := json.Marshal(models.ImportSolutionDto{
		PublishWorkflows:                 true,
		OverwriteUnmanagedCustomizations: false,
		ComponentParameters:              solutionComponents,
		SolutionParameters: models.ImportSolutionSolutionParametersDto{
			StageSolutionUploadId: stageSolutionResponse.StageSolutionResults.StageSolutionUploadId,
		},
	})
	if err != nil {
		return nil, err
	}

	apiUrl = &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   "/api/data/v9.2/ImportSolutionAsync",
	}
	importSolutionRequest, err := http.NewRequestWithContext(ctx, "POST", apiUrl.String(), bytes.NewReader(importSolutionRequestBody))
	if err != nil {
		return nil, err
	}
	apiResponse, err = client.doRequest(ctx, environmentUrl, importSolutionRequest)
	if err != nil {
		return nil, err
	}
	importSolutionResponse := models.ImportSolutionResponseDto{}
	err = apiResponse.MarshallTo(&importSolutionResponse)
	if err != nil {
		return nil, err
	}

	//pull for solution import completion
	time.Sleep(10 * time.Second)

	apiUrl = &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   fmt.Sprintf("/api/data/v9.2/asyncoperations(%s)", importSolutionResponse.AsyncOperationId),
	}
	asyncSolutionImportRequest, err := http.NewRequestWithContext(ctx, "GET", apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	for {
		apiResponse, err := client.doRequest(ctx, environmentUrl, asyncSolutionImportRequest)
		if err != nil {
			return nil, err
		}
		asyncSolutionPullResponse := models.AsyncSolutionPullResponseDto{}
		err = apiResponse.MarshallTo(&asyncSolutionPullResponse)
		if err != nil {
			return nil, err
		}
		if asyncSolutionPullResponse.CompletedOn != "" {
			err = client.validateSolutionImportResult(ctx, environmentUrl, importSolutionResponse.ImportJobKey)
			if err != nil {
				return nil, err
			}
			solution, err := client.GetSolution(ctx, environmentId, stageSolutionResponse.StageSolutionResults.SolutionDetails.SolutionUniqueName)
			if err != nil {
				return nil, err
			}
			return solution, nil
		}
		time.Sleep(10 * time.Second)
	}
}

func (client *DataverseClientImplementation) GetSolution(ctx context.Context, environmentId string, solutionName string) (*models.SolutionDto, error) {
	solutions, err := client.GetSolutions(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	for _, solution := range solutions {
		if strings.EqualFold(solution.Name, solutionName) {
			return &solution, nil
		}
	}
	return nil, fmt.Errorf("solution %s not found in %s", solutionName, environmentId)
}

func (client *DataverseClientImplementation) createSolutionComponentParameters(ctx context.Context, settings []byte) ([]interface{}, error) {
	if len(settings) == 0 {
		return nil, nil
	}

	solutionSettings := models.SolutionSettings{}
	if settings != nil {
		err := json.Unmarshal(settings, &solutionSettings)
		if err != nil {
			return nil, err
		}
	}

	solutionComponents := make([]interface{}, 0)
	for _, connectionReferenceComponent := range solutionSettings.ConnectionReferences {
		solutionComponents = append(solutionComponents, models.ImportSolutionConnectionReferencesDto{
			Type:                           "Microsoft.Dynamics.CRM.connectionreference",
			ConnectionId:                   connectionReferenceComponent.ConnectionId,
			ConnectorId:                    connectionReferenceComponent.ConnectorId,
			ConnectionReferenceLogicalName: connectionReferenceComponent.LogicalName,
			ConnectionReferenceDisplayName: "",
			Description:                    "",
		})
	}
	for _, envVariableComponent := range solutionSettings.EnvironmentVariables {
		solutionComponents = append(solutionComponents, models.ImportSolutionEnvironmentVariablesDto{
			Type:       "Microsoft.Dynamics.CRM.environmentvariablevalue",
			SchemaName: envVariableComponent.SchemaName,
			Value:      envVariableComponent.Value,
		})
	}

	return solutionComponents, nil
}

func (client *DataverseClientImplementation) validateSolutionImportResult(ctx context.Context, environmentUrl, ImportJobKey string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   fmt.Sprintf("/api/data/v9.0/RetrieveSolutionImportResult(ImportJobId=%s)", ImportJobKey),
	}
	validateSolutionImportRequest, err := http.NewRequestWithContext(ctx, "GET", apiUrl.String(), nil)
	if err != nil {
		return err
	}
	apiResponse, err := client.doRequest(ctx, environmentUrl, validateSolutionImportRequest)
	if err != nil {
		return err
	}

	validateSolutionImportResponseDto := models.ValidateSolutionImportResponseDto{}
	err = apiResponse.MarshallTo(&validateSolutionImportResponseDto)
	if err != nil {
		return err
	}
	if validateSolutionImportResponseDto.SolutionOperationResult.Status != "Passed" {
		//todo read error and warning messages
		return fmt.Errorf("solution import failed: %s", validateSolutionImportResponseDto.SolutionOperationResult.Status)
	}
	return nil
}

func (client *DataverseClientImplementation) DeleteSolution(ctx context.Context, environmentId string, solutionName string) error {
	solution, err := client.GetSolution(ctx, environmentId, solutionName)
	if err != nil {
		return err
	}

	environmentUrl, err := client.getEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return err
	}
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   fmt.Sprintf("/api/data/v9.2/solutions(%s)", solution.Id),
	}
	deleteSolutionRequest, err := http.NewRequestWithContext(ctx, "DELETE", apiUrl.String(), nil)
	if err != nil {
		return err
	}
	apiResponse, err := client.doRequest(ctx, environmentUrl, deleteSolutionRequest)
	if err != nil {
		return err
	}
	err = apiResponse.ValidateStatusCode(http.StatusNoContent)
	if err != nil {
		return err
	}
	return nil
}
