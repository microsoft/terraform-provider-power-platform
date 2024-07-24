// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	powerplatform_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/helpers"
)

func NewConnectionsClient(api *api.ApiClient) ConnectionsClient {
	return ConnectionsClient{
		Api: api,
	}
}

type ConnectionsClient struct {
	Api *api.ApiClient
}

func (client *ConnectionsClient) BuildHostUri(environmentId string) string {
	envId := strings.ReplaceAll(environmentId, "-", "")
	realm := string(envId[len(envId)-2:])
	envId = envId[:len(envId)-2]

	return fmt.Sprintf("%s.%s.environment.%s", envId, realm, client.Api.GetConfig().Urls.PowerPlatformUrl)

}

func (client *ConnectionsClient) GetConnectorDefinition(ctx context.Context, environmentId, connectorName string) (*ConnectorDefinition, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s", connectorName),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	connector := ConnectorDefinition{}
	_, err := client.Api.Execute(ctx, "PUT", apiUrl.String(), nil, nil, []int{http.StatusCreated}, &connector)
	if err != nil {
		return nil, err
	}

	return &connector, nil
}

func (client *ConnectionsClient) CreateConnection(ctx context.Context, environmentId, connectorName string, connectionToCreate ConnectionToCreateDto) (*ConnectionDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s", connectorName, strings.ReplaceAll(uuid.New().String(), "-", "")),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	connection := ConnectionDto{}
	_, err := client.Api.Execute(ctx, "PUT", apiUrl.String(), nil, connectionToCreate, []int{http.StatusCreated}, &connection)
	if err != nil {
		return nil, err
	}

	return &connection, nil
}

func (client *ConnectionsClient) UpdateConnection(ctx context.Context, environmentId, connectorName, connectionId, displayName string) (*ConnectionDto, error) {

	conn, err := client.GetConnection(ctx, environmentId, connectorName, connectionId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s", connectorName, connectionId),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	conn.Properties.DisplayName = displayName

	updatedConnection := ConnectionDto{}
	_, err = client.Api.Execute(ctx, "PUT", apiUrl.String(), nil, conn, []int{http.StatusOK}, &updatedConnection)
	if err != nil {
		return nil, err
	}

	return &updatedConnection, nil
}

func (client *ConnectionsClient) GetConnection(ctx context.Context, environmentId, connectorName, connectionId string) (*ConnectionDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s", connectorName, connectionId),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	connection := ConnectionDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connection)
	if err != nil {
		if strings.Contains(err.Error(), "ConnectionNotFound") {
			return nil, powerplatform_helpers.WrapIntoProviderError(err, powerplatform_helpers.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("Connection '%s' not found", connectionId))
		}
		return nil, err
	}
	return &connection, nil
}

func (client *ConnectionsClient) GetConnections(ctx context.Context, environmentId string) ([]ConnectionDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BuildHostUri(environmentId),
		Path:   "/connectivity/connections",
	}

	values := url.Values{}
	//values.Add("$expand", "permissions($filter=maxAssignedTo('<<aaid_objectid_here>>'))")
	//values.Add("$filter", fmt.Sprintf("environment eq '%s' and ApiId not in ('shared_logicflows', 'shared_powerflows', 'shared_pqogenericconnector')", environmentId))
	values.Add("api-version", "1")
	apiUrl.RawQuery = values.Encode()

	connetionsArray := ConnectionDtoArray{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connetionsArray)
	if err != nil {
		return nil, err
	}

	return connetionsArray.Value, nil
}

func (client *ConnectionsClient) DeleteConnection(ctx context.Context, environmentId, connectorName, connectionId string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s", connectorName, connectionId),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	_, err := client.Api.Execute(ctx, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (client *ConnectionsClient) ShareConnection(ctx context.Context, environmentId, connectorName, connectionId, roleName, entraUserObjectId string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s/modifyPermissions", connectorName, connectionId),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	// body1 := interface{}(map[string]interface{}{
	// 	"put": []interface{}{
	// 		map[string]interface{}{
	// 			"properties": map[string]interface{}{
	// 				"roleName":     "CanEdit",
	// 				"capabilities": []interface{}{},
	// 				"principal": map[string]interface{}{
	// 					"id":       "f99f844b-ce3b-49ae-86f3-e374ecae789c",
	// 					"type":     "ServicePrincipal",
	// 					"tenantId": nil,
	// 				},
	// 				"NotifyShareTargetOption": "Notify",
	// 			},
	// 		},
	// 	},
	// 	"delete": []interface{}{},
	// })
	share := ShareConnectionRequestDto{
		Put: []ShareConnectionRequestPutDto{
			{
				Properties: ShareConnectionRequestPutPropertiesDto{
					RoleName:     roleName,
					Capabilities: []interface{}{},
					Principal: ShareConnectionRequestPutPropertiesPrincipalDto{
						Id:       entraUserObjectId,
						Type:     "ServicePrincipal",
						TenantId: "null",
					},
				},
			},
		},
		Delete: []ShareConnectionRequestDeleteDto{},
	}

	_, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, share, []int{http.StatusOK}, nil)
	if err != nil {
		//todo: check if permissions does not exists
		return err
	}
	return nil
}

func (client *ConnectionsClient) GetConnectionShares(ctx context.Context, environmentId, connectorName, connectionId string) (*ShareConnectionResponseArrayDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s/permissions", connectorName, connectionId),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	share := ShareConnectionResponseArrayDto{}

	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &share)
	if err != nil {
		return nil, err
	}
	return &share, nil
}
