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

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

func NewSolutionClient(apiClient *api.Client) Client {
	return Client{
		Api: apiClient,
	}
}

type Client struct {
	Api *api.Client
}

func (client *Client) DataverseExists(ctx context.Context, environmentId string) (bool, error) {
	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return false, err
	}
	return env.Properties.LinkedEnvironmentMetadata.InstanceURL != "", nil
}

func (client *Client) GetSolutionUniqueName(ctx context.Context, environmentId, name string) (*solutionDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.2/solutions",
	}
	values := url.Values{}
	values.Add("$expand", "publisherid")
	values.Add("$filter", fmt.Sprintf("uniquename eq '%s'", name))
	apiUrl.RawQuery = values.Encode()

	solutions := solutionArrayDto{}
	_, err = client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &solutions)
	if err != nil {
		return nil, err
	}
	if len(solutions.Value) == 0 {
		return nil, helpers.WrapIntoProviderError(err, helpers.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("solution with unique name '%s' not found", name))
	}

	solutions.Value[0].EnvironmentId = environmentId

	return &solutions.Value[0], nil
}

func (client *Client) GetSolutionById(ctx context.Context, environmentId, solutionId string) (*solutionDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.2/solutions",
	}
	values := url.Values{}
	values.Add("$expand", "publisherid")
	values.Add("$filter", fmt.Sprintf("solutionid eq %s", solutionId))
	apiUrl.RawQuery = values.Encode()

	solutions := solutionArrayDto{}
	_, err = client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &solutions)
	if err != nil {
		return nil, err
	}
	if len(solutions.Value) == 0 {
		return nil, helpers.WrapIntoProviderError(err, helpers.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("solution with id '%s' not found", solutionId))
	}

	solutions.Value[0].EnvironmentId = environmentId

	return &solutions.Value[0], nil
}

func (client *Client) GetSolutions(ctx context.Context, environmentId string) ([]solutionDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.2/solutions",
	}
	values := url.Values{}
	values.Add("$expand", "publisherid")
	values.Add("$filter", "(isvisible eq true)")
	values.Add("$orderby", "createdon desc")
	apiUrl.RawQuery = values.Encode()

	solutionArray := solutionArrayDto{}
	_, err = client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &solutionArray)
	if err != nil {
		return nil, err
	}

	for inx := range solutionArray.Value {
		solutionArray.Value[inx].EnvironmentId = environmentId
	}

	solutions := make([]solutionDto, 0)
	solutions = append(solutions, solutionArray.Value...)

	return solutions, nil
}

func (client *Client) CreateSolution(ctx context.Context, environmentId string, solutionToCreate importSolutionDto, content []byte, settings []byte) (*solutionDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	if content == nil {
		err = fmt.Errorf("solution content is nil")
		return nil, err
	}

	stageSolutionRequestBody := stageSolutionImportDto{
		CustomizationFile: base64.StdEncoding.EncodeToString(content),
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.2/StageSolution",
	}

	stageSolutionResponse := stageSolutionImportResponseDto{}
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

	solutionComponents, err := client.createSolutionComponentParameters(settings)
	if err != nil {
		return nil, err
	}

	importSolutionRequestBody := importSolutionDto{
		PublishWorkflows:                 true,
		OverwriteUnmanagedCustomizations: false,
		ComponentParameters:              solutionComponents,
		SolutionParameters: importSolutionSolutionParametersDto{
			StageSolutionUploadId: stageSolutionResponse.StageSolutionResults.StageSolutionUploadId,
		},
	}

	apiUrl = &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   "/api/data/v9.2/ImportSolutionAsync",
	}
	importSolutionResponse := importSolutionResponseDto{}
	_, err = client.Api.Execute(ctx, "POST", apiUrl.String(), nil, importSolutionRequestBody, []int{http.StatusOK}, &importSolutionResponse)
	if err != nil {
		return nil, err
	}

	// pull for solution import completion.
	err = client.Api.SleepWithContext(ctx, client.Api.RetryAfterDefault())
	if err != nil {
		return nil, err
	}

	apiUrl = &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/v9.2/asyncoperations(%s)", importSolutionResponse.AsyncOperationId),
	}
	for {
		asyncSolutionPullResponse := asyncSolutionPullResponseDto{}
		_, err = client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &asyncSolutionPullResponse)
		if err != nil {
			return nil, err
		}
		if asyncSolutionPullResponse.CompletedOn != "" {
			err = client.validateSolutionImportResult(ctx, environmentHost, importSolutionResponse.ImportJobKey)
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

func (client *Client) createSolutionComponentParameters(settings []byte) ([]any, error) {
	if len(settings) == 0 {
		return nil, nil
	}

	solutionSettings := solutionSettings{}
	if settings != nil {
		err := json.Unmarshal(settings, &solutionSettings)
		if err != nil {
			return nil, err
		}
	}

	solutionComponents := make([]any, 0)
	for _, connectionReferenceComponent := range solutionSettings.ConnectionReferences {
		solutionComponents = append(solutionComponents, importSolutionConnectionReferencesDto{
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
			solutionComponents = append(solutionComponents, importSolutionEnvironmentVariablesDto{
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

func (client *Client) validateSolutionImportResult(ctx context.Context, environmentHost, importJobKey string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/v9.0/RetrieveSolutionImportResult(ImportJobId=%s)", importJobKey),
	}

	validateSolutionImportResponseDto := validateSolutionImportResponseDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &validateSolutionImportResponseDto)
	if err != nil {
		return err
	}
	if validateSolutionImportResponseDto.SolutionOperationResult.Status != "Passed" {
		return fmt.Errorf("solution import failed: %s", validateSolutionImportResponseDto.SolutionOperationResult.ErrorMessages...)
	}
	return nil
}

func (client *Client) DeleteSolution(ctx context.Context, environmentId, solutionId string) error {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return err
	}
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/v9.2/solutions(%s)", solutionId),
	}
	_, err = client.Api.Execute(ctx, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusNoContent}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (client *Client) GetTableData(ctx context.Context, environmentId, tableName, odataQuery string, responseObj any) error {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return err
	}
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
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

func (client *Client) GetEnvironmentHostById(ctx context.Context, environmentId string) (string, error) {
	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return "", err
	}
	environmentUrl := strings.TrimSuffix(env.Properties.LinkedEnvironmentMetadata.InstanceURL, "/")
	if environmentUrl == "" {
		return "", helpers.WrapIntoProviderError(nil, helpers.ERROR_ENVIRONMENT_URL_NOT_FOUND, "environment url not found, please check if the environment has dataverse linked")
	}

	envUrl, err := url.Parse(environmentUrl)
	if err != nil {
		return "", err
	}
	return envUrl.Host, nil
}

func (client *Client) getEnvironment(ctx context.Context, environmentId string) (*environmentIdDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("$expand", "permissions,properties.capacity,properties/billingPolicy")
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	env := environmentIdDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
	if err != nil {
		if strings.ContainsAny(err.Error(), "404") {
			return nil, helpers.WrapIntoProviderError(err, helpers.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("environment %s not found", environmentId))
		}
		return nil, err
	}

	return &env, nil
}
