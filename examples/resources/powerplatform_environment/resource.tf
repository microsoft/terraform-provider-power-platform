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

resource "powerplatform_environment" "development" {
  display_name     = "example_environment"
  description      = "example environment description"
  location         = "europe"
  azure_region     = "northeurope"
  environment_type = "Sandbox"
  cadence          = "Moderate"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    domain            = "mydomain"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

