// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_settings_test

import (
	"fmt"
	"net/http"
	"regexp"

	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccTestEnvironmentSettingsDataSource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "example_environment_settings" {
					display_name      = "` + mocks.TestName() + `"
					location          = "unitedstates"
					environment_type  = "Sandbox"
					dataverse = {
						language_code     = "1033"
						currency_code     = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}
				  
				data "powerplatform_environment_settings" "settings" {
					environment_id = powerplatform_environment.example_environment_settings.id
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_audit_enabled", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_read_audit_enabled", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_user_access_audit_enabled", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.log_retention_period_in_days", "-1"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "email.email_settings.max_upload_file_size_in_bytes", "5242880"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "audit_and_logs.plugin_trace_log_setting", "Off"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "product.behavior_settings.show_dashboard_cards_in_expanded_state", "false"),
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
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/organizations`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/organisations.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_environment_settings" "settings" {
					environment_id = "00000000-0000-0000-0000-000000000001"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_audit_enabled", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_read_audit_enabled", "true"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_user_access_audit_enabled", "true"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.log_retention_period_in_days", "-1"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "email.email_settings.max_upload_file_size_in_bytes", "5242880"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "audit_and_logs.plugin_trace_log_setting", "Off"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "product.behavior_settings.show_dashboard_cards_in_expanded_state", "true"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "product.features.power_apps_component_framework_for_canvas_apps", "false"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "product.security.allow_application_user_access", "true"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "product.security.allow_microsoft_trusted_service_tags", "true"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "product.security.enable_ip_based_firewall_rule", "true"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "product.security.enable_ip_based_firewall_rule_in_audit_mode", "true"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "product.security.allowed_service_tags_for_firewall.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "product.security.allowed_service_tags_for_firewall.0", "ApiManagement"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "product.security.allowed_service_tags_for_firewall.1", "AppService"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "product.security.reverse_proxy_ip_addresses.#", "2"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "product.security.reverse_proxy_ip_addresses.0", "10.10.1.1"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "product.security.reverse_proxy_ip_addresses.1", "192.168.1.1"),
					resource.TestCheckResourceAttr("data.powerplatform_environment_settings.settings", "product.security.enable_ip_based_cookie_binding", "true"),
				),
			},
		},
	})
}

func TestUnitTestEnvironmentSettingsDataSource_Validate_No_Dataverse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_No_Dataverse/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/organizations`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_No_Dataverse/organisations.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_No_Dataverse/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/datasource/Validate_No_Dataverse/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_No_Dataverse/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "env" {
					display_name                              = "displayname"
					location                                  = "europe"
					environment_type                          = "Sandbox"
				}

				data "powerplatform_environment_settings" "settings" {
					environment_id = powerplatform_environment.env.id
				}`,
				ExpectError: regexp.MustCompile(`No Dataverse exists in environment`),
				Check:       resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}
