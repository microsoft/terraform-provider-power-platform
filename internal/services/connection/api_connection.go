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

func newConnectionsClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}

type client struct {
	Api *api.Client
}

func (client *client) BuildHostUri(environmentId string) string {
	envId := strings.ReplaceAll(environmentId, "-", "")
	realm := string(envId[len(envId)-2:])
	envId = envId[:len(envId)-2]

	return fmt.Sprintf("%s.%s.environment.%s", envId, realm, client.Api.GetConfig().Urls.PowerPlatformUrl)
}

func (client *client) CreateConnection(ctx context.Context, environmentId, connectorName string, connectionToCreate createDto) (*connectionDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s", connectorName, strings.ReplaceAll(uuid.New().String(), "-", "")),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	connection := connectionDto{}
	_, err := client.Api.Execute(ctx, "PUT", apiUrl.String(), nil, connectionToCreate, []int{http.StatusCreated}, &connection)
	if err != nil {
		return nil, err
	}

	return &connection, nil
}

func (client *client) UpdateConnection(ctx context.Context, environmentId, connectorName, connectionId, displayName string, connParams, connParamsSet map[string]any) (*connectionDto, error) {
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

	updatedConnection := connectionDto{}
	_, err = client.Api.Execute(ctx, "PUT", apiUrl.String(), nil, conn, []int{http.StatusOK}, &updatedConnection)
	if err != nil {
		return nil, err
	}

	return &updatedConnection, nil
}

func (client *client) GetConnection(ctx context.Context, environmentId, connectorName, connectionId string) (*connectionDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s", connectorName, connectionId),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	connection := connectionDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connection)
	if err != nil {
		if strings.Contains(err.Error(), "ConnectionNotFound") {
			return nil, helpers.WrapIntoProviderError(err, helpers.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("Connection '%s' not found", connectionId))
		}
		return nil, err
	}
	return &connection, nil
}

func (client *client) GetConnections(ctx context.Context, environmentId string) ([]connectionDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.BuildHostUri(environmentId),
		Path:   "/connectivity/connections",
	}

	values := url.Values{}
	values.Add("api-version", "1")
	apiUrl.RawQuery = values.Encode()

	connetionsArray := connectionArrayDto{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connetionsArray)
	if err != nil {
		return nil, err
	}

	return connetionsArray.Value, nil
}

func (client *client) DeleteConnection(ctx context.Context, environmentId, connectorName, connectionId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
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

func (client *client) ShareConnection(ctx context.Context, environmentId, connectorName, connectionId, roleName, entraUserObjectId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s/modifyPermissions", connectorName, connectionId),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	share := shareConnectionRequestDto{
		Put: []shareConnectionRequestPutDto{
			{
				Properties: shareConnectionRequestPutPropertiesDto{
					RoleName:     roleName,
					Capabilities: []any{},
					Principal: shareConnectionRequestPutPropertiesPrincipalDto{
						Id:       entraUserObjectId,
						Type:     "ServicePrincipal",
						TenantId: nil,
					},
					NotifyShareTargetOption: "Notify",
				},
			},
		},
		Delete: []shareConnectionRequestDeleteDto{},
	}

	_, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, share, []int{http.StatusOK}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (client *client) GetConnectionShares(ctx context.Context, environmentId, connectorName, connectionId string) (*shareConnectionResponseArrayDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s/permissions", connectorName, connectionId),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	share := shareConnectionResponseArrayDto{}

	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &share)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(share.Value, func(i, j int) bool {
		return share.Value[i].Properties.Principal["id"].(string) < share.Value[j].Properties.Principal["id"].(string)
	})

	return &share, nil
}

func (client *client) GetConnectionShare(ctx context.Context, environmentId, connectorName, connectionId, principalId string) (*shareConnectionResponseDto, error) {
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

func (client *client) UpdateConnectionShare(ctx context.Context, environmentId, connectorName, connectionId string, share shareConnectionRequestDto) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s/modifyPermissions", connectorName, connectionId),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	_, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, share, []int{http.StatusOK}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (client *client) DeleteConnectionShare(ctx context.Context, environmentId, connectorName, connectionId, shareId string) error {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.BuildHostUri(environmentId),
		Path:   fmt.Sprintf("/connectivity/connectors/%s/connections/%s/modifyPermissions", connectorName, connectionId),
	}
	values := url.Values{}
	values.Add("api-version", "1")
	values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
	apiUrl.RawQuery = values.Encode()

	share := shareConnectionRequestDto{
		Put: []shareConnectionRequestPutDto{},
		Delete: []shareConnectionRequestDeleteDto{
			{
				Id: fmt.Sprintf("/providers/Microsoft.PowerApps/apis/%s/connections/%s/permissions/%s", connectorName, connectionId, shareId),
			},
		},
	}

	_, err := client.Api.Execute(ctx, "POST", apiUrl.String(), nil, share, []int{http.StatusOK}, nil)
	if err != nil {
		return err
	}
	return nil
}
