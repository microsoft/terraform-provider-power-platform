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
  display_name = "example_environment_group_rule_based_policy"
  description  = "Example environment group"
}

resource "powerplatform_environment_group_rule_based_policy" "example" {
  environment_group_id = powerplatform_environment_group.example_group.id
  rule_sets = {
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
