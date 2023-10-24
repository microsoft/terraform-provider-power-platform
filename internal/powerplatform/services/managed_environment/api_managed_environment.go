package powerplaform

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	environment "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/environment"
)

func NewManagedEnvironmentClient(bapi *api.BapiClientApi, dv *api.DataverseClientApi) ManagedEnvironmentClient {
	return ManagedEnvironmentClient{
		bapiClient:        bapi,
		dataverseClient:   dv,
		environmentClient: environment.NewEnvironmentClient(bapi, dv),
	}
}

type ManagedEnvironmentClient struct {
	bapiClient        *api.BapiClientApi
	dataverseClient   *api.DataverseClientApi
	environmentClient environment.EnvironmentClient
}

func (client *ManagedEnvironmentClient) GetManagedEnvironmentSettings(ctx context.Context, environmentId string) (*environment.GovernanceConfigurationDto, error) {

	managedEnvSettings, err := client.environmentClient.GetEnvironment(ctx, environmentId)
	if err != nil {
		return nil, err
	}
	return &managedEnvSettings.Properties.GovernanceConfiguration, nil
}

func (client *ManagedEnvironmentClient) EnableManagedEnvironment(ctx context.Context, managedEnvSettings environment.GovernanceConfigurationDto, environmentId string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.bapiClient.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/environments/%s/governanceConfiguration", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	_, err := client.bapiClient.Execute(ctx, "POST", apiUrl.String(), managedEnvSettings, []int{http.StatusNoContent, http.StatusAccepted}, nil)
	if err != nil {
		return err
	}

	//todo look at location header and follow "https://switzerland.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/98cd1690-992b-4d7c-b7d6-98c4e1b08fec?api-version=2021-04-01"
	time.Sleep(10 * time.Second)

	return nil
}

func (client *ManagedEnvironmentClient) DisableManagedEnvironment(ctx context.Context, environmentId string) error {
	apiUrl := &url.URL{
		Scheme: "https",
		Host:   client.bapiClient.GetConfig().Urls.BapiUrl,
		Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/environments/%s/governanceConfiguration", environmentId),
	}
	values := url.Values{}
	values.Add("api-version", "2021-04-01")
	apiUrl.RawQuery = values.Encode()

	managedEnv := environment.GovernanceConfigurationDto{
		ProtectionLevel: "Basic",
	}

	_, err := client.bapiClient.Execute(ctx, "POST", apiUrl.String(), managedEnv, []int{http.StatusAccepted}, nil)
	if err != nil {
		return err
	}
	return nil

}
