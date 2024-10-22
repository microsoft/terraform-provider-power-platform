terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  use_cli = true
}

resource "powerplatform_environment_group" "example_group" {
  display_name = "example_environment_group_ruleset"
  description  = "Example environment group"
}

resource "powerplatform_environment_group_rule_set" "example_group_rule_set" {
  environment_group_id = powerplatform_environment_group.example_group.id
  rules = {
    sharing_controls = {
      share_mode      = "exclude sharing with security groups"
      share_max_limit = 42
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
