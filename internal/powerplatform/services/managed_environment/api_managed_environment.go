package powerplatform

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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

	apiResponse, err := client.bapiClient.Execute(ctx, "POST", apiUrl.String(), nil, managedEnvSettings, []int{http.StatusNoContent, http.StatusAccepted}, nil)
	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Managed Environment Enablement Operation HTTP Status: '"+apiResponse.Response.Status+"'")

	tflog.Debug(ctx, "Waiting for Managed Environment Enablement Operation to complete")
	_, err = client.bapiClient.DoWaitForLifecycleOperationStatus(ctx, apiResponse)
	if err != nil {
		return err
	}
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

	apiResponse, err := client.bapiClient.Execute(ctx, "POST", apiUrl.String(), nil, managedEnv, []int{http.StatusAccepted}, nil)
	if err != nil {
		return err
	}

	tflog.Debug(ctx, "Managed Environment Disablement Operation HTTP Status: '"+apiResponse.Response.Status+"'")
	tflog.Debug(ctx, "Waiting for Managed Environment Disablement Operation to complete")

	_, err = client.bapiClient.DoWaitForLifecycleOperationStatus(ctx, apiResponse)
	if err != nil {
		return err
	}
	return nil

}
