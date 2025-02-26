// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package managed_environment_test

import (
	"fmt"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccManagedEnvironmentsResource_Validate_Create(t *testing.T) {
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
					solution_checker_rule_overrides = toset(["meta-avoid-reg-no-attribute", "meta-avoid-reg-retrieve", "app-use-delayoutput-text-input"])
					maker_onboarding_markdown  = "this is test markdown"
					maker_onboarding_url       = "http://www.example.com"
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
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "maker_onboarding_markdown", "this is test markdown"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "maker_onboarding_url", "http://www.example.com"),
				),
			},
		},
	})
}

func TestAccManagedEnvironmentsResource_Validate_Update(t *testing.T) {
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
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					maker_onboarding_markdown  = "this is test markdown"
					maker_onboarding_url       = "http://www.example.com"
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
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = true
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					maker_onboarding_markdown  = "this is test markdown"
					maker_onboarding_url       = "http://www.example.com"
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
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					maker_onboarding_markdown  = "this is test markdown"
					maker_onboarding_url       = "http://www.example.com"
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
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = 10
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					maker_onboarding_markdown  = "this is test markdown"
					maker_onboarding_url       = "http://www.example.com"
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
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = true
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					maker_onboarding_markdown  = "this is test markdown"
					maker_onboarding_url       = "http://www.example.com"
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
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = false
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					maker_onboarding_markdown  = "this is test markdown"
					maker_onboarding_url       = "http://www.example.com"
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
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = false
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					maker_onboarding_markdown  = "this is test markdown changed"
					maker_onboarding_url       = "http://www.example.com"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides", "meta-remove-dup-reg"),
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
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = false
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					maker_onboarding_markdown  = "this is test markdown changed"
					maker_onboarding_url       = "http://www.example.com"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides", "meta-avoid-reg-no-attribute"),
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
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = false
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					maker_onboarding_markdown  = "this is test markdown changed"
					maker_onboarding_url       = "http://www.example.com"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides", "meta-avoid-reg-no-attribute"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides", "meta-remove-dup-reg"),
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
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = false
					solution_checker_rule_overrides = ""
					maker_onboarding_markdown  = "this is test markdown changed"
					maker_onboarding_url       = "http://www.example.com"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides", ""),
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
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = false
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					maker_onboarding_markdown  = "this is test markdown changed"
					maker_onboarding_url       = "http://www.example.com"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "maker_onboarding_markdown", "this is test markdown changed"),
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
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "None"
					suppress_validation_emails = false
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					maker_onboarding_markdown  = "this is test markdown changed"
					maker_onboarding_url       = "http://www.example-changed.com"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "maker_onboarding_url", "http://www.example-changed.com"),
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
				
				resource "powerplatform_managed_environment" "managed_development" {
					environment_id             = powerplatform_environment.development.id
					is_usage_insights_disabled = false
					is_group_sharing_disabled  = false
					limit_sharing_mode         = "NoLimit"
					max_limit_user_sharing     = -1
					solution_checker_mode      = "Warn"
					suppress_validation_emails = false
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])
					maker_onboarding_markdown  = "this is test markdown changed"
					maker_onboarding_url       = "http://www.example-changed.com"
				}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestMatchResourceAttr("powerplatform_environment.development", "id", regexp.MustCompile(helpers.GuidRegex)),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_mode", "Warn"),
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
                    maker_onboarding_markdown  = "this is test markdown"
                    maker_onboarding_url       = "http://www.example.com"
                }`,
				ExpectError: regexp.MustCompile(`Invalid Attribute Value Match`),
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
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
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
					maker_onboarding_markdown  = "this is test markdown"
					maker_onboarding_url       = "http://www.example.com"
					solution_checker_rule_overrides = toset(["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"])

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
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "maker_onboarding_markdown", "this is test markdown"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "maker_onboarding_url", "http://www.example.com"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.#", "2"),
					resource.TestCheckTypeSetElemAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.*", "meta-remove-dup-reg"),
					resource.TestCheckTypeSetElemAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.*", "meta-avoid-reg-no-attribute"),
				),
			},
		},
	})
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

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
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
						maker_onboarding_markdown  = "this is test markdown"
						maker_onboarding_url       = "http://www.example.com"
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
						maker_onboarding_markdown  = "this is test markdown"
						maker_onboarding_url       = "http://www.example.com"
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
						maker_onboarding_markdown  = "this is test markdown"
						maker_onboarding_url       = "http://www.example.com"
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
						maker_onboarding_markdown  = "this is test markdown"
						maker_onboarding_url       = "http://www.example.com"
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
						maker_onboarding_markdown  = "this is test markdown"
						maker_onboarding_url       = "http://www.example.com"
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
						maker_onboarding_markdown  = "this is test markdown 2"
						maker_onboarding_url       = "http://www.example.com"
					}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "maker_onboarding_markdown", "this is test markdown 2"),
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
						maker_onboarding_markdown  = "this is test markdown 2"
						maker_onboarding_url       = "http://www.example2.com"
					}`,

				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "maker_onboarding_url", "http://www.example2.com"),
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
						maker_onboarding_markdown  = "this is test markdown 2"
						maker_onboarding_url       = "http://www.example2.com"
					}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "id", "00000000-0000-0000-0000-000000000001"),
					resource.TestCheckResourceAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.#", "2"),
					resource.TestCheckTypeSetElemAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.*", "meta-remove-dup-reg"),
					resource.TestCheckTypeSetElemAttr("powerplatform_managed_environment.managed_development", "solution_checker_rule_overrides.*", "meta-avoid-reg-no-attribute"),
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
					maker_onboarding_markdown  = "this is test markdown"
					maker_onboarding_url       = "http://www.example.com"
				}`,
				ExpectError: regexp.MustCompile("InvalidLifecycleOperationRequest"),
				Check:       resource.ComposeAggregateTestCheckFunc(),
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
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("services/environment/tests/resource/Validate_Create_Wrong_Solution_Checker_Rule/get_environments_0.json").String()), nil
		})
	httpmock.RegisterResponder("POST", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/environments/00000000-0000-0000-0000-000000000001/governanceConfiguration?api-version=2021-04-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		})

	httpmock.RegisterResponder("GET", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Create_Wrong_Solution_Checker_Rule/get_lifecycle.json").String()), nil
		})

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy&api-version=2023-06-01",
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
					maker_onboarding_markdown  = "this is test markdown 2"
					maker_onboarding_url       = "http://www.example2.com"
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
