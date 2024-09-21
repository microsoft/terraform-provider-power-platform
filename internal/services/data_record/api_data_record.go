// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package data_record

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

func NewDataRecordClient(apiClient *api.Client) DataRecordClient {
	return DataRecordClient{
		Api: apiClient,
	}
}

type DataRecordClient struct {
	Api *api.Client
}

type EnvironmentIdDto struct {
	Id         string                     `json:"id"`
	Name       string                     `json:"name"`
	Properties EnvironmentIdPropertiesDto `json:"properties"`
}

type EnvironmentIdPropertiesDto struct {
	LinkedEnvironmentMetadata LinkedEnvironmentIdMetadataDto `json:"linkedEnvironmentMetadata"`
}

type LinkedEnvironmentIdMetadataDto struct {
	InstanceURL string
}

type EntityDefinitionsDto struct {
	OdataContext          string `json:"@odata.context"`
	PrimaryIDAttribute    string `json:"PrimaryIdAttribute"`
	LogicalCollectionName string `json:"LogicalCollectionName"`
	MetadataID            string `json:"MetadataId"`
}

type RelationApiResponse struct {
	OdataContext string            `json:"@odata.context"`
	Value        []RelationApiBody `json:"value"`
}

type RelationApiBody struct {
	OdataID string `json:"@odata.id"`
}

func GetEntityDefinition(ctx context.Context, client *DataRecordClient, environmentId, entityLogicalName string) (*EntityDefinitionsDto, error) {
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

	entityDefinition := EntityDefinitionsDto{}
	_, err = client.Api.Execute(ctx, nil, "GET", entityDefinitionApiUrl.String(), nil, nil, []int{http.StatusOK}, &entityDefinition)
	if err != nil {
		return nil, err
	}

	return &entityDefinition, nil
}

func (client *DataRecordClient) GetEnvironmentHostById(ctx context.Context, environmentId string) (string, error) {
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

func (client *DataRecordClient) getEnvironment(ctx context.Context, environmentId string) (*EnvironmentIdDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	env := EnvironmentIdDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}

func (client *DataRecordClient) GetDataRecordsByODataQuery(ctx context.Context, environmentId, query string, headers map[string]string) (*ODataQueryResponse, error) {
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
			return nil, fmt.Errorf("value field is not of type []any")
		}
		for _, item := range valueSlice {
			value, ok := item.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("item is not of type map[string]any")
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

func (client *DataRecordClient) GetDataRecord(ctx context.Context, recordId, environmentId, tableName string) (map[string]any, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	entityDefinition, err := GetEntityDefinition(ctx, client, environmentId, tableName)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/%s/%s(%s)", constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName, recordId),
	}

	result := make(map[string]any, 0)

	_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &result)
	if err != nil {
		if strings.ContainsAny(err.Error(), "404") {
			return nil, helpers.WrapIntoProviderError(err, helpers.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("Data Record '%s' not found", recordId))
		}
		return nil, err
	}

	return result, nil
}

func (client *DataRecordClient) GetRelationData(ctx context.Context, environmentId, tableName, recordId, relationName string) ([]any, error) {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	entityDefinition, err := GetEntityDefinition(ctx, client, environmentId, tableName)
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
		return nil, fmt.Errorf("value field not found in result when retrieving relational data")
	}

	value, ok := field.([]any)
	if !ok {
		return nil, fmt.Errorf("value field is not of type []any in relational data")
	}

	return value, nil
}

func (client *DataRecordClient) GetTableSingularNameFromPlural(ctx context.Context, environmentId, logicalCollectionName string) (*string, error) {
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
		return nil, fmt.Errorf("logicalName field not found in result when retrieving table singular name")
	}
	return &result, nil
}

func getEntityRelationDefinitionOneToMany(mapResponse map[string]any, entityLogicalName, relationLogicalName string) (string, error) {
	var tableName string
	oneToMany, ok := mapResponse["OneToManyRelationships"].([]any)
	if !ok {
		return "", fmt.Errorf("OneToManyRelationships field is not of type []any")
	}
	for _, list := range oneToMany {
		item, ok := list.(map[string]any)
		if !ok {
			return "", fmt.Errorf("item is not of type map[string]any")
		}
		if item["ReferencingEntityNavigationPropertyName"] == relationLogicalName && item["Entity1LogicalName"] != entityLogicalName {
			var ok bool
			tableName, ok = item["ReferencedEntity"].(string)
			if !ok {
				return "", fmt.Errorf("ReferencedEntity field is not of type string")
			}
			break
		}
		if item["ReferencedEntityNavigationPropertyName"] == relationLogicalName && item["Entity2LogicalName"] != entityLogicalName {
			var ok bool
			tableName, ok = item["ReferencingEntity"].(string)
			if !ok {
				return "", fmt.Errorf("ReferencedEntity field is not of type string")
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
		return "", fmt.Errorf("ManyToOneRelationships field is not of type []any")
	}
	for _, list := range manyToOne {
		item, ok := list.(map[string]any)
		if !ok {
			return "", fmt.Errorf("item is not of type map[string]any")
		}
		if item["ReferencingEntityNavigationPropertyName"] == relationLogicalName {
			var ok bool
			tableName, ok = item["ReferencedEntity"].(string)
			if !ok {
				return "", fmt.Errorf("ReferencedEntity field is not of type string")
			}
			break
		}
		if item["ReferencedEntityNavigationPropertyName"] == relationLogicalName {
			var ok bool
			tableName, ok = item["ReferencingEntity"].(string)
			if !ok {
				return "", fmt.Errorf("ReferencedEntity field is not of type string")
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
		return "", fmt.Errorf("ManyToManyRelationships field is not of type []any")
	}
	for _, list := range manyToMany {
		item, ok := list.(map[string]any)
		if !ok {
			return "", fmt.Errorf("item is not of type map[string]any")
		}
		if item["Entity1NavigationPropertyName"] == relationLogicalName && item["Entity1LogicalName"] != entityLogicalName {
			tableName, ok = item["Entity1LogicalName"].(string)
			if !ok {
				return "", fmt.Errorf("Entity1LogicalName field is not of type string")
			}
			break
		}
		if item["Entity2NavigationPropertyName"] == relationLogicalName && item["Entity2LogicalName"] != entityLogicalName {
			tableName, ok = item["Entity2LogicalName"].(string)
			if !ok {
				return "", fmt.Errorf("Entity2LogicalName field is not of type string")
			}
			break
		}
	}
	return tableName, nil
}

func (client *DataRecordClient) GetEntityRelationDefinitionInfo(ctx context.Context, environmentId string, entityLogicalName string, relationLogicalName string) (tableName string, err error) {
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

func (client *DataRecordClient) ApplyDataRecord(ctx context.Context, recordId, environmentId, tableName string, columns map[string]any) (*DataRecordDto, error) {
	result := DataRecordDto{}

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

				entityDefinition, err := GetEntityDefinition(ctx, client, environmentId, tableLogicalName)
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

	entityDefinition, err := GetEntityDefinition(ctx, client, environmentId, tableName)
	if err != nil {
		return nil, err
	}

	method := "POST"
	apiPath := fmt.Sprintf("/api/data/%s/%s", constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName)

	if val, ok := columns[entityDefinition.PrimaryIDAttribute]; ok {
		method = "PATCH"
		apiPath = fmt.Sprintf("%s(%s)", apiPath, val)
	} else if recordId != "" {
		method = "PATCH"
		apiPath = fmt.Sprintf("%s(%s)", apiPath, recordId)
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   apiPath,
	}

	response, err := client.Api.Execute(ctx, nil, method, apiUrl.String(), nil, columns, []int{http.StatusOK, http.StatusNoContent}, nil)
	if err != nil {
		return nil, err
	}

	if len(response.BodyAsBytes) != 0 {
		err = json.Unmarshal(response.BodyAsBytes, &result)
		if err != nil {
			return nil, err
		}
	} else if response.Response.Header.Get(constants.HEADER_ODATA_ENTITY_ID) != "" {
		re := regexp.MustCompile("[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}")
		match := re.FindAllStringSubmatch(response.Response.Header.Get(constants.HEADER_ODATA_ENTITY_ID), -1)
		if len(match) == 0 {
			return nil, fmt.Errorf("no entity record id returned from the odata-entityid header")
		}
		result.Id = match[len(match)-1][0]
	} else {
		return nil, fmt.Errorf("no entity record id returned from the API")
	}

	err = applyRelations(ctx, client, relations, environmentId, result.Id, entityDefinition)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *DataRecordClient) DeleteDataRecord(ctx context.Context, recordId string, environmentId string, tableName string, columns map[string]any) error {
	environmentHost, err := client.GetEnvironmentHostById(ctx, environmentId)
	if err != nil {
		return err
	}

	tableEntityDefinition, err := GetEntityDefinition(ctx, client, environmentId, tableName)
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
					return fmt.Errorf("nestedItem is not of type map[string]any")
				}

				dataRecordId, ok := nestedMap["data_record_id"].(string)
				if !ok {
					return fmt.Errorf("data_record_id field is missing or not a string")
				}

				apiUrl := &url.URL{
					Scheme: constants.HTTPS,
					Host:   environmentHost,
					Path:   fmt.Sprintf("/api/data/%s/%s(%s)/%s(%s)/$ref", constants.DATAVERSE_API_VERSION, tableEntityDefinition.LogicalCollectionName, recordId, key, dataRecordId),
				}
				_, err = client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNoContent}, nil)
				if err != nil && !strings.ContainsAny(err.Error(), "404") {
					return err
				}
			}
		}
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/%s/%s(%s)", constants.DATAVERSE_API_VERSION, tableEntityDefinition.LogicalCollectionName, recordId),
	}
	_, err = client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, columns, []int{http.StatusOK, http.StatusNoContent}, nil)
	if err != nil && !strings.ContainsAny(err.Error(), "404") {
		// TODO: 404 is desired state for delete.  We should pass 404 as acceptable status code and not error
		return err
	}
	return nil
}

func getTableLogicalNameAndDataRecordIdFromMap(nestedMap map[string]any) (tableLogicalName string, dataRecordId string, err error) {
	tableLogicalName, ok := nestedMap["table_logical_name"].(string)
	if !ok {
		return "", "", fmt.Errorf("table_logical_name field is missing or not a string")
	}
	id, ok := nestedMap["data_record_id"].(string)
	if !ok {
		return "", "", fmt.Errorf("data_record_id field is missing or not a string")
	}
	return tableLogicalName, id, nil
}

func applyRelations(ctx context.Context, client *DataRecordClient, relations map[string]any, environmentId string, parentRecordId string, entityDefinition *EntityDefinitionsDto) error {
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

func applyRelation(ctx context.Context, environmentHost string, entityDefinition *EntityDefinitionsDto, parentRecordId string, key string, client *DataRecordClient, nestedMapList []any, environmentId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   environmentHost,
		Path:   fmt.Sprintf("/api/data/%s/%s(%s)/%s/$ref", constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName, parentRecordId, key),
	}

	existingRelationsResponse := RelationApiResponse{}

	apiResponse, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNoContent}, nil)
	if err != nil {
		return err
	}

	// TODO: execute will unmarshal the response into the existingRelationsResponse use that instead
	err = json.Unmarshal(apiResponse.BodyAsBytes, &existingRelationsResponse)
	if err != nil {
		return err
	}

	var toBeDeleted = make([]RelationApiBody, 0)

	for _, existingRelation := range existingRelationsResponse.Value {
		shouldDelete := true
		for _, nestedItem := range nestedMapList {
			nestedMap, ok := nestedItem.(map[string]any)
			if !ok {
				return fmt.Errorf("nestedItem is not of type map[string]any")
			}

			tableLogicalName, dataRecordId, err := getTableLogicalNameAndDataRecordIdFromMap(nestedMap)
			if err != nil {
				return err
			}

			relationEntityDefinition, err := GetEntityDefinition(ctx, client, environmentId, tableLogicalName)
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
			return fmt.Errorf("nestedItem is not of type map[string]any")
		}

		tableLogicalName, dataRecordId, err := getTableLogicalNameAndDataRecordIdFromMap(nestedMap)
		if err != nil {
			return err
		}

		entityDefinition, err := GetEntityDefinition(ctx, client, environmentId, tableLogicalName)
		if err != nil {
			return err
		}

		relation := RelationApiBody{
			OdataID: fmt.Sprintf("https://%s/api/data/%s/%s(%s)", environmentHost, constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName, dataRecordId),
		}
		_, err = client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, relation, []int{http.StatusOK, http.StatusNoContent}, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
