package powerplatform_bapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi/models"
)

func (client *ApiClient) GetConnectors(ctx context.Context) ([]models.ConnectorDto, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", "https://api.powerapps.com/providers/Microsoft.PowerApps/apis?api-version=2023-06-01&showApisWithToS=true&hideDlpExemptApis=true&showAllDlpEnforceableApis=true&$filter=environment%20eq%20%27~Default%27", nil)
	if err != nil {
		return nil, err
	}

	body, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}

	connectorArray := models.ConnectorDtoArray{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&connectorArray)
	if err != nil {
		return nil, err
	}

	request, err = http.NewRequestWithContext(ctx, "GET", "https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/connectors/metadata/unblockable", nil)
	if err != nil {
		return nil, err
	}
	body, err = client.doRequest(request)
	if err != nil {
		return nil, err
	}
	unblockableConnectorArray := []models.UnblockableConnectorDto{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&unblockableConnectorArray)
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

	request, err = http.NewRequestWithContext(ctx, "GET", "https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/connectors/metadata/virtual", nil)
	if err != nil {
		return nil, err
	}
	body, err = client.doRequest(request)
	if err != nil {
		return nil, err
	}
	virtualConnectorArray := []models.VirtualConnectorDto{}
	err = json.NewDecoder(bytes.NewReader(body)).Decode(&virtualConnectorArray)
	if err != nil {
		return nil, err
	}
	for _, virutualConnector := range virtualConnectorArray {
		connectorArray.Value = append(connectorArray.Value, models.ConnectorDto{
			Id:   virutualConnector.Id,
			Name: virutualConnector.Metadata.Name,
			Type: virutualConnector.Metadata.Type,
			Properties: models.ConnectorPropertiesDto{
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
