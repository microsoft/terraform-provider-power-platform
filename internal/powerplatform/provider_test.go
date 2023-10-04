package powerplatform

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
	clients "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/clients"
	common "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/common"
	dlp_policy "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/dlp_policy"
	"github.com/stretchr/testify/require"
)

const (
	// ProviderConfig is a shared configuration to combine with the actual
	// test configuration so the Power Platform client is properly configured.
	// It is also possible to use the POWER_PLATFORM_ environment variables instead.
	ProviderConfig = `
provider "powerplatform" {
}
`
	UniTestsProviderConfig = `
provider "powerplatform" {
	tenant_id = "_"
	username = "_"
	password = "_"
	client_id = "_"
	secret = "_"
}
`
)

func powerPlatformProviderServerApiMock(bapiClient api.BapiClientInterface, dvClient api.DataverseClientInterface, ppClient api.PowerPlatformClientApiInterface) func() (tfprotov6.ProviderServer, error) {
	providerMock := providerserver.NewProtocol6WithError(&PowerPlatformProvider{
		Config: &common.ProviderConfig{
			Credentials: &common.ProviderCredentials{},
		},
		BapiApi: &clients.BapiClient{
			Client: bapiClient,
		},
		DataverseApi: &clients.DataverseClient{
			Client: dvClient,
		},
		PowerPlatformApi: &clients.PowerPlatoformApiClient{
			Client: ppClient,
		},
	})
	return providerMock
}

var (
	TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"powerplatform": providerserver.NewProtocol6WithError(NewPowerPlatformProvider()()),
	}
)

func TestUnitPowerPlatformProvider_HasChildDataSources(t *testing.T) {
	expectedDataSources := []datasource.DataSource{
		NewPowerAppsDataSource(),
		NewEnvironmentsDataSource(),
		NewConnectorsDataSource(),
		NewSolutionsDataSource(),
		dlp_policy.NewDataLossPreventionPolicyDataSource(),
	}
	datasources := NewPowerPlatformProvider()().(*PowerPlatformProvider).DataSources(nil)

	require.Equal(t, len(expectedDataSources), len(datasources), "There are an unexpected number of registered data sources")
	for _, d := range datasources {
		require.Contains(t, expectedDataSources, d(), "An unexpected data source was registered")
	}
}

func TestUnitPowerPlatformProvider_HasChildResources(t *testing.T) {
	expectedResources := []resource.Resource{
		NewEnvironmentResource(),
		dlp_policy.NewDataLossPreventionPolicyResource(),
		NewSolutionResource(),
	}
	resources := NewPowerPlatformProvider()().(*PowerPlatformProvider).Resources(nil)

	require.Equal(t, len(expectedResources), len(resources), "There are an unexpected number of registered resources")
	for _, r := range resources {
		require.Contains(t, expectedResources, r(), "An unexpected resource was registered")
	}

}

func TestAccPreCheck(t *testing.T) {
	// if v := os.Getenv("POWER_PLATFORM_TENANT_ID"); v == "" {
	// 	t.Fatal("POWER_PLATFORM_TENANT_ID must be set for acceptance tests")
	// }
	// if v := os.Getenv("POWER_PLATFORM_USERNAME"); v == "" {
	// 	t.Fatal("POWER_PLATFORM_USERNAME must be set for acceptance tests")
	// }
	// if v := os.Getenv("POWER_PLATFORM_PASSWORD"); v == "" {
	// 	t.Fatal("POWER_PLATFORM_PASSWORD must be set for acceptance tests")
	// }
}
