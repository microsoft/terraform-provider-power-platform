// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
	application "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/application"
	auth "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/authorization"
	connectors "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/connectors"
	dlp_policy "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/dlp_policy"
	environment "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/environment"
	licensing "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/licensing"
	managed_environment "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/managed_environment"
	powerapps "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/powerapps"
	solution "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/solution"
	tenant_settings "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/services/tenant_settings"
	"github.com/stretchr/testify/require"
)

const (
	// TestsProviderConfig is a shared configuration to combine with the actual
	// test configuration so the Power Platform client is properly configured.
	// It is also possible to use the POWER_PLATFORM_ environment variables instead.
	//lintignore:AT004
	TestsProviderConfig = `
provider "powerplatform" {
	use_cli = true
}
`
)

var TestUnitTestProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"powerplatform": providerserver.NewProtocol6WithError(NewPowerPlatformProvider(context.Background(), true)()),
}

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"powerplatform": providerserver.NewProtocol6WithError(NewPowerPlatformProvider(context.Background(), false)()),
}

func TestUnitPowerPlatformProviderHasChildDataSources_Basic(t *testing.T) {
	expectedDataSources := []datasource.DataSource{
		powerapps.NewEnvironmentPowerAppsDataSource(),
		environment.NewEnvironmentsDataSource(),
		application.NewEnvironmentApplicationPackagesDataSource(),
		connectors.NewConnectorsDataSource(),
		solution.NewSolutionsDataSource(),
		dlp_policy.NewDataLossPreventionPolicyDataSource(),
		tenant_settings.NewTenantSettingsDataSource(),
		licensing.NewBillingPoliciesDataSource(),
		licensing.NewBillingPoliciesEnvironmetsDataSource(),
		auth.NewSecurityRolesDataSource(),
	}
	datasources := NewPowerPlatformProvider(context.Background())().(*PowerPlatformProvider).DataSources(context.Background())

	require.Equal(t, len(expectedDataSources), len(datasources), "There are an unexpected number of registered data sources")
	for _, d := range datasources {
		require.Contains(t, expectedDataSources, d(), "An unexpected data source was registered")
	}
}

func TestUnitPowerPlatformProviderHasChildResources_Basic(t *testing.T) {
	expectedResources := []resource.Resource{
		environment.NewEnvironmentResource(),
		application.NewEnvironmentApplicationPackageInstallResource(),
		dlp_policy.NewDataLossPreventionPolicyResource(),
		solution.NewSolutionResource(),
		tenant_settings.NewTenantSettingsResource(),
		managed_environment.NewManagedEnvironmentResource(),
		licensing.NewBillingPolicyResource(),
		licensing.NewBillingPolicyEnvironmentResource(),
		auth.NewUserResource(),
	}
	resources := NewPowerPlatformProvider(context.Background())().(*PowerPlatformProvider).Resources(context.Background())

	require.Equal(t, len(expectedResources), len(resources), "There are an unexpected number of registered resources")
	for _, r := range resources {
		require.Contains(t, expectedResources, r(), "An unexpected resource was registered")
	}

}

func TestUnitPowerPlatformProvider_Validate_Telementry_Optout_Is_False(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24expand=properties%2FbillingPolicy&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/datasource/Validate_Read/get_environments.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https:=//api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("User-Agent") != "terraform-provider-power-platform" {
				t.Errorf("User-Agent='terraform-provider-power-platform' is expected when telemetry_optout is set to True")
			}

			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/datasource/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	r.Test(t, r.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []r.TestStep{
			{
				Config: `provider "powerplatform" {
					use_cli = true
					telemetry_optout = false
				}	
				data "powerplatform_environments" "all" {}`,
			},
		},
	})
}

func TestUnitPowerPlatformProvider_Validate_Telementry_Optout_Is_True(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24expand=properties%2FbillingPolicy&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/datasource/Validate_Read/get_environments.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			if req.Header.Get("User-Agent") != "" {
				t.Errorf("User-Agent not expected when telemetry_optout is set to True")
			}

			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/datasource/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	r.Test(t, r.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []r.TestStep{
			{
				Config: `provider "powerplatform" {
					use_cli = true
					telemetry_optout = true
				}	
				data "powerplatform_environments" "all" {}`,
			},
		},
	})
}

func TestAccPreCheck_Basic(t *testing.T) {
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
