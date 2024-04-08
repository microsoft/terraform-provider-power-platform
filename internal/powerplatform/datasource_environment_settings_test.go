// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccTestEnvironmentSettingsDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "example_environment_settings" {
					display_name      = "TestAccTestEnvironmentSettingsDataSource_Validate_Read"
					location          = "europe"
					language_code     = "1033"
					currency_code     = "USD"
					environment_type  = "Sandbox"
					security_group_id = "00000000-0000-0000-0000-000000000000"
				}
				  
				data "powerplatform_environment_settings" "settings" {
					environment_id = powerplatform_environment.example_environment_settings.id
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_audit_enabled", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_read_audit_enabled", "true"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_user_access_audit_enabled", "true"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "email.email_settings.max_upload_file_size_in_bytes", "5242880"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "audit_and_logs.plugin_trace_log_setting", "Off"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "product.behavior_settings.show_dashboard_cards_in_expanded_state", "true"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "product.features.power_apps_component_framework_for_canvas_apps", "false"),
				),
			},
		},
	})
}

func TestUnitTestEnvironmentSettingsDataSource_Validate_Read(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment_settings/tests/datasource/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/organizations`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment_settings/tests/datasource/organisations.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				data "powerplatform_environment_settings" "settings" {
					environment_id = "00000000-0000-0000-0000-000000000001"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_audit_enabled", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_read_audit_enabled", "true"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_user_access_audit_enabled", "true"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "email.email_settings.max_upload_file_size_in_bytes", "5242880"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "audit_and_logs.plugin_trace_log_setting", "Off"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "product.behavior_settings.show_dashboard_cards_in_expanded_state", "true"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "product.features.power_apps_component_framework_for_canvas_apps", "false"),
				),
			},
		},
	})
}
