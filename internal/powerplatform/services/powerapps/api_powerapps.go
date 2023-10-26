package powerplatform

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	environment "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/environment"
)

func NewPowerAppssClient(bapi *api.BapiClientApi, dv *api.DataverseClientApi) PowerAppssClient {
	return PowerAppssClient{
		bapiClient:        bapi,
		dataverseClient:   dv,
		environmentClient: environment.NewEnvironmentClient(bapi, dv),
	}
}

type PowerAppssClient struct {
	bapiClient        *api.BapiClientApi
	dataverseClient   *api.DataverseClientApi
	environmentClient environment.EnvironmentClient
}

func (client *PowerAppssClient) GetPowerApps(ctx context.Context, environmentId string) ([]PowerAppBapi, error) {
	envs, err := client.environmentClient.GetEnvironments(ctx)
	if err != nil {
		return nil, err
	}
	apps := make([]PowerAppBapi, 0)
	for _, env := range envs {
		apiUrl := &url.URL{
			Scheme: "https",
			Host:   client.bapiClient.GetConfig().Urls.PowerAppsUrl,
			Path:   fmt.Sprintf("/providers/Microsoft.PowerApps/scopes/admin/environments/%s/apps", env.Name),
		}
		values := url.Values{}
		values.Add("api-version", "2023-06-01")
		apiUrl.RawQuery = values.Encode()

		appsArray := PowerAppDtoArray{}
		_, err := client.bapiClient.Execute(ctx, "GET", apiUrl.String(), nil, []int{http.StatusOK}, &appsArray)
		if err != nil {
			return nil, err
		}
		apps = append(apps, appsArray.Value...)

	}
	return apps, nil
}
