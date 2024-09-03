// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerapps

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	environment "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/environment"
)

func NewPowerAppssClient(api *api.ApiClient) PowerAppssClient {
	return PowerAppssClient{
		Api:               api,
		environmentClient: environment.NewEnvironmentClient(api),
	}
}

type PowerAppssClient struct {
	Api               *api.ApiClient
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
			Host:   client.Api.GetConfig().Urls.PowerAppsUrl,
			Path:   fmt.Sprintf("/providers/Microsoft.PowerApps/scopes/admin/environments/%s/apps", env.Name),
		}
		values := url.Values{}
		values.Add("api-version", "2023-06-01")
		apiUrl.RawQuery = values.Encode()

		appsArray := PowerAppDtoArray{}
		_, err := client.Api.Execute(ctx, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &appsArray)
		if err != nil {
			return nil, err
		}
		apps = append(apps, appsArray.Value...)

	}
	return apps, nil
}
