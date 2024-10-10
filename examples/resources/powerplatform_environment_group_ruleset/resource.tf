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

# resource "powerplatform_environment_group" "example_group" {
#   display_name = "example_environment_group"
#   description  = "Example environment group"
# }

resource "powerplatform_environment_group_rule_set" "example_group_rule_set" {
  environment_group_id = "bd6b30f1-e31e-4cdd-b82b-689a4b674f2f"
  rules = [
    {
      type = "Sharing controls",
      values = {
        share_mode      = "exclude sharing with security groups"
        share_max_limit = 10
      }
    },
    {
      type = "Usage insights"
      values = {
        insights_enabled = true
      }
    }
  ]
}
