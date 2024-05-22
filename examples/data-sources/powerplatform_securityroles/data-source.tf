terraform {
  required_providers {
    power-platform = {
      source = "microsoft/power-platform"
    }
  }
}

provider "power-platform" {
  use_cli = true
}

resource "powerplatform_environment" "environment" {
  display_name     = "datasource_security_roles_example"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

data "powerplatform_security_roles" "all" {
  environment_id = powerplatform_environment.environment.id
}
