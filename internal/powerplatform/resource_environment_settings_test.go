// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
)

func TestAccTestEnvironmentSettingsResource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				resource "powerplatform_environment" "example_environment_settings" {
					display_name      = "TestAccTestEnvironmentSettingsResource_Validate_Read"
					location          = "europe" 
					language_code     = "1033"
					currency_code     = "USD"
					environment_type  = "Sandbox"
					security_group_id = "00000000-0000-0000-0000-000000000000"
				}
				  
				resource "powerplatform_environment_settings" "settings" {
					environment_id                         = powerplatform_environment.example_environment_settings.id
					audit_and_logs = {
						plugin_trace_log_setting = "All"
						audit_settings = {
						  is_audit_enabled             = true
						  is_user_access_audit_enabled = true
						  is_read_audit_enabled        = true
						}
					  }
					  email = {
						email_settings = {
						  max_upload_file_size_in_bytes = 100
						}
					  }
					  product = {
						behavior_settings = {
						  show_dashboard_cards_in_expanded_state = true
						}
						features = {
						  power_apps_component_framework_for_canvas_apps = true
						}
					  }
				  }`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_audit_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_read_audit_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_user_access_audit_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "email.email_settings.max_upload_file_size_in_bytes", "100"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.plugin_trace_log_setting", "All"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.behavior_settings.show_dashboard_cards_in_expanded_state", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.features.power_apps_component_framework_for_canvas_apps", "true"),
				),
			},
		},
	})
}

func TestUnitTestEnvironmentSettingsResource_Validate_Read(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var getOrgInx = 0

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment_settings/tests/resources/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/organizations`,
		func(req *http.Request) (*http.Response, error) {
			getOrgInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment_settings/tests/resources/get_organisations_%d.json", getOrgInx)).String()), nil
		})

	httpmock.RegisterResponder("PATCH", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/organizations%2843f51247-aee6-ee11-9048-000d3a688755%29`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsProviderConfig + `
				  resource "powerplatform_environment_settings" "settings" {
					environment_id                         = "00000000-0000-0000-0000-000000000001"
					audit_and_logs = {
						plugin_trace_log_setting = "All"
						audit_settings = {
						  is_audit_enabled             = true
						  is_user_access_audit_enabled = true
						  is_read_audit_enabled        = true
						}
					  }
					  email = {
						email_settings = {
						  max_upload_file_size_in_bytes = 100
						}
					  }
					  product = {
						behavior_settings = {
						  show_dashboard_cards_in_expanded_state = true
						}
						features = {
						  power_apps_component_framework_for_canvas_apps = false
						}
					  }
				  }`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_audit_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_read_audit_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_user_access_audit_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "email.email_settings.max_upload_file_size_in_bytes", "100"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.plugin_trace_log_setting", "All"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.behavior_settings.show_dashboard_cards_in_expanded_state", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.features.power_apps_component_framework_for_canvas_apps", "false"),
				),
			},
		},
	})
}
