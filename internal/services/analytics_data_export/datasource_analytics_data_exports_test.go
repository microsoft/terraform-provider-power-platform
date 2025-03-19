// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package analytics_data_export_test

import (
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccAnalyticsDataExportsDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_analytics_data_exports" "test" {}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.source", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.sink.id", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.sink.type", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.sink.resource_name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.package_name", regexp.MustCompile(helpers.StringRegex)),
					resource.TestMatchResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.resource_provider", regexp.MustCompile(helpers.StringRegex)),
				),
			},
		},
	})
}

func TestUnitAnalyticsDataExportsDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Register mock responses
	httpmock.RegisterResponder(
		"GET",
		"https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/default",
		httpmock.NewStringResponder(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_tenant.json").String()))

	// Register responder for gateway cluster API
	httpmock.RegisterResponder(
		"GET",
		`=~^https://admin\.powerplatform\.microsoft\.com/gateway/cluster.*`,
		httpmock.NewStringResponder(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_gateway_cluster.json").String()))

	// Register responder for analytics data exports API
	httpmock.RegisterResponder(
		"GET",
		"https://na.csanalytics.powerplatform.microsoft.com/api/v2/connections",
		httpmock.NewStringResponder(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_analytics_data_exports.json").String()))

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_analytics_data_exports" "test" {}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.id", "12345678-1234-1234-1234-123456789012"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.source", "Power Platform"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.environments.0.environment_id", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.environments.0.organization_id", "00000000-0000-0000-0000-000000000003"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.status.0.name", "DataExport"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.status.0.state", "Active"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.status.0.last_run_on", "2023-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.status.0.message", "Export is operational"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.sink.id", "/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/rg-example/providers/Microsoft.Insights/components/app-insights-example"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.sink.type", "ApplicationInsights"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.sink.subscription_id", "12345678-1234-1234-1234-123456789012"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.sink.resource_group_name", "rg-example"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.sink.resource_name", "app-insights-example"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.sink.key", "EXAMPLE_KEY"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.package_name", "PowerPlatform.Analytics"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.scenarios.0", "Telemetry"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.scenarios.1", "Usage"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.resource_provider", "Microsoft.PowerPlatform"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.ai_type", "ApplicationInsights"),
				),
			},
		},
	})
}
