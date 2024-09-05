terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
    azuread = {
      source = "hashicorp/azuread"
    }
  }
}

provider "powerplatform" {
  use_cli = true
}

provider "azuread" {
  use_cli = true
}

resource "azuread_group" "environment_routing_target_security_group" {
  display_name     = "example_security_group"
  description      = "Example security group"
  mail_enabled     = false
  security_enabled = true
}

resource "powerplatform_environment_group" "example_group" {
  display_name = "example_environment_group"
  description  = "Example environment group"
}

resource "powerplatform_tenant_settings" "environment_routing" {
  power_platform = {
    governance = {
      enable_default_environment_routing              = false
      environment_routing_all_makers                  = false
      environment_routing_target_environment_group_id = powerplatform_environment_group.example_group.id
      environment_routing_target_security_group_id    = azuread_group.environment_routing_target_security_group.id
    }
  }

  walk_me_opt_out = true

}
