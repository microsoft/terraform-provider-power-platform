// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

func NewConnectorsClient(api *api.ApiClient) ConnectorsClient {
	return ConnectorsClient{
		Api: api,
	}
}

type ConnectorsClient struct {
	Api *api.ApiClient
}

func (client *ConnectorsClient) GetConnectors(ctx context.Context) ([]ConnectorDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.PowerAppsUrl,
		Path:   "/providers/Microsoft.PowerApps/apis",
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	values.Add("showApisWithToS", "true")
	values.Add("hideDlpExemptApis", "true")
	values.Add("showAllDlpEnforceableApis", "true")
	values.Add("$filter", "environment eq '~Default'")
	apiUrl.RawQuery = values.Encode()

	connectorArray := ConnectorDtoArray{}
	_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connectorArray)
	if err != nil {
		return nil, err
	}

	apiUrl = &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/PowerPlatform.Governance/v1/connectors/metadata/unblockable",
	}
	unblockableConnectorArray := []UnblockableConnectorDto{}
	_, err = client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &unblockableConnectorArray)
	if err != nil {
		return nil, err
	}

	for inx, connector := range connectorArray.Value {
		for _, unblockableConnector := range unblockableConnectorArray {
			if connector.Id == unblockableConnector.Id {
				connectorArray.Value[inx].Properties.Unblockable = unblockableConnector.Metadata.Unblockable
			}
		}
	}

	apiUrl = &url.URL{
		Scheme: "https",
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/PowerPlatform.Governance/v1/connectors/metadata/virtual",
	}
	virtualConnectorArray := []VirtualConnectorDto{}
	_, err = client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &virtualConnectorArray)
	if err != nil {
		return nil, err
	}
	for _, virutualConnector := range virtualConnectorArray {
		connectorArray.Value = append(connectorArray.Value, ConnectorDto{
			Id:   virutualConnector.Id,
			Name: virutualConnector.Metadata.Name,
			Type: virutualConnector.Metadata.Type,
			Properties: ConnectorPropertiesDto{
				DisplayName: virutualConnector.Metadata.DisplayName,
				Unblockable: false,
				Tier:        "Built-in",
				Publisher:   "Microsoft",
				Description: "",
			},
		})
	}

	for inx, connector := range connectorArray.Value {
		nameSplit := strings.Split(connector.Id, "/")
		connectorArray.Value[inx].Name = nameSplit[len(nameSplit)-1]
	}

	return connectorArray.Value, nil
}