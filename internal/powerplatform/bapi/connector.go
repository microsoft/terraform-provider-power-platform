package powerplatform_bapi

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi/models"
)

func (client *ApiClient) GetConnectors(ctx context.Context) ([]models.ConnectorDto, error) {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   "api.powerapps.com",
		Path:   "/providers/Microsoft.PowerApps/apis",
	}
	values := url.Values{}
	values.Add("api-version", "2023-06-01")
	values.Add("showApisWithToS", "true")
	values.Add("hideDlpExemptApis", "true")
	values.Add("showAllDlpEnforceableApis", "true")
	values.Add("$filter", "environment eq '~Default'")
	apiUrl.RawQuery = values.Encode()

	request, err := http.NewRequestWithContext(ctx, "GET", apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	apiRespose, err := client.doRequest(request)
	if err != nil {
		return nil, err
	}

	connectorArray := models.ConnectorDtoArray{}
	err = apiRespose.MarshallTo(&connectorArray)
	if err != nil {
		return nil, err
	}

	apiUrl = &url.URL{
		Scheme: "https",
		Host:   "api.bap.microsoft.com",
		Path:   "/providers/PowerPlatform.Governance/v1/connectors/metadata/unblockable",
	}
	request, err = http.NewRequestWithContext(ctx, "GET", apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	apiRespose, err = client.doRequest(request)
	if err != nil {
		return nil, err
	}
	unblockableConnectorArray := []models.UnblockableConnectorDto{}
	err = apiRespose.MarshallTo(&unblockableConnectorArray)
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
		Host:   "api.bap.microsoft.com",
		Path:   "/providers/PowerPlatform.Governance/v1/connectors/metadata/virtual",
	}
	request, err = http.NewRequestWithContext(ctx, "GET", apiUrl.String(), nil)
	if err != nil {
		return nil, err
	}
	apiRespose, err = client.doRequest(request)
	if err != nil {
		return nil, err
	}
	virtualConnectorArray := []models.VirtualConnectorDto{}
	err = apiRespose.MarshallTo(&virtualConnectorArray)
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
