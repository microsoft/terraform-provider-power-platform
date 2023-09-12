package powerplatform_bapi

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi/models"
)

func (client *ApiClient) GetSolution(ctx context.Context, environmentId string, solutionName string) (*models.SolutionDto, error) {
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

func (client *ApiClient) GetSolutions(ctx context.Context, environmentId string) ([]models.SolutionDto, error) {
	environmentUrl, token, err := client.getEnvironmentAuthDetails(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(*environmentUrl, "https://"),
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

	request.Header.Set("Authorization", "Bearer "+*token)
	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}

	solutionArray := models.SolutionDtoArray{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&solutionArray)
	if err != nil {
		return nil, err
	}

	for inx, _ := range solutionArray.Value {
		solutionArray.Value[inx].EnvironmentName = environmentId
	}

	solutions := make([]models.SolutionDto, 0)
	solutions = append(solutions, solutionArray.Value...)

	return solutions, nil
}

func (client *ApiClient) CreateSolution(ctx context.Context, environmentId string, solutionToCreate models.ImportSolutionDto, content []byte, settings []byte) (*models.SolutionDto, error) {

	environmentUrl, token, err := client.getEnvironmentAuthDetails(ctx, environmentId)
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
		Host:   strings.TrimPrefix(*environmentUrl, "https://"),
		Path:   "/api/data/v9.2/StageSolution",
	}

	stageSolutionRequest, err := http.NewRequestWithContext(ctx, "POST", apiUrl.String(), bytes.NewReader(stageSolutionRequestBody))
	if err != nil {
		return nil, err
	}
	stageSolutionRequest.Header.Set("Authorization", "Bearer "+*token)
	stageSolutionResponseBody, err := client.doRequest(stageSolutionRequest)
	if err != nil {
		return nil, err
	}
	stageSolutionResponse := models.StageSolutionImportResponseDto{}
	err = json.NewDecoder(bytes.NewReader(stageSolutionResponseBody)).Decode(&stageSolutionResponse)
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
		Host:   strings.TrimPrefix(*environmentUrl, "https://"),
		Path:   "/api/data/v9.2/ImportSolutionAsync",
	}
	importSolutionRequest, err := http.NewRequestWithContext(ctx, "POST", apiUrl.String(), bytes.NewReader(importSolutionRequestBody))
	if err != nil {
		return nil, err
	}
	importSolutionRequest.Header.Set("Authorization", "Bearer "+*token)
	importSolutionResponseBody, err := client.doRequest(importSolutionRequest)
	if err != nil {
		return nil, err
	}
	importSolutionResponse := models.ImportSolutionResponseDto{}
	err = json.NewDecoder(bytes.NewReader(importSolutionResponseBody)).Decode(&importSolutionResponse)
	if err != nil {
		return nil, err
	}

	//pull for solution import completion
	time.Sleep(10 * time.Second)

	apiUrl = &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(*environmentUrl, "https://"),
		Path:   fmt.Sprintf("/api/data/v9.2/asyncoperations(%s)", importSolutionResponse.AsyncOperationId),
	}
	asyncSolutionImportRequest, err := http.NewRequestWithContext(ctx, "GET", apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	asyncSolutionImportRequest.Header.Set("Authorization", "Bearer "+*token)
	for {
		asyncSolutionImportResponseBody, err := client.doRequest(asyncSolutionImportRequest)
		if err != nil {
			return nil, err
		}
		asyncSolutionPullResponse := models.AsyncSolutionPullResponseDto{}
		err = json.NewDecoder(bytes.NewReader(asyncSolutionImportResponseBody)).Decode(&asyncSolutionPullResponse)
		if err != nil {
			return nil, err
		}
		if asyncSolutionPullResponse.CompletedOn != "" {
			err = client.validateSolutionImportResult(ctx, *token, *environmentUrl, importSolutionResponse.ImportJobKey)
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

func (client *ApiClient) DeleteSolution(ctx context.Context, environmentId string, solutionName string) error {
	solution, err := client.GetSolution(ctx, environmentId, solutionName)
	if err != nil {
		return err
	}

	environmentUrl, token, err := client.getEnvironmentAuthDetails(ctx, environmentId)
	if err != nil {
		return err
	}
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(*environmentUrl, "https://"),
		Path:   fmt.Sprintf("/api/data/v9.2/solutions(%s)", solution.Id),
	}
	deleteSolutionRequest, err := http.NewRequestWithContext(ctx, "DELETE", apiUrl.String(), nil)
	if err != nil {
		return err
	}
	deleteSolutionRequest.Header.Set("Authorization", "Bearer "+*token)
	_, err = client.doRequest(deleteSolutionRequest)
	if err != nil {
		return err
	}
	return nil
}

func (client *ApiClient) createSolutionComponentParameters(ctx context.Context, settings []byte) ([]interface{}, error) {
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

func (client *ApiClient) getEnvironmentAuthDetails(ctx context.Context, environmentId string) (*string, *string, error) {
	env, err := client.GetEnvironment(ctx, environmentId)
	if err != nil {
		return nil, nil, err
	}
	environmentUrl := strings.TrimSuffix(env.Properties.LinkedEnvironmentMetadata.InstanceURL, "/")

	auth, err := client.DoAuthClientSecretForDataverse(ctx, environmentUrl)
	if err != nil {
		return nil, nil, err
	}
	tflog.Debug(ctx, "Dataverse_token:   "+auth.Token+"  ")
	return &environmentUrl, &auth.Token, nil
}

func (client *ApiClient) validateSolutionImportResult(ctx context.Context, token, environmentUrl, ImportJobKey string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   fmt.Sprintf("/api/data/v9.0/RetrieveSolutionImportResult(ImportJobId=%s)", ImportJobKey),
	}
	validateSolutionImportRequest, err := http.NewRequestWithContext(ctx, "GET", apiUrl.String(), nil)
	if err != nil {
		return err
	}
	validateSolutionImportRequest.Header.Set("Authorization", "Bearer "+token)
	validateSolutionImportResponseBody, err := client.doRequest(validateSolutionImportRequest)
	if err != nil {
		return err
	}

	validateSolutionImportResponseDto := models.ValidateSolutionImportResponseDto{}
	err = json.NewDecoder(bytes.NewReader(validateSolutionImportResponseBody)).Decode(&validateSolutionImportResponseDto)
	if err != nil {
		return err
	}
	if validateSolutionImportResponseDto.SolutionOperationResult.Status != "Passed" {
		//todo read error and warning messages
		return fmt.Errorf("solution import failed: %s", validateSolutionImportResponseDto.SolutionOperationResult.Status)
	}
	return nil
}
