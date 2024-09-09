// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package solution

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
)

func NewSolutionClient(api *api.ApiClient) SolutionClient {
	return SolutionClient{
		Api: api,
	}
}

type SolutionClient struct {
	Api *api.ApiClient
}

func (client *SolutionClient) DataverseExists(ctx context.Context, environmentId string) (bool, error) {

	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return false, err
	}
	return env.Properties.LinkedEnvironmentMetadata.InstanceURL != "", nil
}

func (client *SolutionClient) GetSolutionUniqueName(ctx context.Context, environmentId, name string) (*SolutionDto, error) {
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
	values.Add("$filter", fmt.Sprintf("uniquename eq '%s'", name))
	apiUrl.RawQuery = values.Encode()

	solutions := SolutionDtoArray{}
	_, err = client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &solutions)
	if err != nil {
		return nil, err
	}
	if len(solutions.Value) == 0 {
		return nil, fmt.Errorf("solution with unique name '%s' not found", name)
	}

	solutions.Value[0].EnvironmentId = environmentId

	return &solutions.Value[0], nil
}

func (client *SolutionClient) GetSolutionById(ctx context.Context, environmentId, solutionId string) (*SolutionDto, error) {
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
	values.Add("$filter", fmt.Sprintf("solutionid eq %s", solutionId))
	apiUrl.RawQuery = values.Encode()

	solutions := SolutionDtoArray{}
	_, err = client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &solutions)
	if err != nil {
		return nil, err
	}
	if len(solutions.Value) == 0 {
		return nil, fmt.Errorf("solution with id '%s' not found", solutionId)
	}

	solutions.Value[0].EnvironmentId = environmentId

	return &solutions.Value[0], nil
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
	_, err = client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &solutionArray)
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

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   "/api/data/v9.2/StageSolution",
	}

	stageSolutionResponse := StageSolutionImportResponseDto{}
	_, err = client.Api.Execute(ctx, "POST", apiUrl.String(), nil, stageSolutionRequestBody, []int{http.StatusOK}, &stageSolutionResponse)
	if err != nil {
		return nil, err
	}
	if stageSolutionResponse.StageSolutionResults.StageSolutionStatus != "Passed" {
		e := fmt.Errorf("solution failed with status: '%s'", stageSolutionResponse.StageSolutionResults.StageSolutionStatus)

		for _, missingDependency := range stageSolutionResponse.StageSolutionResults.MissingDependencies {
			e = errors.Join(fmt.Errorf("missing dependency: '%s'", missingDependency.RequiredComponentSchemaName), e)
		}
		for _, validation := range stageSolutionResponse.StageSolutionResults.SolutionValidationResults {
			e = errors.Join(fmt.Errorf("solution validation failed: %s", validation.Message), e)
		}
		return nil, e
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

	apiUrl = &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   "/api/data/v9.2/ImportSolutionAsync",
	}
	importSolutionResponse := ImportSolutionResponseDto{}
	_, err = client.Api.Execute(ctx, "POST", apiUrl.String(), nil, importSolutionRequestBody, []int{http.StatusOK}, &importSolutionResponse)
	if err != nil {
		return nil, err
	}

	//pull for solution import completion
	err = client.Api.SleepWithContext(ctx, client.Api.RetryAfterDefault())
	if err != nil {
		return nil, err
	}

	apiUrl = &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   fmt.Sprintf("/api/data/v9.2/asyncoperations(%s)", importSolutionResponse.AsyncOperationId),
	}
	for {
		asyncSolutionPullResponse := AsyncSolutionPullResponseDto{}
		_, err = client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &asyncSolutionPullResponse)
		if err != nil {
			return nil, err
		}
		if asyncSolutionPullResponse.CompletedOn != "" {
			err = client.validateSolutionImportResult(ctx, environmentUrl, importSolutionResponse.ImportJobKey)
			if err != nil {
				return nil, err
			}
			solution, err := client.GetSolutionUniqueName(ctx, environmentId, stageSolutionResponse.StageSolutionResults.SolutionDetails.SolutionUniqueName)
			if err != nil {
				return nil, err
			}
			return solution, nil
		}
		err = client.Api.SleepWithContext(ctx, client.Api.RetryAfterDefault())
		if err != nil {
			return nil, err
		}
	}
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
		if envVariableComponent.Value != "" {
			solutionComponents = append(solutionComponents, ImportSolutionEnvironmentVariablesDto{
				Type:       "Microsoft.Dynamics.CRM.environmentvariablevalue",
				SchemaName: envVariableComponent.SchemaName,
				Value:      envVariableComponent.Value,
			})
		}
	}

	if len(solutionComponents) == 0 {
		return nil, nil
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
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &validateSolutionImportResponseDto)
	if err != nil {
		return err
	}
	if validateSolutionImportResponseDto.SolutionOperationResult.Status != "Passed" {
		return fmt.Errorf("solution import failed: %s", validateSolutionImportResponseDto.SolutionOperationResult.ErrorMessages...)
	}
	return nil
}

func (client *SolutionClient) DeleteSolution(ctx context.Context, environmentId, solutionId string) error {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return err
	}
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   strings.TrimPrefix(environmentUrl, "https://"),
		Path:   fmt.Sprintf("/api/data/v9.2/solutions(%s)", solutionId),
	}
	_, err = client.Api.Execute(ctx, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusNoContent}, nil)
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
	_, err = client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &responseObj)
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
	if environmentUrl == "" {
		return "", helpers.WrapIntoProviderError(nil, helpers.ERROR_ENVIRONMENT_URL_NOT_FOUND, "environment url not found, please check if the environment has dataverse linked")
		return "", helpers.WrapIntoProviderError(nil, helpers.ERROR_ENVIRONMENT_URL_NOT_FOUND, "environment url not found, please check if the environment has dataverse linked")
	}
	return environmentUrl, nil
}

func (client *SolutionClient) getEnvironment(ctx context.Context, environmentId string) (*EnvironmentIdDto, error) {

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("$expand", "permissions,properties.capacity,properties/billingPolicy")
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	env := EnvironmentIdDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
	if err != nil {
		if strings.ContainsAny(err.Error(), "404") {
			return nil, helpers.WrapIntoProviderError(err, helpers.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("environment %s not found", environmentId))
			return nil, helpers.WrapIntoProviderError(err, helpers.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("environment %s not found", environmentId))
		}
		return nil, err

	}

	return &env, nil
}
