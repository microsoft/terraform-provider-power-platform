package powerplatform

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewSolutionClient(bapi *api.BapiClientApi, dv *api.DataverseClientApi) SolutionClient {
	return SolutionClient{
		dataverseClient: dv,
		bapiClient:      bapi,
	}
}

type SolutionClient struct {
	dataverseClient *api.DataverseClientApi
	bapiClient      *api.BapiClientApi
}

func (client *SolutionClient) GetSolutions(ctx context.Context, environmentId string) ([]SolutionDto, error) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
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

	solutionArray := SolutionDtoArray{}
	_, err = client.dataverseClient.Execute(ctx, environmentUrl, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &solutionArray)
	if err != nil {
		return nil, err
	}

	for inx := range solutionArray.Value {
		solutionArray.Value[inx].EnvironmentId = environmentId
	}

	solutions := make([]SolutionDto, 0)
	solutions = append(solutions, solutionArray.Value...)

	return solutions, nil
}

func (client *SolutionClient) CreateSolution(ctx context.Context, environmentId string, solutionToCreate ImportSolutionDto, content []byte, settings []byte) (*SolutionDto, error) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	if content == nil {
		err = fmt.Errorf("solution content is nil")
		return nil, err
	}

	stageSolutionRequestBody := StageSolutionImportDto{
		CustomizationFile: base64.StdEncoding.EncodeToString(content),
	}
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   "/api/data/v9.2/StageSolution",
	}

	stageSolutionResponse := StageSolutionImportResponseDto{}
	_, err = client.dataverseClient.Execute(ctx, environmentUrl, "POST", apiUrl.String(), stageSolutionRequestBody, []int{http.StatusOK}, &stageSolutionResponse)
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

	importSolutionRequestBody := ImportSolutionDto{
		PublishWorkflows:                 true,
		OverwriteUnmanagedCustomizations: false,
		ComponentParameters:              solutionComponents,
		SolutionParameters: ImportSolutionSolutionParametersDto{
			StageSolutionUploadId: stageSolutionResponse.StageSolutionResults.StageSolutionUploadId,
		},
	}
	if err != nil {
		return nil, err
	}

	apiUrl = &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   "/api/data/v9.2/ImportSolutionAsync",
	}
	importSolutionResponse := ImportSolutionResponseDto{}
	_, err = client.dataverseClient.Execute(ctx, environmentUrl, "POST", apiUrl.String(), importSolutionRequestBody, []int{http.StatusOK}, &importSolutionResponse)
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
	for {
		asyncSolutionPullResponse := AsyncSolutionPullResponseDto{}
		_, err = client.dataverseClient.Execute(ctx, environmentUrl, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &asyncSolutionPullResponse)
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

func (client *SolutionClient) GetSolution(ctx context.Context, environmentId string, solutionName string) (*SolutionDto, error) {
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

func (client *SolutionClient) createSolutionComponentParameters(ctx context.Context, settings []byte) ([]interface{}, error) {
	if len(settings) == 0 {
		return nil, nil
	}

	solutionSettings := SolutionSettings{}
	if settings != nil {
		err := json.Unmarshal(settings, &solutionSettings)
		if err != nil {
			return nil, err
		}
	}

	solutionComponents := make([]interface{}, 0)
	for _, connectionReferenceComponent := range solutionSettings.ConnectionReferences {
		solutionComponents = append(solutionComponents, ImportSolutionConnectionReferencesDto{
			Type:                           "Microsoft.Dynamics.CRM.connectionreference",
			ConnectionId:                   connectionReferenceComponent.ConnectionId,
			ConnectorId:                    connectionReferenceComponent.ConnectorId,
			ConnectionReferenceLogicalName: connectionReferenceComponent.LogicalName,
			ConnectionReferenceDisplayName: "",
			Description:                    "",
		})
	}
	for _, envVariableComponent := range solutionSettings.EnvironmentVariables {
		solutionComponents = append(solutionComponents, ImportSolutionEnvironmentVariablesDto{
			Type:       "Microsoft.Dynamics.CRM.environmentvariablevalue",
			SchemaName: envVariableComponent.SchemaName,
			Value:      envVariableComponent.Value,
		})
	}

	return solutionComponents, nil
}

func (client *SolutionClient) validateSolutionImportResult(ctx context.Context, environmentUrl, ImportJobKey string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   fmt.Sprintf("/api/data/v9.0/RetrieveSolutionImportResult(ImportJobId=%s)", ImportJobKey),
	}

	validateSolutionImportResponseDto := ValidateSolutionImportResponseDto{}
	_, err := client.dataverseClient.Execute(ctx, environmentUrl, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &validateSolutionImportResponseDto)
	if err != nil {
		return err
	}
	if validateSolutionImportResponseDto.SolutionOperationResult.Status != "Passed" {
		//todo read error and warning messages
		return fmt.Errorf("solution import failed: %s", validateSolutionImportResponseDto.SolutionOperationResult.Status)
	}
	return nil
}

func (client *SolutionClient) DeleteSolution(ctx context.Context, environmentId string, solutionName string) error {
	solution, err := client.GetSolution(ctx, environmentId, solutionName)
	if err != nil {
		return err
	}

	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return err
	}
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   fmt.Sprintf("/api/data/v9.2/solutions(%s)", solution.Id),
	}
	_, err = client.dataverseClient.Execute(ctx, environmentUrl, "DELETE", apiUrl.String(), nil, []int{http.StatusNoContent}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (client *SolutionClient) GetTableData(ctx context.Context, environmentId, tableName, odataQuery string, responseObj interface{}) error {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return err
	}
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   fmt.Sprintf("/api/data/v9.2/%s", tableName),
	}
	if odataQuery != "" {
		apiUrl.RawQuery = odataQuery
	}
	_, err = client.dataverseClient.Execute(ctx, environmentUrl, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &responseObj)
	if err != nil {
		return err
	}
	return nil
}

func (client *SolutionClient) GetEnvironmentUrlById(ctx context.Context, environmentId string) (string, error) {
	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return "", err
	}
	environmentUrl := strings.TrimSuffix(env.Properties.LinkedEnvironmentMetadata.InstanceURL, "/")
	return environmentUrl, nil
}

func (client *SolutionClient) getEnvironment(ctx context.Context, environmentId string) (*EnvironmentIdDto, error) {

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.dataverseClient.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("$expand", "permissions,properties.capacity")
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	env := EnvironmentIdDto{}
	_, err := client.bapiClient.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}
