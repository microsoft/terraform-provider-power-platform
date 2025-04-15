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
	t.Skip("Skipping test due lack of SP support")

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

	// Register mock response for tenant API
	httpmock.RegisterResponder(
		"GET",
		"https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2021-04-01",
		httpmock.NewStringResponder(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_tenant.json").String()))

	// Register responder for gateway cluster API with dynamic hostname pattern
	httpmock.RegisterResponder(
		"GET",
		`=~^https://.*\.tenant\.api\.powerplatform\.com/gateway/cluster.*`,
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
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.source", "AppInsight"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.environments.0", "00000000-0000-0000-0000-000000000002"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.status.0.name", "Plugin executions excep"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.status.0.state", "Connected"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.status.0.last_run_on", "2025-03-08T06:55:56.0481713+00:00"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.status.1.name", "SDK executions excep"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.status.1.state", "Connected"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.status.1.last_run_on", "2025-03-08T06:55:56.0481713+00:00"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.sink.id", "/subscriptions/00000000-0000-0000-0000-000000000005/resourceGroups/analytics/providers/microsoft.insights/components/insights"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.sink.type", "AppInsights"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.sink.key", "00000000-0000-0000-0000-000000000004"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.sink.resource_name", "insights"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.package_name", "dd"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.scenarios.0", "Plugin executions excep"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.scenarios.1", "SDK executions excep"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.resource_provider", "dataverse"),
					resource.TestCheckResourceAttr("data.powerplatform_analytics_data_exports.test", "exports.0.ai_type", "Local"),
				),
			},
		},
	})
}
