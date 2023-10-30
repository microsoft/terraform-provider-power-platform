package powerplatform

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	connectors "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/connectors"
	dlp_policy "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/dlp_policy"
	environment "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/environment"
	powerapps "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/powerapps"
	solution "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/solution"
	tenant_settings "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/tenant_settings"
	"github.com/stretchr/testify/require"
)

const (
	// AcceptanceTestsProviderConfig is a shared configuration to combine with the actual
	// test configuration so the Power Platform client is properly configured.
	// It is also possible to use the POWER_PLATFORM_ environment variables instead.
	AcceptanceTestsProviderConfig = `
provider "powerplatform" {
}
`
	UnitTestsProviderConfig = `
provider "powerplatform" {
	tenant_id = "_"
	username = "_"
	password = "_"
	client_id = "_"
	secret = "_"
}
`
)

var (
	TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"powerplatform": providerserver.NewProtocol6WithError(NewPowerPlatformProvider()()),
	}
)

func TestUnitPowerPlatformProvider_HasChildDataSources(t *testing.T) {
	expectedDataSources := []datasource.DataSource{
		powerapps.NewPowerAppsDataSource(),
		environment.NewEnvironmentsDataSource(),
		connectors.NewConnectorsDataSource(),
		solution.NewSolutionsDataSource(),
		dlp_policy.NewDataLossPreventionPolicyDataSource(),
		tenant_settings.NewTenantSettingsDataSource(),
	}
	datasources := NewPowerPlatformProvider()().(*PowerPlatformProvider).DataSources(context.Background())

	require.Equal(t, len(expectedDataSources), len(datasources), "There are an unexpected number of registered data sources")
	for _, d := range datasources {
		require.Contains(t, expectedDataSources, d(), "An unexpected data source was registered")
	}
}

func TestUnitPowerPlatformProvider_HasChildResources(t *testing.T) {
	expectedResources := []resource.Resource{
		environment.NewEnvironmentResource(),
		dlp_policy.NewDataLossPreventionPolicyResource(),
		solution.NewSolutionResource(),
		tenant_settings.NewTenantSettingsResource(),
	}
	resources := NewPowerPlatformProvider()().(*PowerPlatformProvider).Resources(context.Background())

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
