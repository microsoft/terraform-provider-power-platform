// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package provider_test

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	test "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
	"github.com/microsoft/terraform-provider-power-platform/internal/provider"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/admin_management_application"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/analytics_data_export"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/application"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/authorization"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/capacity"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/connection"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/connectors"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/copilot_studio_application_insights"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/currencies"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/data_record"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/dlp_policy"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/enterprise_policy"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment_group_rule_set"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment_groups"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment_settings"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment_templates"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/environment_wave"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/languages"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/licensing"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/locations"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/managed_environment"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/powerapps"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/rest"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/solution"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/solution_checker_rules"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/tenant"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/tenant_isolation_policy"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/tenant_settings"
	"github.com/stretchr/testify/require"
)

func TestUnitPowerPlatformProviderHasChildDataSources_Basic(t *testing.T) {
	expectedDataSources := []datasource.DataSource{
		analytics_data_export.NewAnalyticsExportDataSource(),
		powerapps.NewEnvironmentPowerAppsDataSource(),
		environment.NewEnvironmentsDataSource(),
		environment_templates.NewEnvironmentTemplatesDataSource(),
		application.NewEnvironmentApplicationPackagesDataSource(),
		connectors.NewConnectorsDataSource(),
		solution.NewSolutionsDataSource(),
		dlp_policy.NewDataLossPreventionPolicyDataSource(),
		tenant_settings.NewTenantSettingsDataSource(),
		licensing.NewBillingPoliciesDataSource(),
		licensing.NewBillingPoliciesEnvironmetsDataSource(),
		locations.NewLocationsDataSource(),
		languages.NewLanguagesDataSource(),
		currencies.NewCurrenciesDataSource(),
		authorization.NewSecurityRolesDataSource(),
		environment_settings.NewEnvironmentSettingsDataSource(),
		application.NewTenantApplicationPackagesDataSource(),
		connection.NewConnectionsDataSource(),
		connection.NewConnectionSharesDataSource(),
		data_record.NewDataRecordDataSource(),
		rest.NewDataverseWebApiDatasource(),
		capacity.NewTenantCapcityDataSource(),
		tenant.NewTenantDataSource(),
		solution_checker_rules.NewSolutionCheckerRulesDataSource(),
	}
	datasources := provider.NewPowerPlatformProvider(context.Background())().(*provider.PowerPlatformProvider).DataSources(context.Background())

	require.Equalf(t, len(expectedDataSources), len(datasources), "Expected %d data sources, got %d", len(expectedDataSources), len(datasources))
	for _, d := range datasources {
		require.Containsf(t, expectedDataSources, d(), "Data source %+v was not expected", d())
	}
}

func TestUnitPowerPlatformProviderHasChildResources_Basic(t *testing.T) {
	expectedResources := []resource.Resource{
		environment.NewEnvironmentResource(),
		environment_groups.NewEnvironmentGroupResource(),
		application.NewEnvironmentApplicationPackageInstallResource(),
		dlp_policy.NewDataLossPreventionPolicyResource(),
		solution.NewSolutionResource(),
		tenant_settings.NewTenantSettingsResource(),
		managed_environment.NewManagedEnvironmentResource(),
		licensing.NewBillingPolicyResource(),
		licensing.NewBillingPolicyEnvironmentResource(),
		authorization.NewUserResource(),
		environment_settings.NewEnvironmentSettingsResource(),
		data_record.NewDataRecordResource(),
		rest.NewDataverseWebApiResource(),
		connection.NewConnectionResource(),
		connection.NewConnectionShareResource(),
		admin_management_application.NewAdminManagementApplicationResource(),
		environment_group_rule_set.NewEnvironmentGroupRuleSetResource(),
		enterprise_policy.NewEnterpisePolicyResource(),
		copilot_studio_application_insights.NewCopilotStudioApplicationInsightsResource(),
		tenant_isolation_policy.NewTenantIsolationPolicyResource(),
		environment_wave.NewEnvironmentWaveResource(),
	}
	resources := provider.NewPowerPlatformProvider(context.Background())().(*provider.PowerPlatformProvider).Resources(context.Background())

	require.Equalf(t, len(expectedResources), len(resources), "Expected %d resources, got %d", len(expectedResources), len(resources))
	for _, r := range resources {
		require.Containsf(t, expectedResources, r(), "Resource %+v was not expected", r())
	}
}

func TestUnitPowerPlatformProvider_Validate_Telementry_Optout_Is_False(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/(00000000-0000-0000-0000-000000000001|00000000-0000-0000-0000-000000000002)\?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../services/environment/tests/datasource/Validate_Read/get_environments.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24expand=properties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../services/environment/tests/datasource/Validate_Read/get_environments.json").String()), nil
		})

	test.Test(t, test.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []test.TestStep{
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

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/(00000000-0000-0000-0000-000000000001|00000000-0000-0000-0000-000000000002)\?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01$`),
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../services/environment/tests/datasource/Validate_Read/get_environments.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24expand=properties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../services/environment/tests/datasource/Validate_Read/get_environments.json").String()), nil
		})

	test.Test(t, test.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []test.TestStep{
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
