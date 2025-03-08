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

func TestUnitTestEnvironmentSettingsResource_Validate_Create_Empty_Settings(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	var getOrgInx = 0

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resources/Validate_Create_Empty_Settings/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/organizations`,
		func(req *http.Request) (*http.Response, error) {
			getOrgInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resources/Validate_Create_Empty_Settings/get_organisations_%d.json", getOrgInx)).String()), nil
		})

	httpmock.RegisterResponder("PATCH", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/organizations%2843f51247-aee6-ee11-9048-000d3a688755%29`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				  resource "powerplatform_environment_settings" "settings" {
					environment_id                         = "00000000-0000-0000-0000-000000000001"
				  }`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_audit_enabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_read_audit_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_user_access_audit_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.log_retention_period_in_days", "-1"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "email.email_settings.max_upload_file_size_in_bytes", "5242880"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.plugin_trace_log_setting", "Off"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.behavior_settings.show_dashboard_cards_in_expanded_state", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.features.power_apps_component_framework_for_canvas_apps", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allow_application_user_access", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allow_microsoft_trusted_service_tags", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.enable_ip_based_firewall_rule", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.enable_ip_based_firewall_rule_in_audit_mode", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allowed_ip_range_for_firewall.#", "0"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allowed_service_tags_for_firewall.#", "0"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.reverse_proxy_ip_addresses.#", "0"),
				),
			},
		},
	})
}

func TestAccTestEnvironmentSettingsResource_Validate_Create_Empty_Settings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		IsUnitTest:               false,
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

					resource "powerplatform_environment_settings" "settings" {
						environment_id                         = powerplatform_environment.example_environment_settings.id
					}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_audit_enabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_read_audit_enabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_user_access_audit_enabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.log_retention_period_in_days", "-1"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "email.email_settings.max_upload_file_size_in_bytes", "5242880"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.plugin_trace_log_setting", "Off"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.behavior_settings.show_dashboard_cards_in_expanded_state", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.features.power_apps_component_framework_for_canvas_apps", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allow_application_user_access", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allow_microsoft_trusted_service_tags", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.enable_ip_based_firewall_rule", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.enable_ip_based_firewall_rule_in_audit_mode", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allowed_ip_range_for_firewall.#", "0"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allowed_service_tags_for_firewall.#", "0"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.reverse_proxy_ip_addresses.#", "0"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.enable_ip_based_cookie_binding", "false"),
				),
			},
		},
	})
}

func TestAccTestEnvironmentSettingsResource_Validate_Read(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		ExternalProviders: map[string]resource.ExternalProvider{
			"time": {
				Source: "hashicorp/time",
			},
		},
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

				resource "powerplatform_managed_environment" "managed_environment" {
					environment_id                  = powerplatform_environment.example_environment_settings.id
					is_usage_insights_disabled      = true
					is_group_sharing_disabled       = true
					limit_sharing_mode              = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing          = 10
					solution_checker_mode           = "Warn"
					suppress_validation_emails      = true
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg"])
					maker_onboarding_markdown       = "this is example markdown"
					maker_onboarding_url            = "https://www.microsoft.com"
				}

				resource "time_sleep" "wait_60_seconds" {
					depends_on = [powerplatform_managed_environment.managed_environment]
					create_duration = "60s"
				}
			
				resource "powerplatform_environment_settings" "settings" {
					depends_on = [time_sleep.wait_60_seconds]

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
						security = {
						  allow_application_user_access               = true
						  allow_microsoft_trusted_service_tags        = true
						  allowed_ip_range_for_firewall               = toset(["10.10.0.0/16", "192.168.0.0/24"])
						  allowed_service_tags_for_firewall           = toset(["ApiManagement", "AppService"])
						  enable_ip_based_firewall_rule               = true
						  enable_ip_based_firewall_rule_in_audit_mode = true
						  reverse_proxy_ip_addresses                  = toset(["10.10.1.1", "192.168.1.1"])
						  enable_ip_based_cookie_binding              = true
						}
					}
				  }`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_audit_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_read_audit_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_user_access_audit_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.log_retention_period_in_days", "-1"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "email.email_settings.max_upload_file_size_in_bytes", "100"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.plugin_trace_log_setting", "All"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.behavior_settings.show_dashboard_cards_in_expanded_state", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.features.power_apps_component_framework_for_canvas_apps", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allow_application_user_access", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allow_microsoft_trusted_service_tags", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.enable_ip_based_firewall_rule", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.enable_ip_based_firewall_rule_in_audit_mode", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allowed_ip_range_for_firewall.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allowed_ip_range_for_firewall.0", "10.10.0.0/16"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allowed_ip_range_for_firewall.1", "192.168.0.0/24"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allowed_service_tags_for_firewall.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allowed_service_tags_for_firewall.0", "ApiManagement"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allowed_service_tags_for_firewall.1", "AppService"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.reverse_proxy_ip_addresses.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.reverse_proxy_ip_addresses.0", "10.10.1.1"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.reverse_proxy_ip_addresses.1", "192.168.1.1"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.enable_ip_based_cookie_binding", "true"),
				),
			},
		},
	})
}

func TestAccTestEnvironmentSettingsResource_Validate_No_Dataverse(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "example_environment_settings" {
					display_name      = "` + mocks.TestName() + `"
					location          = "unitedstates" 
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
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resources/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/organizations`,
		func(req *http.Request) (*http.Response, error) {
			getOrgInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resources/Validate_Read/get_organisations_%d.json", getOrgInx)).String()), nil
		})

	httpmock.RegisterResponder("PATCH", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/organizations%2843f51247-aee6-ee11-9048-000d3a688755%29`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
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
						security = {
						  allow_application_user_access               = true
						  allow_microsoft_trusted_service_tags        = true
						  allowed_ip_range_for_firewall               = toset(["10.10.0.0/16", "192.168.0.0/24"])
						  allowed_service_tags_for_firewall           = toset(["ApiManagement", "AppService"])
						  enable_ip_based_firewall_rule               = true
						  enable_ip_based_firewall_rule_in_audit_mode = true
						  reverse_proxy_ip_addresses                  = toset(["10.10.1.1", "192.168.1.1"])
						  enable_ip_based_cookie_binding              = true
						}
					  }
				  }`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_audit_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_read_audit_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.is_user_access_audit_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.audit_settings.log_retention_period_in_days", "-1"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "email.email_settings.max_upload_file_size_in_bytes", "100"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "audit_and_logs.plugin_trace_log_setting", "All"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.behavior_settings.show_dashboard_cards_in_expanded_state", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.features.power_apps_component_framework_for_canvas_apps", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allow_application_user_access", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allow_microsoft_trusted_service_tags", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.enable_ip_based_firewall_rule", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.enable_ip_based_firewall_rule_in_audit_mode", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allowed_ip_range_for_firewall.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allowed_ip_range_for_firewall.0", "10.10.0.0/16"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allowed_ip_range_for_firewall.1", "192.168.0.0/24"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allowed_service_tags_for_firewall.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allowed_service_tags_for_firewall.0", "ApiManagement"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.allowed_service_tags_for_firewall.1", "AppService"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.reverse_proxy_ip_addresses.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.reverse_proxy_ip_addresses.0", "10.10.1.1"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.reverse_proxy_ip_addresses.1", "192.168.1.1"),
					resource.TestCheckResourceAttr("powerplatform_environment_settings.settings", "product.security.enable_ip_based_cookie_binding", "true"),
				),
			},
		},
	})
}

func TestUnitTestEnvironmentSettingsResource_Validate_No_Dataverse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	var getOrgInx = 0

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resources/Validate_No_Dataverse/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.0/organizations`,
		func(req *http.Request) (*http.Response, error) {
			getOrgInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resources/Validate_No_Dataverse/get_organisations_%d.json", getOrgInx)).String()), nil
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
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resources/Validate_No_Dataverse/get_lifecycle_delete.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resources/Validate_No_Dataverse/get_environment_%s.json", id)).String()), nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resources/Validate_No_Dataverse/get_lifecycle.json").String()), nil
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
