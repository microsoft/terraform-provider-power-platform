package powerplatform_bapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi/models"
)

func (client *ApiClient) GetPowerApps(ctx context.Context, environmentId string) ([]models.PowerAppBapi, error) {

	envs, err := client.GetEnvironments(ctx)
	if err != nil {
		return nil, err
	}
	apps := make([]models.PowerAppBapi, 0)
	for _, env := range envs {
		request, err := http.NewRequestWithContext(ctx, "GET", "https://api.powerapps.com/providers/Microsoft.PowerApps/scopes/admin/environments/"+env.Name+"/apps?api-version=2023-06-01", nil)
		if err != nil {
			return nil, err
		}

		body, err := client.doRequest(request)
		if err != nil {
			return nil, err
		}

		appsArray := models.PowerAppDtoArray{}
		err = json.NewDecoder(bytes.NewReader(body)).Decode(&appsArray)
		if err != nil {
			return nil, err
		}
		apps = append(apps, appsArray.Value...)

	}
	return apps, nil
}
