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

resource "powerplatform_environment" "env" {
  display_name     = "Example Environment"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "powerplatform_application_user" "application_user" {
  environment_id = powerplatform_environment.env.id
  application_id = var.application_id
}

resource "powerplatform_role_assignment" "example" {
  environment_id     = powerplatform_environment.env.id
  principal_id       = powerplatform_application_user.application_user.id
  security_role_name = "Basic User"
}
