// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"encoding/json"
	"fmt"
	http "net/http"
	"net/url"
	"regexp"
	"strings"

	constants "github.com/microsoft/terraform-provider-power-platform/constants"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
)

func NewDataRecordClient(api *api.ApiClient) DataRecordClient {
	return DataRecordClient{
		Api: api,
	}
}

type DataRecordClient struct {
	Api *api.ApiClient
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
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	e, err := url.Parse(environmentUrl)
	if err != nil {
		return nil, err
	}

	entityDefinitionApiUrl := &url.URL{
		Scheme:   e.Scheme,
		Host:     e.Host,
		Path:     fmt.Sprintf("/api/data/%s/EntityDefinitions(LogicalName='%s')", constants.DATAVERSE_API_VERSION, entityLogicalName),
		Fragment: "$select=PrimaryIdAttribute,LogicalCollectionName",
	}

	entityDefinition := EntityDefinitionsDto{}
	_, err = client.Api.Execute(ctx, "GET", entityDefinitionApiUrl.String(), nil, nil, []int{http.StatusOK}, &entityDefinition)
	if err != nil {
		return nil, err
	}

	return &entityDefinition, nil
}

func (client *DataRecordClient) GetEnvironmentUrlById(ctx context.Context, environmentId string) (string, error) {
	env, err := client.getEnvironment(ctx, environmentId)
	if err != nil {
		return "", err
	}
	environmentUrl := strings.TrimSuffix(env.Properties.LinkedEnvironmentMetadata.InstanceURL, "/")
	return environmentUrl, nil
}

func (client *DataRecordClient) getEnvironment(ctx context.Context, environmentId string) (*EnvironmentIdDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/%s", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	apiUrl.RawQuery = values.Encode()

	env := EnvironmentIdDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &env)
	if err != nil {
		return nil, err
	}

	return &env, nil
}

func (client *DataRecordClient) GetDataRecordsByODataQuery(ctx context.Context, environmentId, query string, headers map[string]string) (*ODataQueryResponse, error) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	var h http.Header = make(http.Header)
	for k, v := range headers {
		h.Add(k, v)
	}

	e, _ := url.Parse(environmentUrl)
	apiUrl := fmt.Sprintf("%s://%s/api/data/%s/%s", e.Scheme, e.Host, constants.DATAVERSE_API_VERSION, query)

	response := map[string]interface{}{}
	_, err = client.Api.Execute(ctx, "GET", apiUrl, h, nil, []int{http.StatusOK}, &response)
	if err != nil {
		return nil, err
	}

	var totalRecords *int64 = nil
	if response["@Microsoft.Dynamics.CRM.totalrecordcount"] != nil {
		count := int64(response["@Microsoft.Dynamics.CRM.totalrecordcount"].(float64))
		totalRecords = &count
	}
	var totalRecordsCountLimitExceeded *bool = nil
	if response["@Microsoft.Dynamics.CRM.totalrecordcountlimitexceeded"] != nil {
		isLimitExceeded := response["@Microsoft.Dynamics.CRM.totalrecordcountlimitexceeded"].(bool)
		totalRecordsCountLimitExceeded = &isLimitExceeded
	}

	records := []map[string]interface{}{}
	if response["value"] != nil {
		for _, item := range response["value"].([]interface{}) {
			value := item.(map[string]interface{})
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
		//url will be as example: https://org.crm4.dynamics.com/api/data/v9.2/$metadata#tablepluralname/$entity
		TablePluralName: pluralName,
	}, nil
}

type ODataQueryResponse struct {
	Records                  []map[string]interface{}
	TotalRecord              *int64
	TotalRecordLimitExceeded *bool
	TableMetadataUrl         string
	TablePluralName          string
}

func (client *DataRecordClient) GetDataRecord(ctx context.Context, recordId, environmentId, tableName string) (map[string]interface{}, error) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	entityDefinition, err := GetEntityDefinition(ctx, client, environmentId, tableName)
	if err != nil {
		return nil, err
	}

	e, err := url.Parse(environmentUrl)
	if err != nil {
		return nil, err
	}
	apiUrl := &url.URL{
		Scheme: e.Scheme,
		Host:   e.Host,
		Path:   fmt.Sprintf("/api/data/%s/%s(%s)", constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName, recordId),
	}

	result := make(map[string]interface{}, 0)

	_, err = client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &result)
	if err != nil {
		if strings.ContainsAny(err.Error(), "404") {
			return nil, powerplatform_helpers.WrapIntoProviderError(err, powerplatform_helpers.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("Data Record '%s' not found", recordId))
		}
		return nil, err
	}

	return result, nil
}

func (client *DataRecordClient) GetRelationData(ctx context.Context, environmentId, tableName, recordId, relationName string) ([]interface{}, error) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	entityDefinition, err := GetEntityDefinition(ctx, client, environmentId, tableName)
	if err != nil {
		return nil, err
	}

	e, err := url.Parse(environmentUrl)
	if err != nil {
		return nil, err
	}
	apiUrl := &url.URL{
		Scheme:   e.Scheme,
		Host:     e.Host,
		Path:     fmt.Sprintf("/api/data/%s/%s(%s)/%s", constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName, recordId, relationName),
		RawQuery: "$select=createdon",
	}

	result := make(map[string]interface{}, 0)

	_, err = client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &result)
	if err != nil {
		return nil, err
	}

	field, ok := result["value"]
	if !ok {
		return nil, fmt.Errorf("value field not found in result when retrieving relational data")
	}

	value, ok := field.([]interface{})
	if !ok {
		return nil, fmt.Errorf("value field is not of type []interface{} in relational data")
	}

	return value, nil
}

func (client *DataRecordClient) GetTableSingularNameFromPlural(ctx context.Context, environmentId, logicalCollectionName string) (*string, error) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	e, err := url.Parse(environmentUrl)
	if err != nil {
		return nil, err
	}
	apiUrl := &url.URL{
		Scheme: e.Scheme,
		Host:   e.Host,
		Path:   fmt.Sprintf("/api/data/%s/EntityDefinitions", constants.DATAVERSE_API_VERSION),
	}
	q := apiUrl.Query()
	q.Add("$filter", fmt.Sprintf("LogicalCollectionName eq '%s'", logicalCollectionName))
	q.Add("$select", "PrimaryIdAttribute,LogicalCollectionName,LogicalName")
	apiUrl.RawQuery = q.Encode()

	response, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return nil, err
	}

	var mapResponse map[string]interface{}
	json.Unmarshal(response.BodyAsBytes, &mapResponse)

	var result string
	if mapResponse["value"] != nil && len(mapResponse["value"].([]interface{})) > 0 &&
		mapResponse["value"].([]interface{})[0].(map[string]interface{}) != nil &&
		mapResponse["value"].([]interface{})[0].(map[string]interface{})["LogicalName"] != nil {
		result = mapResponse["value"].([]interface{})[0].(map[string]interface{})["LogicalName"].(string)
	} else if mapResponse["LogicalName"] != nil {
		result = mapResponse["LogicalName"].(string)
	} else {
		return nil, fmt.Errorf("logicalName field not found in result when retrieving table singular name")
	}
	return &result, nil
}

func (client *DataRecordClient) GetEntityRelationDefinitionInfo(ctx context.Context, environmentId string, entityLogicalName string, relationLogicalName string) (tableName string, err error) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return "", err
	}

	apiUrl := fmt.Sprintf("%s/api/data/%s/EntityDefinitions(LogicalName='%s')?$expand=OneToManyRelationships,ManyToManyRelationships,ManyToOneRelationships", environmentUrl, constants.DATAVERSE_API_VERSION, entityLogicalName)

	response, err := client.Api.Execute(ctx, "GET", apiUrl, nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return "", err
	}

	var mapResponse map[string]interface{}
	json.Unmarshal(response.BodyAsBytes, &mapResponse)

	oneToMany, ok := mapResponse["OneToManyRelationships"].([]interface{})
	if !ok {
		return "", fmt.Errorf("OneToManyRelationships field is not of type []interface{}")
	}
	for _, list := range oneToMany {
		item := list.(map[string]interface{})
		if item["ReferencingEntityNavigationPropertyName"] == relationLogicalName && item["Entity1LogicalName"] != entityLogicalName {
			tableName = item["ReferencedEntity"].(string)
			break
		}
		if item["ReferencedEntityNavigationPropertyName"] == relationLogicalName && item["Entity2LogicalName"] != entityLogicalName {
			tableName = item["ReferencingEntity"].(string)
			break
		}
	}

	manyToOne, ok := mapResponse["ManyToOneRelationships"].([]interface{})
	if !ok {
		return "", fmt.Errorf("ManyToOneRelationships field is not of type []interface{}")
	}
	for _, list := range manyToOne {
		item := list.(map[string]interface{})
		if item["ReferencingEntityNavigationPropertyName"] == relationLogicalName {
			tableName = item["ReferencedEntity"].(string)
			break
		}
		if item["ReferencedEntityNavigationPropertyName"] == relationLogicalName {
			tableName = item["ReferencingEntity"].(string)
			break
		}
	}

	manyToMany, ok := mapResponse["ManyToManyRelationships"].([]interface{})
	if !ok {
		return "", fmt.Errorf("ManyToManyRelationships field is not of type []interface{}")
	}
	for _, list := range manyToMany {
		item := list.(map[string]interface{})
		if item["Entity1NavigationPropertyName"] == relationLogicalName && item["Entity1LogicalName"] != entityLogicalName {
			tableName = item["Entity1LogicalName"].(string)
			break
		}
		if item["Entity2NavigationPropertyName"] == relationLogicalName && item["Entity2LogicalName"] != entityLogicalName {
			tableName = item["Entity2LogicalName"].(string)
			break
		}
	}
	return tableName, nil
}

func (client *DataRecordClient) ApplyDataRecord(ctx context.Context, recordId, environmentId, tableName string, columns map[string]interface{}) (*DataRecordDto, error) {
	result := DataRecordDto{}

	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	relations := make(map[string]interface{}, 0)

	for key, value := range columns {
		if nestedMap, ok := value.(map[string]interface{}); ok {
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
		} else if nestedMapList, ok := value.([]interface{}); ok {
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

	e, err := url.Parse(environmentUrl)
	if err != nil {
		return nil, err
	}
	apiUrl := &url.URL{
		Scheme: e.Scheme,
		Host:   e.Host,
		Path:   apiPath,
	}

	response, err := client.Api.Execute(ctx, method, apiUrl.String(), nil, columns, []int{http.StatusOK, http.StatusNoContent}, nil)
	if err != nil {
		return nil, err
	}

	if len(response.BodyAsBytes) != 0 {
		json.Unmarshal(response.BodyAsBytes, &result)
	} else if response.Response.Header.Get(constants.HEADER_ODATA_ENTITY_ID) != "" {
		re := regexp.MustCompile("[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}")
		match := re.FindAllStringSubmatch(response.Response.Header.Get(constants.HEADER_ODATA_ENTITY_ID), -1)
		if len(match) > 0 {
			result.Id = match[len(match)-1][0]
		} else {
			return nil, fmt.Errorf("no entity record id returned from the odata-entityid header")
		}
	} else {
		return nil, fmt.Errorf("no entity record id returned from the API")
	}

	err = applyRelations(ctx, client, relations, environmentUrl, result.Id, entityDefinition)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (client *DataRecordClient) DeleteDataRecord(ctx context.Context, recordId string, environmentId string, tableName string, columns map[string]interface{}) error {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return err
	}

	tableEntityDefinition, err := GetEntityDefinition(ctx, client, environmentId, tableName)
	if err != nil {
		return err
	}

	e, err := url.Parse(environmentUrl)
	if err != nil {
		return err
	}

	for key, value := range columns {
		if _, ok := value.(map[string]interface{}); ok {
			delete(columns, key)
		}
		if nestedMapList, ok := value.([]interface{}); ok {
			delete(columns, key)

			for _, nestedItem := range nestedMapList {
				nestedMap := nestedItem.(map[string]interface{})

				dataRecordId, ok := nestedMap["data_record_id"].(string)
				if !ok {
					return fmt.Errorf("data_record_id field is missing or not a string")
				}

				apiUrl := &url.URL{
					Scheme: e.Scheme,
					Host:   e.Host,
					Path:   fmt.Sprintf("/api/data/%s/%s(%s)/%s(%s)/$ref", constants.DATAVERSE_API_VERSION, tableEntityDefinition.LogicalCollectionName, recordId, key, dataRecordId),
				}
				_, err = client.Api.Execute(ctx, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNoContent}, nil)
				if err != nil && !strings.ContainsAny(err.Error(), "404") {
					return err
				}
			}
		}
	}

	apiUrl := &url.URL{
		Scheme: e.Scheme,
		Host:   e.Host,
		Path:   fmt.Sprintf("/api/data/%s/%s(%s)", constants.DATAVERSE_API_VERSION, tableEntityDefinition.LogicalCollectionName, recordId),
	}
	_, err = client.Api.Execute(ctx, "DELETE", apiUrl.String(), nil, columns, []int{http.StatusOK, http.StatusNoContent}, nil)
	if err != nil && !strings.ContainsAny(err.Error(), "404") {
		return err
	}
	return nil
}

func getTableLogicalNameAndDataRecordIdFromMap(nestedMap map[string]interface{}) (string, string, error) {
	tableLogicalName, ok := nestedMap["table_logical_name"].(string)
	if !ok {
		return "", "", fmt.Errorf("table_logical_name field is missing or not a string")
	}
	dataRecordId, ok := nestedMap["data_record_id"].(string)
	if !ok {
		return "", "", fmt.Errorf("data_record_id field is missing or not a string")
	}
	return tableLogicalName, dataRecordId, nil
}

func applyRelations(ctx context.Context, client *DataRecordClient, relations map[string]interface{}, environmentId string, parentRecordId string, entityDefinition *EntityDefinitionsDto) error {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return err
	}

	for key, value := range relations {
		if nestedMapList, ok := value.([]interface{}); ok {
			e, err := url.Parse(environmentUrl)
			if err != nil {
				return err
			}
			apiUrl := &url.URL{
				Scheme: e.Scheme,
				Host:   e.Host,
				Path:   fmt.Sprintf("/api/data/%s/%s(%s)/%s/$ref", constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName, parentRecordId, key),
			}

			existingRelationsResponse := RelationApiResponse{}

			apiResponse, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNoContent}, nil)
			if err != nil {
				return err
			}

			json.Unmarshal(apiResponse.BodyAsBytes, &existingRelationsResponse)

			var toBeDeleted []RelationApiBody = make([]RelationApiBody, 0)

			for _, existingRelation := range existingRelationsResponse.Value {
				delete := true
				for _, nestedItem := range nestedMapList {
					nestedMap := nestedItem.(map[string]interface{})

					tableLogicalName, dataRecordId, err := getTableLogicalNameAndDataRecordIdFromMap(nestedMap)
					if err != nil {
						return err
					}

					relationEntityDefinition, err := GetEntityDefinition(ctx, client, environmentId, tableLogicalName)
					if err != nil {
						return err
					}
					if existingRelation.OdataID == fmt.Sprintf("%s/api/data/%s/%s(%s)", environmentUrl, constants.DATAVERSE_API_VERSION, relationEntityDefinition.LogicalCollectionName, dataRecordId) {
						delete = false
						break
					}
				}
				if delete {
					toBeDeleted = append(toBeDeleted, existingRelation)
				}
			}

			for _, relation := range toBeDeleted {
				_, err = client.Api.Execute(ctx, "DELETE", relation.OdataID, nil, nil, []int{http.StatusOK, http.StatusNoContent}, nil)
				if err != nil {
					return err
				}
			}

			for _, nestedItem := range nestedMapList {
				nestedMap := nestedItem.(map[string]interface{})

				tableLogicalName, dataRecordId, err := getTableLogicalNameAndDataRecordIdFromMap(nestedMap)
				if err != nil {
					return err
				}

				entityDefinition, err := GetEntityDefinition(ctx, client, environmentId, tableLogicalName)
				if err != nil {
					return err
				}

				relation := RelationApiBody{
					OdataID: fmt.Sprintf("%s/api/data/%s/%s(%s)", environmentUrl, constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName, dataRecordId),
				}
				_, err = client.Api.Execute(ctx, "POST", apiUrl.String(), nil, relation, []int{http.StatusOK, http.StatusNoContent}, nil)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}
