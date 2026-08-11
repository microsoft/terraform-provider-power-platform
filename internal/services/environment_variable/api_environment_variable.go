// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environmentvariable

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/solution"
)

func newEnvironmentVariableClient(apiClient *api.Client) client {
	return client{
		Api:            apiClient,
		SolutionClient: solution.NewSolutionClient(apiClient),
	}
}

type client struct {
	Api            *api.Client
	SolutionClient solution.Client
}

func (client *client) DataverseExists(ctx context.Context, environmentId string) (bool, error) {
	return client.SolutionClient.DataverseExists(ctx, environmentId)
}

func (client *client) GetEnvironmentVariable(ctx context.Context, environmentId, schemaName string) (*resolvedEnvironmentVariableDto, error) {
	definition, err := client.getDefinitionBySchemaName(ctx, environmentId, schemaName)
	if err != nil {
		return nil, err
	}

	value, err := client.getCurrentValueByDefinitionID(ctx, environmentId, definition.EnvironmentVariableDefinitionId)
	if err != nil {
		return nil, err
	}

	return &resolvedEnvironmentVariableDto{
		Definition: *definition,
		Value:      value,
	}, nil
}

func (client *client) UpsertEnvironmentVariableValue(ctx context.Context, environmentId, schemaName, value string) (*resolvedEnvironmentVariableDto, error) {
	definition, err := client.getDefinitionBySchemaName(ctx, environmentId, schemaName)
	if err != nil {
		return nil, err
	}

	currentValue, err := client.getCurrentValueByDefinitionID(ctx, environmentId, definition.EnvironmentVariableDefinitionId)
	if err != nil {
		return nil, err
	}

	if currentValue == nil {
		currentValue, err = client.createCurrentValue(ctx, environmentId, definition.EnvironmentVariableDefinitionId, schemaName, value)
		if err != nil {
			return nil, err
		}
	} else {
		currentValue, err = client.updateCurrentValue(ctx, environmentId, currentValue.EnvironmentVariableValueId, value)
		if err != nil {
			return nil, err
		}
	}

	return &resolvedEnvironmentVariableDto{
		Definition: *definition,
		Value:      currentValue,
	}, nil
}

func (client *client) DeleteEnvironmentVariableValue(ctx context.Context, environmentId, schemaName string) error {
	definition, err := client.getDefinitionBySchemaName(ctx, environmentId, schemaName)
	if err != nil {
		if isEnvironmentVariableDefinitionNotFound(err) {
			return nil
		}
		return err
	}

	currentValue, err := client.getCurrentValueByDefinitionID(ctx, environmentId, definition.EnvironmentVariableDefinitionId)
	if err != nil {
		return err
	}
	if currentValue == nil {
		return nil
	}

	environmentHost, err := client.SolutionClient.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return err
	}

	apiURL := helpers.BuildApiUrl(environmentHost, fmt.Sprintf("/api/data/%s/environmentvariablevalues(%s)", constants.DATAVERSE_API_VERSION, currentValue.EnvironmentVariableValueId), nil)
	resp, err := client.Api.Execute(ctx, nil, "DELETE", apiURL, nil, nil, []int{http.StatusNoContent, http.StatusNotFound, http.StatusForbidden}, nil)
	if err != nil {
		return err
	}
	return client.Api.HandleForbiddenResponse(resp)
}

func (client *client) getDefinitionBySchemaName(ctx context.Context, environmentId, schemaName string) (*environmentVariableDefinitionDto, error) {
	environmentHost, err := client.SolutionClient.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Add("$select", "environmentvariabledefinitionid,schemaname,displayname,description,defaultvalue,type,valueschema,secretstore")
	values.Add("$filter", fmt.Sprintf("schemaname eq '%s'", escapeODataStringLiteral(schemaName)))

	apiURL := helpers.BuildApiUrl(environmentHost, fmt.Sprintf("/api/data/%s/environmentvariabledefinitions", constants.DATAVERSE_API_VERSION), values)

	response := environmentVariableDefinitionArrayDto{}
	resp, err := client.Api.Execute(ctx, nil, "GET", apiURL, nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &response)
	if err != nil {
		return nil, err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := client.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, err
	}

	switch len(response.Value) {
	case 0:
		return nil, customerrors.WrapIntoProviderError(
			fmt.Errorf("environment variable definition with schema name '%s' not found", schemaName),
			customerrors.ErrorCode(constants.ERROR_OBJECT_NOT_FOUND),
			fmt.Sprintf("environment variable definition with schema name '%s' not found", schemaName),
		)
	case 1:
		return &response.Value[0], nil
	default:
		return nil, fmt.Errorf("multiple environment variable definitions found for schema name '%s'", schemaName)
	}
}

func (client *client) getCurrentValueByDefinitionID(ctx context.Context, environmentId, definitionID string) (*environmentVariableValueDto, error) {
	environmentHost, err := client.SolutionClient.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Add("$select", "environmentvariablevalueid,schemaname,value")
	values.Add("$filter", fmt.Sprintf("_environmentvariabledefinitionid_value eq %s", definitionID))

	apiURL := helpers.BuildApiUrl(environmentHost, fmt.Sprintf("/api/data/%s/environmentvariablevalues", constants.DATAVERSE_API_VERSION), values)

	response := environmentVariableValueArrayDto{}
	resp, err := client.Api.Execute(ctx, nil, "GET", apiURL, nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &response)
	if err != nil {
		return nil, err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := client.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, err
	}

	switch len(response.Value) {
	case 0:
		return nil, nil
	case 1:
		return &response.Value[0], nil
	default:
		return nil, fmt.Errorf("multiple current values found for environment variable definition '%s'", definitionID)
	}
}

func (client *client) createCurrentValue(ctx context.Context, environmentId, definitionID, schemaName, value string) (*environmentVariableValueDto, error) {
	environmentHost, err := client.SolutionClient.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	requestBody := createEnvironmentVariableValueDto{
		SchemaName:                            schemaName,
		Value:                                 value,
		EnvironmentVariableDefinitionODataRef: fmt.Sprintf("/environmentvariabledefinitions(%s)", definitionID),
	}

	apiURL := helpers.BuildApiUrl(environmentHost, fmt.Sprintf("/api/data/%s/environmentvariablevalues", constants.DATAVERSE_API_VERSION), nil)
	resp, err := client.Api.Execute(ctx, nil, "POST", apiURL, nil, requestBody, []int{http.StatusCreated, http.StatusNoContent, http.StatusForbidden, http.StatusNotFound}, nil)
	if err != nil {
		return nil, err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := client.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, err
	}

	entityID := extractEntityIDFromResponse(resp)
	if entityID == "" {
		return nil, fmt.Errorf("no environment variable value id returned after creating schema name '%s'", schemaName)
	}

	return client.getCurrentValueByID(ctx, environmentId, entityID)
}

func (client *client) updateCurrentValue(ctx context.Context, environmentId, valueID, value string) (*environmentVariableValueDto, error) {
	environmentHost, err := client.SolutionClient.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	requestBody := updateEnvironmentVariableValueDto{
		Value: value,
	}

	apiURL := helpers.BuildApiUrl(environmentHost, fmt.Sprintf("/api/data/%s/environmentvariablevalues(%s)", constants.DATAVERSE_API_VERSION, valueID), nil)
	resp, err := client.Api.Execute(ctx, nil, "PATCH", apiURL, nil, requestBody, []int{http.StatusNoContent, http.StatusOK, http.StatusForbidden, http.StatusNotFound}, nil)
	if err != nil {
		return nil, err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := client.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, err
	}

	return client.getCurrentValueByID(ctx, environmentId, valueID)
}

func (client *client) getCurrentValueByID(ctx context.Context, environmentId, valueID string) (*environmentVariableValueDto, error) {
	environmentHost, err := client.SolutionClient.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	values.Add("$select", "environmentvariablevalueid,schemaname,value")
	apiURL := helpers.BuildApiUrl(environmentHost, fmt.Sprintf("/api/data/%s/environmentvariablevalues(%s)", constants.DATAVERSE_API_VERSION, valueID), values)

	response := environmentVariableValueDto{}
	resp, err := client.Api.Execute(ctx, nil, "GET", apiURL, nil, nil, []int{http.StatusOK, http.StatusForbidden, http.StatusNotFound}, &response)
	if err != nil {
		return nil, err
	}
	if err := client.Api.HandleForbiddenResponse(resp); err != nil {
		return nil, err
	}
	if err := client.Api.HandleNotFoundResponse(resp); err != nil {
		return nil, err
	}

	return &response, nil
}

func environmentVariableResourceID(environmentID, schemaName string) string {
	return fmt.Sprintf("%s/%s", environmentID, schemaName)
}

func parseEnvironmentVariableResourceID(id string) (environmentID string, schemaName string, err error) {
	parts := strings.SplitN(id, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid import id %q, expected {environment_id}/{schema_name}", id)
	}

	return parts[0], parts[1], nil
}

func escapeODataStringLiteral(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

func extractEntityIDFromResponse(resp *api.Response) string {
	if resp == nil || resp.HttpResponse == nil {
		return ""
	}

	entityIDHeader := resp.HttpResponse.Header.Get(constants.HEADER_ODATA_ENTITY_ID)
	if entityIDHeader == "" {
		entityIDHeader = resp.HttpResponse.Header.Get(constants.HEADER_LOCATION)
	}
	if entityIDHeader == "" {
		return ""
	}

	re := regexp.MustCompile("[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}")
	matches := re.FindAllString(entityIDHeader, -1)
	if len(matches) == 0 {
		return ""
	}
	return matches[len(matches)-1]
}

func isEnvironmentVariableDefinitionNotFound(err error) bool {
	return err != nil && strings.Contains(err.Error(), "environment variable definition with schema name")
}

func environmentVariableTypeLabel(value int64) string {
	switch value {
	case 100000000:
		return "String"
	case 100000001:
		return "Number"
	case 100000002:
		return "Boolean"
	case 100000003:
		return "JSON"
	case 100000004:
		return "Data Source"
	case 100000005:
		return "Secret"
	default:
		return fmt.Sprintf("%d", value)
	}
}

func environmentVariableSecretStoreLabel(value int64) string {
	switch value {
	case 0:
		return "Azure Key Vault"
	case 1:
		return "Microsoft Dataverse"
	default:
		return fmt.Sprintf("%d", value)
	}
}
