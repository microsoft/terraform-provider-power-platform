// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_set_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/jarcoal/httpmock"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

const mockPolicyId = "00000000-0000-0000-0000-0000000000aa"

// mockRuleBasedPolicy stands in for the rule-based policies API, which now backs
// maker_welcome_content as well as the connector and CSP rules. It echoes back whatever
// rule sets are submitted so the tests exercise the real request/response shapes.
func mockRuleBasedPolicy() {
	policy := map[string]any{}
	created := false

	readBody := func(req *http.Request) map[string]any {
		body := map[string]any{}
		if raw, err := io.ReadAll(req.Body); err == nil {
			_ = json.Unmarshal(raw, &body)
		}
		return body
	}

	store := func(body map[string]any) {
		policy = body
		policy["id"] = mockPolicyId
		created = true
	}

	base := `https://api\.powerplatform\.com/governance/ruleBasedPolicies`

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(base+`/environmentGroups/[^/]+/assignments`),
		func(_ *http.Request) (*http.Response, error) {
			if !created {
				return httpmock.NewStringResponse(http.StatusOK, `{"value":[]}`), nil
			}
			return httpmock.NewStringResponse(http.StatusOK, fmt.Sprintf(`{"value":[{"policyId":%q,"resourceType":"EnvironmentGroup"}]}`, mockPolicyId)), nil
		})

	httpmock.RegisterRegexpResponder("POST", regexp.MustCompile(base+`/[0-9a-f-]+/environmentGroups/[^/]+/assignments`),
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, fmt.Sprintf(`{"policyId":%q}`, mockPolicyId)), nil
		})

	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(base+`/[0-9a-f-]+/environmentGroups/[^/]+/assignments`),
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterRegexpResponder("PATCH", regexp.MustCompile(base+`/[0-9a-f-]+/removeRule`),
		func(req *http.Request) (*http.Response, error) {
			removed := map[string]bool{}
			for _, rs := range asRuleSets(readBody(req)) {
				removed[rs["id"].(string)] = true
			}
			kept := []any{}
			for _, rs := range asRuleSets(policy) {
				if !removed[rs["id"].(string)] {
					kept = append(kept, rs)
				}
			}
			policy["ruleSets"] = kept
			return httpmock.NewStringResponse(http.StatusOK, ""), nil
		})

	httpmock.RegisterRegexpResponder("PUT", regexp.MustCompile(base+`/[0-9a-f-]+\?`),
		func(req *http.Request) (*http.Response, error) {
			store(readBody(req))
			out, _ := json.Marshal(policy)
			return httpmock.NewStringResponse(http.StatusOK, string(out)), nil
		})

	httpmock.RegisterRegexpResponder("DELETE", regexp.MustCompile(base+`/[0-9a-f-]+\?`),
		func(_ *http.Request) (*http.Response, error) {
			policy, created = map[string]any{}, false
			return httpmock.NewStringResponse(http.StatusNoContent, ""), nil
		})

	httpmock.RegisterRegexpResponder("GET", regexp.MustCompile(base+`/[0-9a-f-]+\?`),
		func(_ *http.Request) (*http.Response, error) {
			if !created {
				return httpmock.NewStringResponse(http.StatusNotFound, `{}`), nil
			}
			out, _ := json.Marshal(policy)
			return httpmock.NewStringResponse(http.StatusOK, string(out)), nil
		})

	httpmock.RegisterRegexpResponder("POST", regexp.MustCompile(base+`\?`),
		func(req *http.Request) (*http.Response, error) {
			store(readBody(req))
			out, _ := json.Marshal(policy)
			return httpmock.NewStringResponse(http.StatusCreated, string(out)), nil
		})
}

func asRuleSets(body map[string]any) []map[string]any {
	raw, ok := body["ruleSets"].([]any)
	if !ok {
		return nil
	}
	out := make([]map[string]any, 0, len(raw))
	for _, r := range raw {
		if m, ok := r.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}

func TestAccEnvironmentGroupRuleSetResource_Validate_Create(t *testing.T) {
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
	mockRuleBasedPolicy()

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
	mockRuleBasedPolicy()

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

func TestUnitEnvironmentGroupRuleSetResource_Validate_Import(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()
	mockRuleBasedPolicy()

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
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "environment_group_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
			{
				ResourceName:      "powerplatform_environment_group_rule_set.example_group_rule_set",
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "00000000-0000-0000-0000-000000000000",
			},
		},
	})
}

func TestUnitEnvironmentGroupRuleSetResource_Validate_Import_Empty_Ruleset(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()
	mockRuleBasedPolicy()

	httpmock.RegisterResponder("POST", `https://000000000000000000000000000000.01.tenant.api.powerplatform.com/governance/environmentGroups/00000000-0000-0000-0000-000000000000/ruleSets?api-version=2021-10-01-preview`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusCreated, httpmock.File("tests/Validate_Import_Empty_Ruleset/post_rule_set.json").String()), nil
		})

	httpmock.RegisterResponder("GET", `https://000000000000000000000000000000.01.tenant.api.powerplatform.com/governance/environmentGroups/00000000-0000-0000-0000-000000000000/ruleSets?api-version=2021-10-01-preview`,
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/Validate_Import_Empty_Ruleset/get_rule_set.json").String()), nil
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
						usage_insights = null
						maker_welcome_content = null
						solution_checker_enforcement = null
						backup_retention = null
						ai_generated_descriptions = null
						ai_generative_settings = null
					}
				}`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "environment_group_id", "00000000-0000-0000-0000-000000000000"),
				),
			},
			{
				ResourceName:      "powerplatform_environment_group_rule_set.example_group_rule_set",
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "00000000-0000-0000-0000-000000000000",
			},
		},
	})
}

// TestAccEnvironmentGroupRuleSetResource_Validate_Policy_Rules covers the rules served by the
// rule-based policies API alongside the legacy ones, so both backends are exercised in one resource.
func TestAccEnvironmentGroupRuleSetResource_Validate_Policy_Rules(t *testing.T) {
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
						backup_retention = {
							period_in_days = 21
						}
						advanced_connector_policies_only = {
							enabled = true
						}
						content_security_policy = {
							enabled               = true
							enabled_for_canvas    = true
							enabled_for_code_apps = false

							configuration = {
								img_src     = ["https://example.com"]
								connect_src = ["https://api.example.com"]
								strict_csp  = true
							}
							configuration_for_canvas = {
								connect_src = ["https://canvas-api.example.com"]
							}
						}
						advanced_connector_policies = {
							allowed_connectors = [
								{
									connector_id = "shared_commondataservice"
									actions_mode = "all_allowed"
								},
								{
									connector_id    = "shared_office365"
									actions_mode    = "some_allowed"
									allowed_actions = ["SendEmail", "GetEvents"]
								}
							]
						}
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("powerplatform_environment_group_rule_set.example_group_rule_set", "id"),
					resource.TestCheckResourceAttrSet("powerplatform_environment_group_rule_set.example_group_rule_set", "policy_id"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.backup_retention.period_in_days", "21"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.advanced_connector_policies_only.enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.content_security_policy.enabled", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.content_security_policy.enabled_for_code_apps", "false"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.content_security_policy.configuration.strict_csp", "true"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.content_security_policy.configuration.img_src.0", "https://example.com"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.advanced_connector_policies.allowed_connectors.#", "2"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.advanced_connector_policies.allowed_connectors.1.allowed_actions.0", "SendEmail"),
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
						backup_retention = {
							period_in_days = 28
						}
						advanced_connector_policies_only = {
							enabled = false
						}
						advanced_connector_policies = {
							allowed_connectors = [
								{
									connector_id = "shared_commondataservice"
									actions_mode = "all_allowed"
								}
							]
						}
					}
				}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.backup_retention.period_in_days", "28"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.advanced_connector_policies_only.enabled", "false"),
					resource.TestCheckNoResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.content_security_policy.enabled"),
					resource.TestCheckResourceAttr("powerplatform_environment_group_rule_set.example_group_rule_set", "rules.advanced_connector_policies.allowed_connectors.#", "1"),
				),
			},
		},
	})
}
