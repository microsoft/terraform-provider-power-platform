// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_group_rule_set_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/microsoft/terraform-provider-power-platform/internal/mocks"
)

func TestAccEnvironmentGroupRuleSetResource_Validate_Create(t *testing.T) {
	t.Setenv("TF_ACC", "1")
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `

				resource "powerplatform_environment_group" "example_group" {
  display_name = "example_environment_group_ruleset"
  description  = "Example environment group"
}

      resource "powerplatform_environment_group_rule_set" "example_group_rule_set" {
  environment_group_id = powerplatform_environment_group.example_group.id
  rules = {
    sharing_controls = {
      share_mode      = "exclude sharing with security groups"
      share_max_limit = 10
    }
    usage_insights = {
       insights_enabled = true
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
				Check: resource.ComposeAggregateTestCheckFunc(),
			},
		},
	})
}
