// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package data_record

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/customerrors"
)

func newDataRecordClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}

type client struct {
	Api *api.Client
}

func getEntityDefinition(ctx context.Context, client *client, environmentId, entityLogicalName string) (*entityDefinitionsDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	entityDefinitionApiUrl := &url.URL{
		Scheme:   constants.HTTPS,
		Host:     environmentHost,
		Path:     fmt.Sprintf("/api/data/%s/EntityDefinitions(LogicalName='%s')", constants.DATAVERSE_API_VERSION, entityLogicalName),
		Fragment: "$select=PrimaryIdAttribute,LogicalCollectionName",
	}

	entityDefinition := entityDefinitionsDto{}
	_, err = client.Api.Execute(ctx, nil, "GET", entityDefinitionApiUrl.String(), nil, nil, []int{http.StatusOK}, &entityDefinition)
	if err != nil {
		return nil, err
	}

	return &entityDefinition, nil
}

func (client *client) GetEnvironmentHostById(ctx context.Context, environmentId string) (string, error) {
	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return "", err
	}
	environmentUrl := strings.TrimSuffix(env.Properties.LinkedEnvironmentMetadata.InstanceURL, "/")
	if environmentUrl == "" {
		return "", customerrors.WrapIntoProviderError(nil, customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND, "environment url not found, please check if the environment has dataverse linked")
	}

	envUrl, err := url.Parse(environmentUrl)
	if err != nil {
		return "", err
	}
	return envUrl.Host, nil
}

func (client *client) getEnvironment(ctx context.Context, environmentId string) (*environmentIdDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	env := environmentIdDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}

func (client *client) GetDataRecordsByODataQuery(ctx context.Context, environmentId, query string, headers map[string]string) (*ODataQueryResponse, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	var h = make(http.Header)
	for k, v := range headers {
		h.Add(k, v)
	}

	apiUrl := fmt.Sprintf("https://%s/api/data/%s/%s", environmentHost, constants.DATAVERSE_API_VERSION, query)

	response := map[string]any{}
	_, err = client.Api.Execute(ctx, nil, "GET", apiUrl, h, nil, []int{http.StatusOK}, &response)
	if err != nil {
		return nil, err
	}

	var totalRecords *int64
	if response["@Microsoft.Dynamics.CRM.totalrecordcount"] != nil {
		count := int64(response["@Microsoft.Dynamics.CRM.totalrecordcount"].(float64))
		totalRecords = &count
	}
	var totalRecordsCountLimitExceeded *bool
	if val, ok := response["@Microsoft.Dynamics.CRM.totalrecordcountlimitexceeded"].(bool); ok {
		isLimitExceeded := val
		totalRecordsCountLimitExceeded = &isLimitExceeded
	}

	records := []map[string]any{}
	if response["value"] != nil {
		valueSlice, ok := response["value"].([]any)
		if !ok {
			return nil, errors.New("value field is not of type []any")
		}
		for _, item := range valueSlice {
			value, ok := item.(map[string]any)
			if !ok {
				return nil, errors.New("item is not of type map[string]any")
			}
			records = append(records, value)
		}
	} else {
		records = append(records, response)
	}

	pluralName := strings.Split(response["@odata.context"].(string), "#")[1]
	if index := strings.IndexAny(pluralName, "(/"); index != -1 {
		pluralName = pluralName[:index]
	}

	return &ODataQueryResponse{
		Records:                  records,
		TotalRecord:              totalRecords,
		TotalRecordLimitExceeded: totalRecordsCountLimitExceeded,
		TableMetadataUrl:         response["@odata.context"].(string),
		// url will be as example: https://org.crm4.dynamics.com/api/data/v9.2/$metadata#tablepluralname/$entity.
		TablePluralName: pluralName,
	}, nil
}

type ODataQueryResponse struct {
	Records                  []map[string]any
	TotalRecord              *int64
	TotalRecordLimitExceeded *bool
	TableMetadataUrl         string
	TablePluralName          string
}

func (client *client) GetDataRecord(ctx context.Context, recordId, environmentId, tableName string) (map[string]any, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	entityDefinition, err := getEntityDefinition(ctx, client, environmentId, tableName)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/%s/%s(%s)", constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName, recordId),
	}

	result := make(map[string]any, 0)

	_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNotFound}, &result)
	if err != nil {
		var httpError *customerrors.UnexpectedHttpStatusCodeError
		if errors.As(err, &httpError) && httpError.StatusCode == http.StatusNotFound {
			return nil, customerrors.WrapIntoProviderError(err, customerrors.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("Data Record '%s' not found", recordId))
		}
		return nil, err
	}

	return result, nil
}

func (client *client) GetRelationData(ctx context.Context, environmentId, tableName, recordId, relationName string) ([]any, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	entityDefinition, err := getEntityDefinition(ctx, client, environmentId, tableName)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme:   constants.HTTPS,
		Host:     environmentHost,
		Path:     fmt.Sprintf("/api/data/%s/%s(%s)/%s", constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName, recordId, relationName),
		RawQuery: "$select=createdon",
	}

	result := make(map[string]any, 0)

	_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &result)
	if err != nil {
		return nil, err
	}

	field, ok := result["value"]
	if !ok {
		return nil, errors.New("value field not found in result when retrieving relational data")
	}

	value, ok := field.([]any)
	if !ok {
		return nil, errors.New("value field is not of type []any in relational data")
	}

	return value, nil
}

func (client *client) GetTableSingularNameFromPlural(ctx context.Context, environmentId, logicalCollectionName string) (*string, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/%s/EntityDefinitions", constants.DATAVERSE_API_VERSION),
	}
	q := apiUrl.Query()
	q.Add("$filter", fmt.Sprintf("LogicalCollectionName eq '%s'", logicalCollectionName))
	q.Add("$select", "PrimaryIdAttribute,LogicalCollectionName,LogicalName")
	apiUrl.RawQuery = q.Encode()

	response, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return nil, err
	}

	var mapResponse map[string]any
	err = json.Unmarshal(response.BodyAsBytes, &mapResponse)
	if err != nil {
		return nil, err
	}

	var result string
	if mapResponse["value"] != nil && len(mapResponse["value"].([]any)) > 0 {
		if value, ok := mapResponse["value"].([]any)[0].(map[string]any); ok {
			if logicalName, ok := value["LogicalName"].(string); ok {
				result = logicalName
			}
		}
	} else if logicalName, ok := mapResponse["LogicalName"].(string); ok {
		result = logicalName
	} else {
		return nil, errors.New("logicalName field not found in result when retrieving table singular name")
	}
	return &result, nil
}

func getEntityRelationDefinitionOneToMany(mapResponse map[string]any, entityLogicalName, relationLogicalName string) (string, error) {
	var tableName string
	oneToMany, ok := mapResponse["OneToManyRelationships"].([]any)
	if !ok {
		return "", errors.New("OneToManyRelationships field is not of type []any")
	}
	for _, list := range oneToMany {
		item, ok := list.(map[string]any)
		if !ok {
			return "", errors.New("item is not of type map[string]any")
		}
		if item["ReferencingEntityNavigationPropertyName"] == relationLogicalName && item["Entity1LogicalName"] != entityLogicalName {
			var ok bool
			tableName, ok = item["ReferencedEntity"].(string)
			if !ok {
				return "", errors.New("ReferencedEntity field is not of type string")
			}
			break
		}
		if item["ReferencedEntityNavigationPropertyName"] == relationLogicalName && item["Entity2LogicalName"] != entityLogicalName {
			var ok bool
			tableName, ok = item["ReferencingEntity"].(string)
			if !ok {
				return "", errors.New("ReferencedEntity field is not of type string")
			}
			break
		}
	}
	return tableName, nil
}

func getEntityRelationDefinitionManyToOne(mapResponse map[string]any, relationLogicalName string) (string, error) {
	var tableName string
	manyToOne, ok := mapResponse["ManyToOneRelationships"].([]any)
	if !ok {
		return "", errors.New("ManyToOneRelationships field is not of type []any")
	}
	for _, list := range manyToOne {
		item, ok := list.(map[string]any)
		if !ok {
			return "", errors.New("item is not of type map[string]any")
		}
		if item["ReferencingEntityNavigationPropertyName"] == relationLogicalName {
			var ok bool
			tableName, ok = item["ReferencedEntity"].(string)
			if !ok {
				return "", errors.New("ReferencedEntity field is not of type string")
			}
			break
		}
		if item["ReferencedEntityNavigationPropertyName"] == relationLogicalName {
			var ok bool
			tableName, ok = item["ReferencingEntity"].(string)
			if !ok {
				return "", errors.New("ReferencedEntity field is not of type string")
			}
			break
		}
	}
	return tableName, nil
}

func getEntityRelationDefinitionManyToMany(mapResponse map[string]any, entityLogicalName, relationLogicalName string) (string, error) {
	var tableName string
	manyToMany, ok := mapResponse["ManyToManyRelationships"].([]any)
	if !ok {
		return "", errors.New("ManyToManyRelationships field is not of type []any")
	}
	for _, list := range manyToMany {
		item, ok := list.(map[string]any)
		if !ok {
			return "", errors.New("item is not of type map[string]any")
		}
		if item["Entity1NavigationPropertyName"] == relationLogicalName && item["Entity1LogicalName"] != entityLogicalName {
			tableName, ok = item["Entity1LogicalName"].(string)
			if !ok {
				return "", errors.New("Entity1LogicalName field is not of type string")
			}
			break
		}
		if item["Entity2NavigationPropertyName"] == relationLogicalName && item["Entity2LogicalName"] != entityLogicalName {
			tableName, ok = item["Entity2LogicalName"].(string)
			if !ok {
				return "", errors.New("Entity2LogicalName field is not of type string")
			}
			break
		}
	}
	return tableName, nil
}

func (client *client) GetEntityAttributesDefinition(ctx context.Context, environmentId string, entityLogicalName string) ([]attributesApiBodyDto, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	apiUrl := fmt.Sprintf("https://%s/api/data/%s/EntityDefinitions(LogicalName='%s')/Attributes?$select=LogicalName", environmentHost, constants.DATAVERSE_API_VERSION, entityLogicalName)

	results := attributesApiResponseDto{}
	_, err = client.Api.Execute(ctx, nil, "GET", apiUrl, nil, nil, []int{http.StatusOK}, &results)
	if err != nil {
		return nil, err
	}
	return results.Value, nil
}

func (client *client) GetEntityRelationDefinitionInfo(ctx context.Context, environmentId string, entityLogicalName string, relationLogicalName string) (tableName string, err error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return "", err
	}

	apiUrl := fmt.Sprintf("https://%s/api/data/%s/EntityDefinitions(LogicalName='%s')?$expand=OneToManyRelationships,ManyToManyRelationships,ManyToOneRelationships", environmentHost, constants.DATAVERSE_API_VERSION, entityLogicalName)

	response, err := client.Api.Execute(ctx, nil, "GET", apiUrl, nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return "", err
	}

	var mapResponse map[string]any
	err = json.Unmarshal(response.BodyAsBytes, &mapResponse)
	if err != nil {
		return "", err
	}

	tableName, err = getEntityRelationDefinitionOneToMany(mapResponse, entityLogicalName, relationLogicalName)
	if err != nil {
		return "", err
	}
	if tableName != "" {
		return tableName, nil
	}

	tableName, err = getEntityRelationDefinitionManyToOne(mapResponse, relationLogicalName)
	if err != nil {
		return "", err
	}
	if tableName != "" {
		return tableName, nil
	}

	tableName, err = getEntityRelationDefinitionManyToMany(mapResponse, entityLogicalName, relationLogicalName)
	if err != nil {
		return "", err
	}

	return tableName, nil
}

func (client *client) ApplyDataRecord(ctx context.Context, recordId, environmentId, tableName string, columns map[string]any) (*dataRecordDto, error) {
	result := dataRecordDto{}

	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	relations := make(map[string]any, 0)

	for key, value := range columns {
		if nestedMap, ok := value.(map[string]any); ok {
			delete(columns, key)
			if len(nestedMap) > 0 {
				tableLogicalName, dataRecordId, err := getTableLogicalNameAndDataRecordIdFromMap(nestedMap)
				if err != nil {
					return nil, err
				}

				entityDefinition, err := getEntityDefinition(ctx, client, environmentId, tableLogicalName)
				if err != nil {
					return nil, err
				}

				columns[fmt.Sprintf("%s@odata.bind", key)] = fmt.Sprintf("/%s(%s)", entityDefinition.LogicalCollectionName, dataRecordId)
			}
		} else if nestedMapList, ok := value.([]any); ok {
			delete(columns, key)
			relations[key] = nestedMapList
		}
	}

	entityDefinition, err := getEntityDefinition(ctx, client, environmentId, tableName)
	if err != nil {
		return nil, err
	}

	// we will send create operation as default.
	method := "POST"
	apiPath := fmt.Sprintf("/api/data/%s/%s", constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName)

	if val, ok := columns[entityDefinition.PrimaryIDAttribute]; ok {
		// if one of the sent attributes is the primaryId then send an update.
		method = "PATCH"
		apiPath = fmt.Sprintf("%s(%s)", apiPath, val)
	} else if recordId != "" {
		// if we are referencing the record by its primaryId and its not empty (update or delete) then we send an update
		method = "PATCH"
		apiPath = fmt.Sprintf("%s(%s)", apiPath, recordId)
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   apiPath,
	}

	response, err := client.Api.Execute(ctx, nil, method, apiUrl.String(), nil, columns, []int{http.StatusOK, http.StatusNoContent, http.StatusPreconditionFailed}, nil)
	if err != nil {
		return nil, err
	}
	if response.HttpResponse.StatusCode == http.StatusPreconditionFailed {
		// if record was already found then we will try to update it.
		response, err = client.Api.Execute(ctx, nil, "PATCH", apiUrl.String(), nil, columns, []int{http.StatusOK, http.StatusNoContent}, nil)
		// if the PATCH URL is not pointing to specific recording using PrimaryIDAttribute then we have to inform a user about it.
		if response.HttpResponse.StatusCode == http.StatusMethodNotAllowed {
			return nil, errors.New("record already exists. To update an existing record, primaryId must be provided in the columns attribute")
		}
		if err != nil {
			return nil, err
		}
	}

	if len(response.BodyAsBytes) != 0 {
		err = json.Unmarshal(response.BodyAsBytes, &result)
		if err != nil {
			return nil, err
		}
	} else if response.HttpResponse.Header.Get(constants.HEADER_ODATA_ENTITY_ID) != "" {
		re := regexp.MustCompile("[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}")
		match := re.FindAllStringSubmatch(response.HttpResponse.Header.Get(constants.HEADER_ODATA_ENTITY_ID), -1)
		if len(match) == 0 {
			return nil, errors.New("no entity record id returned from the odata-entityid header")
		}
		result.Id = match[len(match)-1][0]
	} else {
		return nil, errors.New("no entity record id returned from the API")
	}

	err = applyRelations(ctx, client, relations, environmentId, result.Id, entityDefinition)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *client) DeleteDataRecord(ctx context.Context, recordId string, environmentId string, tableName string, columns map[string]any) error {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return err
	}

	tableEntityDefinition, err := getEntityDefinition(ctx, client, environmentId, tableName)
	if err != nil {
		return err
	}

	for key, value := range columns {
		if _, ok := value.(map[string]any); ok {
			delete(columns, key)
		}
		if nestedMapList, ok := value.([]any); ok {
			delete(columns, key)

			for _, nestedItem := range nestedMapList {
				nestedMap, ok := nestedItem.(map[string]any)
				if !ok {
					return errors.New("nestedItem is not of type map[string]any")
				}

				dataRecordId, ok := nestedMap["data_record_id"].(string)
				if !ok {
					return errors.New("data_record_id field is missing or not a string")
				}

				apiUrl := &url.URL{
					Scheme: constants.HTTPS,
					Host:   environmentHost,
					Path:   fmt.Sprintf("/api/data/%s/%s(%s)/%s(%s)/$ref", constants.DATAVERSE_API_VERSION, tableEntityDefinition.LogicalCollectionName, recordId, key, dataRecordId),
				}
				_, err = client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNoContent, http.StatusNotFound}, nil)
				if err != nil {
					return errors.New("error while deleting data record. %w")
				}
			}
		}
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/%s/%s(%s)", constants.DATAVERSE_API_VERSION, tableEntityDefinition.LogicalCollectionName, recordId),
	}

	// 200, 201, or 404 are acceptable status codes for delete and not error
	_, err = client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, columns, []int{http.StatusOK, http.StatusNoContent, http.StatusNotFound}, nil)
	if err != nil {
		return err
	}
	return nil
}

func getTableLogicalNameAndDataRecordIdFromMap(nestedMap map[string]any) (tableLogicalName string, dataRecordId string, err error) {
	tableLogicalName, ok := nestedMap["table_logical_name"].(string)
	if !ok {
		return "", "", errors.New("table_logical_name field is missing or not a string")
	}
	id, ok := nestedMap["data_record_id"].(string)
	if !ok {
		return "", "", errors.New("data_record_id field is missing or not a string")
	}
	return tableLogicalName, id, nil
}

func applyRelations(ctx context.Context, client *client, relations map[string]any, environmentId string, parentRecordId string, entityDefinition *entityDefinitionsDto) error {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return err
	}

	for key, value := range relations {
		if nestedMapList, ok := value.([]any); ok {
			err := applyRelation(ctx, environmentHost, entityDefinition, parentRecordId, key, client, nestedMapList, environmentId)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func applyRelation(ctx context.Context, environmentHost string, entityDefinition *entityDefinitionsDto, parentRecordId string, key string, client *client, nestedMapList []any, environmentId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/%s/%s(%s)/%s/$ref", constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName, parentRecordId, key),
	}

	existingRelationsResponse := relationApiResponseDto{}

	apiResponse, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNoContent}, nil)
	if err != nil {
		return err
	}

	// TODO: execute will unmarshal the response into the existingRelationsResponse use that instead
	err = json.Unmarshal(apiResponse.BodyAsBytes, &existingRelationsResponse)
	if err != nil {
		return err
	}

	var toBeDeleted = make([]relationApiBodyDto, 0)

	for _, existingRelation := range existingRelationsResponse.Value {
		shouldDelete := true
		for _, nestedItem := range nestedMapList {
			nestedMap, ok := nestedItem.(map[string]any)
			if !ok {
				return errors.New("nestedItem is not of type map[string]any")
			}

			tableLogicalName, dataRecordId, err := getTableLogicalNameAndDataRecordIdFromMap(nestedMap)
			if err != nil {
				return err
			}

			relationEntityDefinition, err := getEntityDefinition(ctx, client, environmentId, tableLogicalName)
			if err != nil {
				return err
			}
			if existingRelation.OdataID == fmt.Sprintf("https://%s/api/data/%s/%s(%s)", environmentHost, constants.DATAVERSE_API_VERSION, relationEntityDefinition.LogicalCollectionName, dataRecordId) {
				shouldDelete = false
				break
			}
		}
		if shouldDelete {
			toBeDeleted = append(toBeDeleted, existingRelation)
		}
	}

	for _, relation := range toBeDeleted {
		_, err = client.Api.Execute(ctx, nil, "DELETE", relation.OdataID, nil, nil, []int{http.StatusOK, http.StatusNoContent}, nil)
		if err != nil {
			return err
		}
	}

	for _, nestedItem := range nestedMapList {
		nestedMap, ok := nestedItem.(map[string]any)
		if !ok {
			return errors.New("nestedItem is not of type map[string]any")
		}

		tableLogicalName, dataRecordId, err := getTableLogicalNameAndDataRecordIdFromMap(nestedMap)
		if err != nil {
			return err
		}

		entityDefinition, err := getEntityDefinition(ctx, client, environmentId, tableLogicalName)
		if err != nil {
			return err
		}

		relation := relationApiBodyDto{
			OdataID: fmt.Sprintf("https://%s/api/data/%s/%s(%s)", environmentHost, constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName, dataRecordId),
		}
		_, err = client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, relation, []int{http.StatusOK, http.StatusNoContent}, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
