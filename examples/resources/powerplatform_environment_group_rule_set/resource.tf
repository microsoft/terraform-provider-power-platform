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

    advanced_connector_policies_only = {
      enabled = true
    }

    content_security_policy = {
      enabled               = true
      enabled_for_canvas    = true
      enabled_for_code_apps = true

      configuration = {
        img_src     = ["https://example.com"]
        connect_src = ["https://api.example.com"]
        strict_csp  = true
      }

      configuration_for_canvas = {
        connect_src = ["https://canvas-api.example.com"]
        strict_csp  = true
      }

      configuration_for_code_apps = {
        script_src = ["https://cdn.example.com"]
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
