package powerplatform

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	bapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api/bapi"
	dvapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api/dataverse"
	ppapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api/ppapi"
	common "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/common"
	"github.com/stretchr/testify/require"
)

const (
	// providerConfig is a shared configuration to combine with the actual
	// test configuration so the Power Platform client is properly configured.
	// It is also possible to use the POWER_PLATFORM_ environment variables instead.
	providerConfig = `
provider "powerplatform" {
}
`
	uniTestsProviderConfig = `
provider "powerplatform" {
	tenant_id = "_"
	username = "_"
	password = "_"
}
`
)

func powerPlatformProviderServerApiMock(bapiClient bapi.BapiClientInterface, dvClient dvapi.DataverseClientInterface, ppClient ppapi.PowerPlatformClientApiInterface) func() (tfprotov6.ProviderServer, error) {
	providerMock := providerserver.NewProtocol6WithError(&PowerPlatformProvider{
		Config: &common.ProviderConfig{
			Credentials: &common.ProviderCredentials{},
		},
		BapiApi: &BapiClient{
			Client: bapiClient,
		},
		DataverseApi: &DataverseClient{
			Client: dvClient,
		},
		PowerPlatformApi: &PowerPlatoformApiClient{
			Client: ppClient,
		},
	})
	return providerMock
}

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"powerplatform": providerserver.NewProtocol6WithError(NewPowerPlatformProvider()()),
	}
)

func TestUnitPowerPlatformProvider_HasChildDataSources(t *testing.T) {
	expectedDataSources := []datasource.DataSource{
		NewPowerAppsDataSource(),
		NewEnvironmentsDataSource(),
		NewConnectorsDataSource(),
		NewSolutionsDataSource(),
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
		NewDataLossPreventionPolicyResource(),
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
