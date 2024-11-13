// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_set_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccEnvironmentGroupRuleSetResource_Validate_Create(t *testing.T) {
	t.Skip("creating rule sets with SP is NOT yet supported")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group" "example_group" {
					display_name = "` + mocks.TestName() + `"
					description  = "` + mocks.TestName() + `"
				}

				resource "powerplatform_environment_group_rule_set" "example_group_rule_set" {
					environment_group_id = powerplatform_environment_group.example_group.id
					rules = {
						sharing_controls = {
							share_mode      = "exclude sharing with security groups"
							share_max_limit = 42
						}
						usage_insights = {
							insights_enabled = false
						}
						maker_welcome_content = {
							maker_onboarding_url      = "https://contoso.com/onboarding"
							maker_onboarding_markdown = "## Welcome to the environment!\n\n**This is a markdown description.**"
						}
						solution_checker_enforcement = {
							solution_checker_mode = "block"
							send_emails_enabled   = true
						}
						backup_retention = {
							period_in_days = 21
						}
						ai_generated_descriptions = {
							ai_description_enabled = false
						}
						ai_generative_settings = {
							move_data_across_regions_enabled = true
							bing_search_enabled              = false
						}
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.sharing_controls.share_mode", "exclude sharing with security groups"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.sharing_controls.share_max_limit", "42"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.usage_insights.insights_enabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.maker_welcome_content.maker_onboarding_url", "https://contoso.com/onboarding"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.maker_welcome_content.maker_onboarding_markdown", "## Welcome to the environment!\n\n**This is a markdown description.**"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.solution_checker_enforcement.solution_checker_mode", "block"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.solution_checker_enforcement.send_emails_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.backup_retention.period_in_days", "21"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generated_descriptions.ai_description_enabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generative_settings.move_data_across_regions_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generative_settings.bing_search_enabled", "false"),
				),
			},
		},
	})
}

func TestUnitEnvironmentGroupRuleSetResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("POST", `https://000000000000000000000000000000.01.tenant.api.powerplatform.com/governance/environmentGroups/00000000-0000-0000-0000-000000000000/ruleSets?api-version=2021-10-01-preview`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/Validate_Create/post_rule_set.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://000000000000000000000000000000.01.tenant.api.powerplatform.com/governance/environmentGroups/00000000-0000-0000-0000-000000000000/ruleSets?api-version=2021-10-01-preview`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Create/get_rule_set.json").String()), nil
		})

	httpmock.RegisterResponder("DELETE", `https://000000000000000000000000000000.01.tenant.api.powerplatform.com/governance/ruleSets/?api-version=2021-10-01-preview`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group_rule_set" "example_group_rule_set" {
					environment_group_id = "00000000-0000-0000-0000-000000000000"
					rules = {
						sharing_controls = {
							share_mode      = "exclude sharing with security groups"
							share_max_limit = 42
						}
						usage_insights = {
							insights_enabled = false
						}
						maker_welcome_content = {
							maker_onboarding_url      = "https://contoso.com/onboarding"
							maker_onboarding_markdown = "## Welcome to the environment!\n\n**This is a markdown description.**"
						}
						solution_checker_enforcement = {
							solution_checker_mode = "block"
							send_emails_enabled   = true
						}
						backup_retention = {
							period_in_days = 21
						}
						ai_generated_descriptions = {
							ai_description_enabled = false
						}
						ai_generative_settings = {
							move_data_across_regions_enabled = true
							bing_search_enabled              = false
						}
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.sharing_controls.share_mode", "exclude sharing with security groups"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.sharing_controls.share_max_limit", "42"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.usage_insights.insights_enabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.maker_welcome_content.maker_onboarding_url", "https://contoso.com/onboarding"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.maker_welcome_content.maker_onboarding_markdown", "## Welcome to the environment!\n\n**This is a markdown description.**"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.solution_checker_enforcement.solution_checker_mode", "block"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.solution_checker_enforcement.send_emails_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.backup_retention.period_in_days", "21"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generated_descriptions.ai_description_enabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generative_settings.move_data_across_regions_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generative_settings.bing_search_enabled", "false"),
				),
			},
		},
	})
}

func TestAccEnvironmentGroupRuleSetResource_Validate_Update(t *testing.T) {
	t.Skip("creating rule sets with SP is NOT yet supported")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group" "example_group" {
					display_name = "` + mocks.TestName() + `"
					description  = "` + mocks.TestName() + `"
				}

				resource "powerplatform_environment_group_rule_set" "example_group_rule_set" {
					environment_group_id = powerplatform_environment_group.example_group.id
					rules = {
						sharing_controls = {
							share_mode      = "exclude sharing with security groups"
							share_max_limit = 42
						}
						usage_insights = {
							insights_enabled = false
						}
						maker_welcome_content = {
							maker_onboarding_url      = "https://contoso.com/onboarding"
							maker_onboarding_markdown = "## Welcome to the environment!\n\n**This is a markdown description.**"
						}
						solution_checker_enforcement = {
							solution_checker_mode = "block"
							send_emails_enabled   = true
						}
						backup_retention = {
							period_in_days = 21
						}
						ai_generated_descriptions = {
							ai_description_enabled = false
						}
						ai_generative_settings = {
							move_data_across_regions_enabled = true
							bing_search_enabled              = false
						}
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.sharing_controls.share_mode", "exclude sharing with security groups"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.sharing_controls.share_max_limit", "42"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.usage_insights.insights_enabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.maker_welcome_content.maker_onboarding_url", "https://contoso.com/onboarding"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.maker_welcome_content.maker_onboarding_markdown", "## Welcome to the environment!\n\n**This is a markdown description.**"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.solution_checker_enforcement.solution_checker_mode", "block"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.solution_checker_enforcement.send_emails_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.backup_retention.period_in_days", "21"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generated_descriptions.ai_description_enabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generative_settings.move_data_across_regions_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generative_settings.bing_search_enabled", "false"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment_group" "example_group" {
					display_name = "` + mocks.TestName() + `"
					description  = "` + mocks.TestName() + `"
				}

				resource "powerplatform_environment_group_rule_set" "example_group_rule_set" {
					environment_group_id = powerplatform_environment_group.example_group.id
					rules = {
						sharing_controls = {
							share_mode      = "no limit"
						}
						usage_insights = {
							insights_enabled = true
						}
						maker_welcome_content = {
							maker_onboarding_url      = "https://contoso.com/onboarding1"
							maker_onboarding_markdown = "## Welcome to the environment!\n\n**This is a markdown description1.**"
						}
						solution_checker_enforcement = {
							solution_checker_mode = "warn"
							send_emails_enabled   = false
						}
						backup_retention = {
							period_in_days = 28
						}
						ai_generated_descriptions = {
							ai_description_enabled = true
						}
						ai_generative_settings = {
							move_data_across_regions_enabled = false
							bing_search_enabled              = true
						}
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.sharing_controls.share_mode", "no limit"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.sharing_controls.share_max_limit"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.usage_insights.insights_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.maker_welcome_content.maker_onboarding_url", "https://contoso.com/onboarding1"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.maker_welcome_content.maker_onboarding_markdown", "## Welcome to the environment!\n\n**This is a markdown description1.**"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.solution_checker_enforcement.solution_checker_mode", "warn"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.solution_checker_enforcement.send_emails_enabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.backup_retention.period_in_days", "28"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generated_descriptions.ai_description_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generative_settings.move_data_across_regions_enabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generative_settings.bing_search_enabled", "true"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment_group" "example_group" {
					display_name = "` + mocks.TestName() + `"
					description  = "` + mocks.TestName() + `"
				}

				resource "powerplatform_environment_group_rule_set" "example_group_rule_set" {
					environment_group_id = powerplatform_environment_group.example_group.id
					rules = {
						sharing_controls = {
							share_mode      = "no limit"
						}
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.sharing_controls.share_mode", "no limit"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.sharing_controls.share_max_limit"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.usage_insights.insights_enabled"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.maker_welcome_content.maker_onboarding_url"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.maker_welcome_content.maker_onboarding_markdown"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.solution_checker_enforcement.solution_checker_mode"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.solution_checker_enforcement.send_emails_enabled"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.backup_retention.period_in_days"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generated_descriptions.ai_description_enabled"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generative_settings.move_data_across_regions_enabled"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generative_settings.bing_search_enabled"),
				),
			},
		},
	})
}

func TestUnitEnvironmentGroupRuleSetResource_Validate_Update(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()

	post_rule_set_inx := -1
	get_rule_set_inx := -1
	put_rule_set_inx := -1

	httpmock.RegisterResponder("POST", `https://000000000000000000000000000000.01.tenant.api.powerplatform.com/governance/environmentGroups/00000000-0000-0000-0000-000000000000/ruleSets?api-version=2021-10-01-preview`,
		func(_ *http.Request) (*http.Response, error) {
			post_rule_set_inx++
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File(fmt.Sprintf("tests/Validate_Update/post_rule_set_%d.json", post_rule_set_inx)).String()), nil
		})

	httpmock.RegisterResponder("PUT", `https://000000000000000000000000000000.01.tenant.api.powerplatform.com/governance/ruleSets/00000000-0000-0000-0000-000000000001?api-version=2021-10-01-preview`,
		func(_ *http.Request) (*http.Response, error) {
			put_rule_set_inx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/Validate_Update/put_rule_set_%d.json", put_rule_set_inx)).String()), nil
		})

	httpmock.RegisterResponder("GET", `https://000000000000000000000000000000.01.tenant.api.powerplatform.com/governance/environmentGroups/00000000-0000-0000-0000-000000000000/ruleSets?api-version=2021-10-01-preview`,
		func(_ *http.Request) (*http.Response, error) {
			get_rule_set_inx++
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/Validate_Update/get_rule_set_%d.json", get_rule_set_inx)).String()), nil
		})

	httpmock.RegisterResponder("DELETE", `https://000000000000000000000000000000.01.tenant.api.powerplatform.com/governance/ruleSets/?api-version=2021-10-01-preview`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterResponder("DELETE", `https://000000000000000000000000000000.01.tenant.api.powerplatform.com/governance/ruleSets/00000000-0000-0000-0000-000000000001?api-version=2021-10-01-preview`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "powerplatform_environment_group_rule_set" "example_group_rule_set" {
					environment_group_id = "00000000-0000-0000-0000-000000000000"
					rules = {
						sharing_controls = {
							share_mode      = "exclude sharing with security groups"
							share_max_limit = 42
						}
						usage_insights = {
							insights_enabled = false
						}
						maker_welcome_content = {
							maker_onboarding_url      = "https://contoso.com/onboarding"
							maker_onboarding_markdown = "## Welcome to the environment!\n\n**This is a markdown description.**"
						}
						solution_checker_enforcement = {
							solution_checker_mode = "block"
							send_emails_enabled   = true
						}
						backup_retention = {
							period_in_days = 21
						}
						ai_generated_descriptions = {
							ai_description_enabled = false
						}
						ai_generative_settings = {
							move_data_across_regions_enabled = true
							bing_search_enabled              = false
						}
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.sharing_controls.share_mode", "exclude sharing with security groups"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.sharing_controls.share_max_limit", "42"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.usage_insights.insights_enabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.maker_welcome_content.maker_onboarding_url", "https://contoso.com/onboarding"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.maker_welcome_content.maker_onboarding_markdown", "## Welcome to the environment!\n\n**This is a markdown description.**"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.solution_checker_enforcement.solution_checker_mode", "block"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.solution_checker_enforcement.send_emails_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.backup_retention.period_in_days", "21"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generated_descriptions.ai_description_enabled", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generative_settings.move_data_across_regions_enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generative_settings.bing_search_enabled", "false"),
				),
			},
			{
				Config: `
				resource "powerplatform_environment_group_rule_set" "example_group_rule_set" {
					environment_group_id = "00000000-0000-0000-0000-000000000000"
					rules = {
						sharing_controls = {
							share_mode      = "no limit"
						}
						usage_insights = {
							insights_enabled = true
						}
						maker_welcome_content = {
							maker_onboarding_url      = "https://contoso.com/onboarding1"
							maker_onboarding_markdown = "## Welcome to the environment!\n\n**This is a markdown description1.**"
						}
						solution_checker_enforcement = {
							solution_checker_mode = "warn"
							send_emails_enabled   = false
						}
						backup_retention = {
							period_in_days = 28
						}
						ai_generated_descriptions = {
							ai_description_enabled = true
						}
						ai_generative_settings = {
							move_data_across_regions_enabled = false
							bing_search_enabled              = true
						}
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
			{
				Config: `
				resource "powerplatform_environment_group_rule_set" "example_group_rule_set" {
					environment_group_id = "00000000-0000-0000-0000-000000000000"
					rules = {
						sharing_controls = {
							share_mode      = "no limit"
						}
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.sharing_controls.share_mode", "no limit"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.sharing_controls.share_max_limit"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.usage_insights.insights_enabled"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.maker_welcome_content.maker_onboarding_url"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.maker_welcome_content.maker_onboarding_markdown"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.solution_checker_enforcement.solution_checker_mode"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.solution_checker_enforcement.send_emails_enabled"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.backup_retention.period_in_days"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generated_descriptions.ai_description_enabled"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generative_settings.move_data_across_regions_enabled"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.ai_generative_settings.bing_search_enabled"),
				),
			},
		},
	})
}
