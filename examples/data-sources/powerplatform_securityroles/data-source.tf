terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  use_cli = true
}

resource "powerplatform_environment" "environment" {
  display_name      = "datasource_security_roles_example"
  location          = "europe"
  language_code     = "1033"
  currency_code     = "USD"
  environment_type  = "Sandbox"
  security_group_id = "00000000-0000-0000-0000-000000000000"
}

data "powerplatform_securityroles" "all" {
  environment_id = powerplatform_environment.environment.id
}
