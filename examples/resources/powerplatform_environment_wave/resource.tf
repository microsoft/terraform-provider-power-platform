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

resource "powerplatform_environment" "example" {
  display_name     = "wave_feature_example"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "powerplatform_environment_wave" "example" {
  environment_id = powerplatform_environment.example.id
  feature_name   = "April2025Update"

  timeouts {
    create = "45m" # Allow up to 45 minutes for the feature to be installed
  }
}
