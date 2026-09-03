// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managed_environment_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccManagedEnvironmentsResource_Validate_Create(t *testing.T) {
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
				resource "powerplatform_environment" "development" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code    = "1033"
						currency_code    = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_for_dataverse" {
					create_duration = "120s"

					depends_on = [powerplatform_environment.development]
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = true
					is_group_sharing_disabled  = true
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					solution_checker_rule_overrides = toset(["meta-avoid-reg-no-attribute", "meta-avoid-reg-retrieve", "app-use-delayoutput-text-input"])
					power_automate_is_sharing_disabled                 = true
  					copilot_allow_grant_editor_permissions_when_shared = false
  					copilot_limit_sharing_mode                         = "ExcludeSharingToSecurityGroups"
  					copilot_max_limit_user_sharing                     = 55

					depends_on = [time_sleep.wait_for_dataverse]
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),

					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "is_usage_insights_disabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "protection_level", "Standard"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "is_group_sharing_disabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "limit_sharing_mode", "ExcludeSharingToSecurityGroups"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "max_limit_user_sharing", "10"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_mode", "None"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "suppress_validation_emails", "true"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.#", "3"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.1", "meta-avoid-reg-no-attribute"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.2", "meta-avoid-reg-retrieve"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.0", "app-use-delayoutput-text-input"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "power_automate_is_sharing_disabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "copilot_allow_grant_editor_permissions_when_shared", "false"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "copilot_limit_sharing_mode", "ExcludeSharingToSecurityGroups"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "copilot_max_limit_user_sharing", "55"),
				),
			},
		},
	})
}

func TestAccManagedEnvironmentsResource_Validate_Update(t *testing.T) {
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
				resource "powerplatform_environment" "development" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code    = "1033"
						currency_code    = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_for_dataverse" {
					create_duration = "120s"

					depends_on = [powerplatform_environment.development]
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = true
					is_group_sharing_disabled  = true
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])

					depends_on = [time_sleep.wait_for_dataverse]
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code    = "1033"
						currency_code    = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_for_dataverse" {
					create_duration = "120s"

					depends_on = [powerplatform_environment.development]
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = true
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "is_usage_insights_disabled", "false"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code    = "1033"
						currency_code    = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_for_dataverse" {
					create_duration = "120s"

					depends_on = [powerplatform_environment.development]
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "is_group_sharing_disabled", "false"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code    = "1033"
						currency_code    = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_for_dataverse" {
					create_duration = "120s"

					depends_on = [powerplatform_environment.development]
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "limit_sharing_mode", "NoLimit"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code    = "1033"
						currency_code    = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_for_dataverse" {
					create_duration = "120s"

					depends_on = [powerplatform_environment.development]
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "max_limit_user_sharing", "-1"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code    = "1033"
						currency_code    = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_for_dataverse" {
					create_duration = "120s"

					depends_on = [powerplatform_environment.development]
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = false
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "suppress_validation_emails", "false"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code    = "1033"
						currency_code    = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_for_dataverse" {
					create_duration = "120s"

					depends_on = [powerplatform_environment.development]
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = false
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg"])
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.#", "1"),
					resource.TestCheckTypeSetElemAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.*", "meta-remove-dup-reg"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code    = "1033"
						currency_code    = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_for_dataverse" {
					create_duration = "120s"

					depends_on = [powerplatform_environment.development]
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = false
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.#", "2"),
					resource.TestCheckTypeSetElemAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.*", "meta-remove-dup-reg"),
					resource.TestCheckTypeSetElemAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.*", "meta-avoid-reg-no-attribute"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code    = "1033"
						currency_code    = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_for_dataverse" {
					create_duration = "120s"

					depends_on = [powerplatform_environment.development]
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "Warn"
					suppress_validation_emails = false
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_mode", "Warn"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code    = "1033"
						currency_code    = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_for_dataverse" {
					create_duration = "120s"

					depends_on = [powerplatform_environment.development]
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "Warn"
					suppress_validation_emails = false
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					power_automate_is_sharing_disabled = true
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "power_automate_is_sharing_disabled", "true"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code    = "1033"
						currency_code    = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_for_dataverse" {
					create_duration = "120s"

					depends_on = [powerplatform_environment.development]
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "Warn"
					suppress_validation_emails = false
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					power_automate_is_sharing_disabled = true
					copilot_allow_grant_editor_permissions_when_shared = false
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "copilot_allow_grant_editor_permissions_when_shared", "false"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code    = "1033"
						currency_code    = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_for_dataverse" {
					create_duration = "120s"

					depends_on = [powerplatform_environment.development]
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "Warn"
					suppress_validation_emails = false
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					power_automate_is_sharing_disabled = true
					copilot_allow_grant_editor_permissions_when_shared = false
					copilot_limit_sharing_mode = "ExcludeSharingToSecurityGroups"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "copilot_limit_sharing_mode", "ExcludeSharingToSecurityGroups"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
					dataverse = {
						language_code    = "1033"
						currency_code    = "USD"
						security_group_id = "00000000-0000-0000-0000-000000000000"
					}
				}

				resource "time_sleep" "wait_for_dataverse" {
					create_duration = "120s"

					depends_on = [powerplatform_environment.development]
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "Warn"
					suppress_validation_emails = false
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					power_automate_is_sharing_disabled = true
					copilot_allow_grant_editor_permissions_when_shared = false
					copilot_limit_sharing_mode = "ExcludeSharingToSecurityGroups"
					copilot_max_limit_user_sharing = "15"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "copilot_max_limit_user_sharing", "15"),
				),
			},
		},
	})
}

func TestAccManagedEnvironmentsResource_Validate_Wrong_Solution_Checker_Rule_Overrides(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
                resource "powerplatform_environment" "development" {
                    display_name     = "` + mocks.TestName() + `"
                    location         = "unitedstates"
                    environment_type = "Sandbox"
                    dataverse = {
                        language_code    = "1033"
                        currency_code    = "USD"
                        security_group_id = "00000000-0000-0000-0000-000000000000"
                    }
                }
                
                resource "powerplatform_managed_environment" "managed_development" {
                    environment_id             = powerplatform_environment.development.id
                    is_usage_insights_disabled = true
                    is_group_sharing_disabled  = true
                    limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
                    max_limit_user_sharing     = 10
                    solution_checker_mode      = "None"
                    suppress_validation_emails = true
                    solution_checker_rule_overrides = toset(["invalid-rule", "meta-avoid-reg-no-attribute"])
                }`,
				ExpectError: regexp.MustCompile(".*Invalid Solution Checker Rule Override.*"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.#", "2"),
					resource.TestCheckTypeSetElemAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.*", "invalid-rule"),
					resource.TestCheckTypeSetElemAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.*", "meta-avoid-reg-no-attribute"),
				),
			},
		},
	})
}

func TestUnitManagedEnvironmentsResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	patchResponseInx := 0

	httpmock.RegisterResponder("GET", "https://europe.api.advisor.powerapps.com/api/rule?api-version=2.0&ruleset=0ad12346-e108-40b8-a956-9a8f95ea18c9",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/get_rulesset.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment/tests/resource/Validate_Create_And_Update/get_environments_%d.json", patchResponseInx)).String()), nil
		})
	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_lifecycle.json").String()), nil
		})
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_environment_create_response.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = "00000000-0000-0000-0000-000000000001"
					is_usage_insights_disabled = true
					is_group_sharing_disabled  = true
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					power_automate_is_sharing_disabled                 = true
  					copilot_allow_grant_editor_permissions_when_shared = false
  					copilot_limit_sharing_mode                         = "ExcludeSharingToSecurityGroups"
  					copilot_max_limit_user_sharing                     = 55
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),

					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "is_usage_insights_disabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "protection_level", "Standard"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "is_group_sharing_disabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "limit_sharing_mode", "ExcludeSharingToSecurityGroups"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "max_limit_user_sharing", "10"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_mode", "None"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "suppress_validation_emails", "true"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.#", "2"),
					resource.TestCheckTypeSetElemAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.*", "meta-remove-dup-reg"),
					resource.TestCheckTypeSetElemAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.*", "meta-avoid-reg-no-attribute"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "power_automate_is_sharing_disabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "copilot_allow_grant_editor_permissions_when_shared", "false"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "copilot_limit_sharing_mode", "ExcludeSharingToSecurityGroups"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "copilot_max_limit_user_sharing", "55"),
				),
			},
		},
	})
}

// A 409 from the enablement POST means the request was rejected, not applied. It is common straight
// after an environment is created. The provider must retry rather than report success, otherwise the
// environment is left unmanaged and reading state back fails.
func TestUnitManagedEnvironmentsResource_Validate_Create_Retries_When_Enablement_Conflicts(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	enableAttempts := 0
	environmentReads := 0

	httpmock.RegisterResponder("GET", "https://europe.api.advisor.powerapps.com/api/rule?api-version=2.0&ruleset=0ad12346-e108-40b8-a956-9a8f95ea18c9",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/get_rulesset.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Update/get_environments_0.json").String()), nil
		})

	// First enablement attempt is rejected, the retry is accepted.
	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			enableAttempts++
			if enableAttempts == 1 {
				return httpmock.NewStringResponse(http.StatusConflict, ""), nil
			}
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_lifecycle.json").String()), nil
		})

	// The read taken while handling the 409 shows an environment that is not managed yet.
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			environmentReads++
			// The environment only becomes managed once an enablement request is accepted, so until
			// then every read shows an environment without governance settings. That is what makes the
			// provider retry the rejected request rather than treat the 409 as success.
			if enableAttempts < 2 {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_environment_create_response_not_managed.json").String()), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_environment_create_response.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = "00000000-0000-0000-0000-000000000001"
					is_usage_insights_disabled = true
					is_group_sharing_disabled  = true
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					power_automate_is_sharing_disabled                 = true
					copilot_allow_grant_editor_permissions_when_shared = false
					copilot_limit_sharing_mode                         = "ExcludeSharingToSecurityGroups"
					copilot_max_limit_user_sharing                     = 55
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "protection_level", "Standard"),
					func(_ *terraform.State) error {
						if enableAttempts < 2 {
							return fmt.Errorf("expected the rejected enablement to be retried, got %d attempt(s) across %d environment read(s)", enableAttempts, environmentReads)
						}
						return nil
					},
				),
			},
		},
	})
}

// A 409 is also returned when the environment is already managed. That is not a failure, and it must
// not be retried: the environment already has the governance settings we asked for.
func TestUnitManagedEnvironmentsResource_Validate_Create_Conflict_When_Already_Managed(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	enableAttempts := 0

	httpmock.RegisterResponder("GET", "https://europe.api.advisor.powerapps.com/api/rule?api-version=2.0&ruleset=0ad12346-e108-40b8-a956-9a8f95ea18c9",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/get_rulesset.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Update/get_environments_0.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			enableAttempts++
			if enableAttempts == 1 {
				return httpmock.NewStringResponse(http.StatusConflict, ""), nil
			}
			return httpmock.NewStringResponse(http.StatusAccepted, ""), nil
		})

	// The environment is already managed, so the 409 is a no-op rather than a rejection.
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_environment_create_response.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = "00000000-0000-0000-0000-000000000001"
					is_usage_insights_disabled = true
					is_group_sharing_disabled  = true
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					power_automate_is_sharing_disabled                 = true
					copilot_allow_grant_editor_permissions_when_shared = false
					copilot_limit_sharing_mode                         = "ExcludeSharingToSecurityGroups"
					copilot_max_limit_user_sharing                     = 55
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					func(_ *terraform.State) error {
						if enableAttempts != 1 {
							return fmt.Errorf("expected no retry when the environment is already managed, got %d attempt(s)", enableAttempts)
						}
						return nil
					},
				),
			},
		},
	})
}

// An environment in a group is already managed and its governance configuration is locked, so the
// enablement 409 is permanent. The provider must not retry it, and must say what is actually wrong
// rather than claiming the managed environment feature is not enabled.
func TestUnitManagedEnvironmentsResource_Validate_Create_In_Environment_Group(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	enableAttempts := 0

	httpmock.RegisterResponder("GET", "https://europe.api.advisor.powerapps.com/api/rule?api-version=2.0&ruleset=0ad12346-e108-40b8-a956-9a8f95ea18c9",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/get_rulesset.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Update/get_environments_0.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			enableAttempts++
			return httpmock.NewStringResponse(http.StatusConflict, ""), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_environment_in_group.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = "00000000-0000-0000-0000-000000000001"
					is_usage_insights_disabled = true
					is_group_sharing_disabled  = true
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
				}`,
				ExpectError: regexp.MustCompile(`in an environment group`),
			},
		},
	})

	if enableAttempts > 1 {
		t.Errorf("expected no retries for an environment that is already managed, got %d attempts", enableAttempts)
	}
}

func TestUnitManagedEnvironmentsResource_Validate_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	patchResponseInx := -1

	httpmock.RegisterResponder("GET", "https://europe.api.advisor.powerapps.com/api/rule?api-version=2.0&ruleset=0ad12346-e108-40b8-a956-9a8f95ea18c9",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/get_rulesset.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_And_Update/get_environments_0.json").String()), nil
		})
	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			patchResponseInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Create_And_Update/get_environment_create_response_extended_%d.json", patchResponseInx)).String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
						resource "powerplatform_managed_environment" "managed_development" {
						environment_id             = "00000000-0000-0000-0000-000000000001"
						is_usage_insights_disabled = true
						is_group_sharing_disabled  = false
						limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
						max_limit_user_sharing     = 10
						solution_checker_mode      = "None"
						suppress_validation_emails = true
						solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "is_group_sharing_disabled", "false"),
				),
			},
			{
				Config: `
					resource "powerplatform_managed_environment" "managed_development" {
						environment_id             = "00000000-0000-0000-0000-000000000001"
						is_usage_insights_disabled = true
						is_group_sharing_disabled  = false
						limit_sharing_mode         = "NoLimit"
						max_limit_user_sharing     = 10
						solution_checker_mode      = "None"
						suppress_validation_emails = true
						solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "limit_sharing_mode", "NoLimit"),
				),
			},
			{
				Config: `
					resource "powerplatform_managed_environment" "managed_development" {
						environment_id             = "00000000-0000-0000-0000-000000000001"
						is_usage_insights_disabled = true
						is_group_sharing_disabled  = false
						limit_sharing_mode         = "NoLimit"
						max_limit_user_sharing     = -1
						solution_checker_mode      = "None"
						suppress_validation_emails = true
						solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "max_limit_user_sharing", "-1"),
				),
			},
			{
				Config: `
					resource "powerplatform_managed_environment" "managed_development" {
						environment_id             = "00000000-0000-0000-0000-000000000001"
						is_usage_insights_disabled = true
						is_group_sharing_disabled  = false
						limit_sharing_mode         = "NoLimit"
						max_limit_user_sharing     = -1
						solution_checker_mode      = "Warn"
						suppress_validation_emails = true
						solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_mode", "Warn"),
				),
			},
			{
				Config: `
					resource "powerplatform_managed_environment" "managed_development" {
						environment_id             = "00000000-0000-0000-0000-000000000001"
						is_usage_insights_disabled = true
						is_group_sharing_disabled  = false
						limit_sharing_mode         = "NoLimit"
						max_limit_user_sharing     = -1
						solution_checker_mode      = "Warn"
						suppress_validation_emails = false
						solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "suppress_validation_emails", "false"),
				),
			},
			{
				Config: `
					resource "powerplatform_managed_environment" "managed_development" {
						environment_id             = "00000000-0000-0000-0000-000000000001"
						is_usage_insights_disabled = true
						is_group_sharing_disabled  = false
						limit_sharing_mode         = "NoLimit"
						max_limit_user_sharing     = -1
						solution_checker_mode      = "Warn"
						suppress_validation_emails = false
						solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.#", "2"),
					resource.TestCheckTypeSetElemAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.*", "meta-remove-dup-reg"),
					resource.TestCheckTypeSetElemAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.*", "meta-avoid-reg-no-attribute"),
				),
			},
			{
				Config: `
					resource "powerplatform_managed_environment" "managed_development" {
						environment_id             = "00000000-0000-0000-0000-000000000001"
						is_usage_insights_disabled = true
						is_group_sharing_disabled  = false
						limit_sharing_mode         = "NoLimit"
						max_limit_user_sharing     = -1
						solution_checker_mode      = "Warn"
						suppress_validation_emails = false
						solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
						power_automate_is_sharing_disabled                 = true
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "power_automate_is_sharing_disabled", "true"),
				),
			},
		},
	})
}

func TestAccManagedEnvironmentsResource_Validate_No_Dataverse(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment" "development" {
					display_name     = "` + mocks.TestName() + `"
					location         = "unitedstates"
					environment_type = "Sandbox"
				}
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = true
					is_group_sharing_disabled  = true
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
				}`,
				ExpectError: regexp.MustCompile("(?s).*requires Dataverse.*"),
			},
		},
	})
}

func TestUnitManagedEnvironmentsResource_Validate_Import(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	patchResponseInx := 0

	httpmock.RegisterResponder("GET", "https://europe.api.advisor.powerapps.com/api/rule?api-version=2.0&ruleset=0ad12346-e108-40b8-a956-9a8f95ea18c9",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/get_rulesset.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("services/environment/tests/resource/Validate_Create_And_Update/get_environments_%d.json", patchResponseInx)).String()), nil
		})
	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_lifecycle.json").String()), nil
		})
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_environment_create_response.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = "00000000-0000-0000-0000-000000000001"
					is_usage_insights_disabled = true
					is_group_sharing_disabled  = true
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					power_automate_is_sharing_disabled                 = true
					copilot_allow_grant_editor_permissions_when_shared = false
					copilot_limit_sharing_mode                         = "ExcludeSharingToSecurityGroups"
					copilot_max_limit_user_sharing                     = 55
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "environment_id", "00000000-0000-0000-0000-000000000001"),
				),
			},
			{
				ResourceName:      "powerplatform_managed_environment.managed_development",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateId:     "00000000-0000-0000-0000-000000000001",
			},
		},
	})
}

func TestUnitManagedEnvironmentsResource_Validate_Update_Wrong_Solution_Checker_Rule_Overrides(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	patchResponseInx := -1
	httpmock.RegisterResponder("GET", "https://europe.api.advisor.powerapps.com/api/rule?api-version=2.0&ruleset=0ad12346-e108-40b8-a956-9a8f95ea18c9",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/get_rulesset.json").String()), nil
		})
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_Wrong_Solution_Checker_Rule/get_environments_0.json").String()), nil
		})
	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Wrong_Solution_Checker_Rule/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			patchResponseInx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/resource/Validate_Create_Wrong_Solution_Checker_Rule/get_environment_create_response_extended_%d.json", patchResponseInx)).String()), nil
		})

	// Define the test case
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				// Expect an error indicating that the solution checker rule override is invalid `Invalid Solution Checker Rule Override`.
				ExpectError: regexp.MustCompile(`Invalid Solution Checker Rule Override`),
				Config: `
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = "00000000-0000-0000-0000-000000000001"
					is_usage_insights_disabled = true
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "Warn"
					suppress_validation_emails = false
					solution_checker_rule_overrides = toset(["invalid-rule", "meta-avoid-reg-no-attribute"])
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.#", "2"),
					resource.TestCheckTypeSetElemAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.*", "invalid-rule"),
					resource.TestCheckTypeSetElemAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.*", "meta-avoid-reg-no-attribute"),
				),
			},
		},
	})
}

func TestUnitManagedEnvironmentsResource_Validate_Create_No_Dataverse(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://europe.api.advisor.powerapps.com/api/rule?api-version=2.0&ruleset=0ad12346-e108-40b8-a956-9a8f95ea18c9",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/get_rulesset.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Dataverse/get_environment.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile(`requires Dataverse to be`),
				Config: `
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = "00000000-0000-0000-0000-000000000001"
					is_usage_insights_disabled = true
					is_group_sharing_disabled  = true
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
				}`,
			},
		},
	})
}

func TestUnitManagedEnvironmentsResource_Validate_Create_No_Managed_Env(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", "https://europe.api.advisor.powerapps.com/api/rule?api-version=2.0&ruleset=0ad12346-e108-40b8-a956-9a8f95ea18c9",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/get_rulesset.json").String()), nil
		})

	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_And_Update/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_No_Managed_Env/get_environment.json").String()), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				ExpectError: regexp.MustCompile(`doesn't have managed`),
				Config: `
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = "00000000-0000-0000-0000-000000000001"
					is_usage_insights_disabled = true
					is_group_sharing_disabled  = true
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
				}`,
			},
		},
	})
}
