// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

func getEntityDefinition(ctx context.Context, client *DataRecordClient, environmentUrl string, entityLogicalName string) *EntityDefinitionsDto {
	e, _ := url.Parse(environmentUrl)
	entityDefinitionApiUrl := &url.URL{
		Scheme:   e.Scheme,
		Host:     e.Host,
		Path:     fmt.Sprintf("/api/data/%s/EntityDefinitions(LogicalName='%s')", constants.DATAVERSE_API_VERSION, entityLogicalName),
		Fragment: "$select=PrimaryIdAttribute,LogicalCollectionName",
	}
	entityDefinition := EntityDefinitionsDto{}
	_, err := client.Api.Execute(ctx, "GET", entityDefinitionApiUrl.String(), nil, nil, []int{http.StatusOK}, &entityDefinition)
	if err != nil {
		return nil
	}

	return &entityDefinition
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

func (client *DataRecordClient) GetDataRecord(ctx context.Context, recordId string, environmentId string, tableName string) (map[string]interface{}, error) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	entityDefinition := getEntityDefinition(ctx, client, environmentUrl, tableName)

	e, _ := url.Parse(environmentUrl)
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

func (client *DataRecordClient) GetRelationData(ctx context.Context, recordId string, environmentId string, tableName string, relationName string) ([]interface{}, error) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return nil, err
	}

	entityDefinition := getEntityDefinition(ctx, client, environmentUrl, tableName)

	e, _ := url.Parse(environmentUrl)
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
		return nil, fmt.Errorf("value field not found in result")
	}

	value, ok := field.([]interface{})
	if !ok {
		return nil, fmt.Errorf("value field is not of type []interface{}")
	}

	return value, nil
}

func (client *DataRecordClient) GetEntityRelationDefinitionInfo(ctx context.Context, environmentId string, entityLogicalName string, relationLogicalName string) (tableName string) {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return ""
	}

	apiUrl := fmt.Sprintf("%s/api/data/%s/EntityDefinitions(LogicalName='%s')?$expand=OneToManyRelationships,ManyToManyRelationships,ManyToOneRelationships", environmentUrl, constants.DATAVERSE_API_VERSION, entityLogicalName)

	response, err := client.Api.Execute(ctx, "GET", apiUrl, nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return ""
	}

	var mapResponse map[string]interface{}
	json.Unmarshal(response.BodyAsBytes, &mapResponse)

	oneToMany, _ := mapResponse["OneToManyRelationships"].([]interface{})
	for _, list := range oneToMany {
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

	manyToOne, _ := mapResponse["ManyToOneRelationships"].([]interface{})
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

	manyToMany, _ := mapResponse["ManyToManyRelationships"].([]interface{})
	for _, list := range manyToMany {
		item := list.(map[string]interface{})
		if item["Entity1NavigationPropertyName"] == relationLogicalName {
			tableName = item["Entity1LogicalName"].(string)
			break
		}
		if item["Entity2NavigationPropertyName"] == relationLogicalName {
			tableName = item["Entity2LogicalName"].(string)
			break
		}
	}

	return tableName
}

func (client *DataRecordClient) ApplyDataRecord(ctx context.Context, recordId string, environmentId string, tableName string, columns map[string]interface{}) (*DataRecordDto, error) {
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
				tableLogicalName := nestedMap["table_logical_name"].(string)

				entityDefinition := getEntityDefinition(ctx, client, environmentUrl, tableLogicalName)

				columns[fmt.Sprintf("%s@odata.bind", key)] = fmt.Sprintf("%s/api/data/%s/%s(%s)", environmentUrl, constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName, nestedMap["data_record_id"])
			}
		}
		if nestedMapList, ok := value.([]interface{}); ok {
			delete(columns, key)
			relations[key] = nestedMapList
		}
	}

	entityDefinition := getEntityDefinition(ctx, client, environmentUrl, tableName)

	method := "POST"
	apiPath := fmt.Sprintf("/api/data/%s/%s", constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName)

	if val, ok := columns[entityDefinition.PrimaryIDAttribute]; ok {
		method = "PATCH"
		apiPath = fmt.Sprintf("%s(%s)", apiPath, val)
	} else if recordId != "" {
		method = "PATCH"
		apiPath = fmt.Sprintf("%s(%s)", apiPath, recordId)
	}

	e, _ := url.Parse(environmentUrl)
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
	} else if response.Response.Header.Get("OData-EntityId") != "" {
		re := regexp.MustCompile(`\(([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})\)`)
		match := re.FindStringSubmatch(response.Response.Header.Get("OData-EntityId"))
		if len(match) > 1 {
			result.Id = match[1]
		} else {
			return nil, fmt.Errorf("no entity record id returned from the odata-entityid header")
		}
	} else {
		return nil, fmt.Errorf("no entity record id returned from the API")
	}

	locationHeader := response.GetHeader("Location")
	locationHeader = strings.TrimPrefix(locationHeader, fmt.Sprintf("%s/api/data/%s/%s(", environmentUrl, constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName))
	locationHeader = strings.TrimSuffix(locationHeader, ")")

	result.Id = locationHeader

	for key, value := range relations {
		if nestedMapList, ok := value.([]interface{}); ok {
			e, _ := url.Parse(environmentUrl)
			apiUrl := &url.URL{
				Scheme: e.Scheme,
				Host:   e.Host,
				Path:   fmt.Sprintf("/api/data/%s/%s(%s)/%s/$ref", constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName, result.Id, key),
			}

			existingRelationsResponse := RelationApiResponse{}

			apiResponse, _ := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNoContent}, nil)

			json.Unmarshal(apiResponse.BodyAsBytes, &existingRelationsResponse)

			var toBeDeleted []RelationApiBody = make([]RelationApiBody, 0)

			for _, existingRelation := range existingRelationsResponse.Value {
				delete := true
				for _, nestedItem := range nestedMapList {
					nestedMap := nestedItem.(map[string]interface{})
					relationEntityDefinition := getEntityDefinition(ctx, client, environmentUrl, nestedMap["table_logical_name"].(string))
					if existingRelation.OdataID == fmt.Sprintf("%s/api/data/%s/%s(%s)", environmentUrl, constants.DATAVERSE_API_VERSION, relationEntityDefinition.LogicalCollectionName, nestedMap["data_record_id"]) {
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
					return nil, err
				}
			}

			for _, nestedItem := range nestedMapList {
				nestedMap := nestedItem.(map[string]interface{})

				tableLogicalName := nestedMap["table_logical_name"].(string)

				entityDefinition := getEntityDefinition(ctx, client, environmentUrl, tableLogicalName)

				relation := RelationApiBody{
					OdataID: fmt.Sprintf("%s/api/data/%s/%s(%s)", environmentUrl, constants.DATAVERSE_API_VERSION, entityDefinition.LogicalCollectionName, nestedMap["data_record_id"]),
				}
				_, err = client.Api.Execute(ctx, "POST", apiUrl.String(), nil, relation, []int{http.StatusOK, http.StatusNoContent}, nil)
				if err != nil {
					return &result, err
				}
			}
		}
	}

	return &result, nil
}

func (client *DataRecordClient) DeleteDataRecord(ctx context.Context, recordId string, environmentId string, tableName string, columns map[string]interface{}) error {
	environmentUrl, err := client.GetEnvironmentUrlById(ctx, environmentId)
	if err != nil {
		return err
	}

	tableEntityDefinition := getEntityDefinition(ctx, client, environmentUrl, tableName)

	e, _ := url.Parse(environmentUrl)
	apiUrl := &url.URL{
		Scheme: e.Scheme,
		Host:   e.Host,
		Path:   fmt.Sprintf("/api/data/%s/%s(%s)", constants.DATAVERSE_API_VERSION, tableEntityDefinition.LogicalCollectionName, recordId),
	}

	for key, value := range columns {
		if _, ok := value.(map[string]interface{}); ok {
			delete(columns, key)
		}
		if nestedMapList, ok := value.([]interface{}); ok {
			delete(columns, key)

			for _, nestedItem := range nestedMapList {
				nestedMap := nestedItem.(map[string]interface{})

				apiUrl = &url.URL{
					Scheme: e.Scheme,
					Host:   e.Host,
					Path:   fmt.Sprintf("/api/data/%s/%s(%s)/%s(%s)/$ref", constants.DATAVERSE_API_VERSION, tableEntityDefinition.LogicalCollectionName, recordId, key, nestedMap["data_record_id"].(string)),
				}
				_, err = client.Api.Execute(ctx, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK, http.StatusNoContent}, nil)
				if err != nil {
					return err
				}
			}
		}
	}

	_, err = client.Api.Execute(ctx, "DELETE", apiUrl.String(), nil, columns, []int{http.StatusOK, http.StatusNoContent}, nil)
	if err != nil {
		return err
	}

	return nil
}
