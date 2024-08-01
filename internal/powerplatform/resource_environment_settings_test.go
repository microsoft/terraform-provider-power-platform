// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"fmt"
	"net/http"
	"regexp"

	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	mock_helpers "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/mocks"
)

func TestAccTestEnvironmentSettingsResource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsAcceptanceProviderConfig + `
				resource "powerplatform_environment" "example_environment_settings" {
					display_name      = "` + mock_helpers.TestName() + `"
					location          = "europe" 
					environment_type  = "Sandbox"
					dataverse = {
						language_code     = "1033"
						currency_code     = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
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

func TestAccTestEnvironmentSettingsResource_Validate_No_Dataverse(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { TestAccPreCheck_Basic(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsAcceptanceProviderConfig + `
				resource "powerplatform_environment" "example_environment_settings" {
					display_name      = "` + mock_helpers.TestName() + `"
					location          = "europe" 
					environment_type  = "Sandbox"
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
				ExpectError: regexp.MustCompile("No Dataverse exists in environment"),
				Check:       resource.ComposeAggregateTestCheckFunc(),
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
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment_settings/tests/resources/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil

		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/organizations`,
		func(req *http.Request) (*http.Response, error) {
			getOrgInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment_settings/tests/resources/Validate_Read/get_organisations_%d.json", getOrgInx)).String()), nil

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
				Config: TestsUnitProviderConfig + `
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

func TestUnitTestEnvironmentSettingsResource_Validate_No_Dataverse(t *testing.T) {

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mock_helpers.ActivateEnvironmentHttpMocks()

	var getOrgInx = 0

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment_settings/tests/resources/Validate_No_Dataverse/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/organizations`,
		func(req *http.Request) (*http.Response, error) {
			getOrgInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment_settings/tests/resources/Validate_No_Dataverse/get_organisations_%d.json", getOrgInx)).String()), nil
		})

	httpmock.RegisterResponder("PATCH", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/organizations%2843f51247-aee6-ee11-9048-000d3a688755%29`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment_settings/tests/resources/Validate_No_Dataverse/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment_settings/tests/resources/Validate_No_Dataverse/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment_settings/tests/resources/Validate_No_Dataverse/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: TestsUnitProviderConfig + `
				  resource "powerplatform_environment" "env" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
				  }

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
				ExpectError: regexp.MustCompile("No Dataverse exists in environment '00000000-0000-0000-0000-000000000001'"),
				Check:       resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}
