package powerplatform_bapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

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

	return connectorArray.Value, nil
}
