package powerplatform_bapi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi/models"
)

func (client *ApiClient) GetPowerApps(ctx context.Context, environmentId string) ([]models.PowerAppBapi, error) {

	envs, err := client.GetEnvironments(ctx)
	if err != nil {
		return nil, err
	}
	apps := make([]models.PowerAppBapi, 0)
	for _, env := range envs {
		apiUrl := &url.URL{
			Scheme: "https",
			Host:   "api.powerapps.com",
			Path:   fmt.Sprintf("/providers/Microsoft.PowerApps/scopes/admin/environments/%s/apps", env.Name),
		}
		values := url.Values{}
		values.Add("api-version", "2023-06-01")
		apiUrl.RawQuery = values.Encode()
		request, err := http.NewRequestWithContext(ctx, "GET", apiUrl.String(), nil)
		if err != nil {
			return nil, err
		}

		apiResponse, err := client.doRequest(request)
		if err != nil {
			return nil, err
		}

		appsArray := models.PowerAppDtoArray{}
		err = apiResponse.MarshallTo(&appsArray)
		if err != nil {
			return nil, err
		}
		apps = append(apps, appsArray.Value...)

	}
	return apps, nil
}
