// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connectors

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

func newConnectorsClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}

type client struct {
	Api *api.Client
}

func (client *client) GetConnectors(ctx context.Context) ([]connectorDto, error) {
	apiUrl := &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.PowerAppsUrl,
		Path:   "/providers/Microsoft.PowerApps/apis",
	}
	values := url.Values{}
	values.Add(constants.API_VERSION_PARAM, constants.CONNECTORS_API_VERSION)
	values.Add("showApisWithToS", "true")
	values.Add("hideDlpExemptApis", "true")
	values.Add("showAllDlpEnforceableApis", "true")
	values.Add("$filter", "environment eq '~Default'")

	apiUrl.RawQuery = values.Encode()

	connectorArray := connectorArrayDto{}
	_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &connectorArray)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PowerApps connectors: %w", err)
	}

	apiUrl = &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/PowerPlatform.Governance/v1/connectors/metadata/unblockable",
	}
	unblockableConnectorArray := []unblockableConnectorDto{}
	_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &unblockableConnectorArray)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch unblockable connectors metadata: %w", err)
	}

	for inx, connector := range connectorArray.Value {
		for _, unblockableConnector := range unblockableConnectorArray {
			if connector.Id == unblockableConnector.Id {
				connectorArray.Value[inx].Properties.Unblockable = unblockableConnector.Metadata.Unblockable
			}
		}
	}

	apiUrl = &url.URL{
		Scheme: constants.HTTPS,
		Host:   client.Api.GetConfig().Urls.BapiUrl,
		Path:   "/providers/PowerPlatform.Governance/v1/connectors/metadata/virtual",
	}
	virtualConnectorArray := []virtualConnectorDto{}
	_, err = client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &virtualConnectorArray)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch virtual connectors metadata: %w", err)
	}
	for _, virutualConnector := range virtualConnectorArray {
		connectorArray.Value = append(connectorArray.Value, connectorDto{
			Id:   virutualConnector.Id,
			Name: virutualConnector.Metadata.Name,
			Type: virutualConnector.Metadata.Type,
			Properties: connectorPropertiesDto{
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
