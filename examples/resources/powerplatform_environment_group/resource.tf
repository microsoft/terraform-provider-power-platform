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
  display_name = "example_environment_group"
  description  = "Example environment group"
}

resource "powerplatform_tenant_settings" "environment_routing" {
  walk_me_opt_out = true
  power_platform = {
    governance = {
      enable_default_environment_routing              = true
      environment_routing_all_makers                  = false
      environment_routing_target_environment_group_id = powerplatform_environment_group.example_group.id
      //environment_routing_target_environment_group_id = "00000000-0000-0000-0000-000000000000"
      environment_routing_target_security_group_id = "00000000-0000-0000-0000-000000000000"
    }
  }
}
