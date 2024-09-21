// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connection

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

func NewConnectionsClient(apiClient *api.Client) ConnectionsClient {
	return ConnectionsClient{
		Api: apiClient,
	}
}

type ConnectionsClient struct {
	Api *api.Client
}

func (client *ConnectionsClient) BuildHostUri(environmentId string) string {
	envId := strings.ReplaceAll(environmentId, "-", "")
	realm := string(envId[len(envId)-2:])
	envId = envId[:len(envId)-2]

	return fmt.Sprintf("%s.%s.environment.%s", envId, realm, client.Api.GetConfig().Urls.PowerPlatformUrl)
}

func (client *ConnectionsClient) CreateConnection(ctx context.Context, environmentId, connectorName string, connectionToCreate CreateDto) (*Dto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s", connectorName, strings.ReplaceAll(uuid.New().String(), "-", "")),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	connection := Dto{}
	_, err := client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, connectionToCreate, []int{http.StatusCreated}, &connection)
	if err != nil {
		return nil, err
	}

	return &connection, nil
}

func (client *ConnectionsClient) UpdateConnection(ctx context.Context, environmentId, connectorName, connectionId, displayName string, connParams, connParamsSet map[string]any) (*Dto, error) {
	conn, err := client.GetConnection(ctx, environmentId, connectorName, connectionId)
	if err != nil {
		return nil, err
	}

	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s", connectorName, connectionId),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	conn.Properties.DisplayName = displayName
	conn.Properties.ConnectionParametersSet = connParamsSet
	conn.Properties.ConnectionParameters = connParams

	updatedConnection := Dto{}
	_, err = client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, conn, []int{http.StatusOK}, &updatedConnection)
	if err != nil {
		return nil, err
	}

	return &updatedConnection, nil
}

func (client *ConnectionsClient) GetConnection(ctx context.Context, environmentId, connectorName, connectionId string) (*Dto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s", connectorName, connectionId),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	connection := Dto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connection)
	if err != nil {
		if strings.Contains(err.Error(), "ConnectionNotFound") {
			return nil, helpers.WrapIntoProviderError(err, helpers.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("Connection '%s' not found", connectionId))
		}
		return nil, err
	}
	return &connection, nil
}

func (client *ConnectionsClient) GetConnections(ctx context.Context, environmentId string) ([]Dto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.BuildHostUri(environmentId),
		Path:   "/connectivity/connections",
	}

	values := url.Values{}
	values.Add("api-version", "1")
	apiUrl.RawQuery = values.Encode()

	connetionsArray := DtoArray{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connetionsArray)
	if err != nil {
		return nil, err
	}

	return connetionsArray.Value, nil
}

func (client *ConnectionsClient) DeleteConnection(ctx context.Context, environmentId, connectorName, connectionId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s", connectorName, connectionId),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	_, err := client.Api.Execute(ctx, nil, "DELETE", apiUrl.String(), nil, nil, []int{http.StatusOK}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (client *ConnectionsClient) ShareConnection(ctx context.Context, environmentId, connectorName, connectionId, roleName, entraUserObjectId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s/modifyPermissions", connectorName, connectionId),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	share := ShareConnectionRequestDto{
		Put: []ShareConnectionRequestPutDto{
			{
				Properties: ShareConnectionRequestPutPropertiesDto{
					RoleName:     roleName,
					Capabilities: []any{},
					Principal: ShareConnectionRequestPutPropertiesPrincipalDto{
						Id:       entraUserObjectId,
						Type:     "ServicePrincipal",
						TenantId: nil,
					},
					NotifyShareTargetOption: "Notify",
				},
			},
		},
		Delete: []ShareConnectionRequestDeleteDto{},
	}

	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, share, []int{http.StatusOK}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (client *ConnectionsClient) GetConnectionShares(ctx context.Context, environmentId, connectorName, connectionId string) (*ShareConnectionResponseArrayDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s/permissions", connectorName, connectionId),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	share := ShareConnectionResponseArrayDto{}

	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &share)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(share.Value, func(i, j int) bool {
		return share.Value[i].Properties.Principal["id"].(string) < share.Value[j].Properties.Principal["id"].(string)
	})

	return &share, nil
}

func (client *ConnectionsClient) GetConnectionShare(ctx context.Context, environmentId, connectorName, connectionId, principalId string) (*ShareConnectionResponseDto, error) {
	shares, err := client.GetConnectionShares(ctx, environmentId, connectorName, connectionId)
	if err != nil {
		return nil, err
	}

	for _, share := range shares.Value {
		if id, ok := share.Properties.Principal["id"].(string); ok && id == principalId {
			return &share, nil
		}
	}
	return nil, helpers.WrapIntoProviderError(err, helpers.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("Share for principal '%s' not found", principalId))
}

func (client *ConnectionsClient) UpdateConnectionShare(ctx context.Context, environmentId, connectorName, connectionId string, share ShareConnectionRequestDto) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s/modifyPermissions", connectorName, connectionId),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, share, []int{http.StatusOK}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (client *ConnectionsClient) DeleteConnectionShare(ctx context.Context, environmentId, connectorName, connectionId, shareId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s/modifyPermissions", connectorName, connectionId),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	share := ShareConnectionRequestDto{
		Put: []ShareConnectionRequestPutDto{},
		Delete: []ShareConnectionRequestDeleteDto{
			{
				Id: fmt.Sprintf("/providers/Microsoft.PowerApps/apis/%s/connections/%s/permissions/%s", connectorName, connectionId, shareId),
			},
		},
	}

	_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, share, []int{http.StatusOK}, nil)
	if err != nil {
		return err
	}
	return nil
}
