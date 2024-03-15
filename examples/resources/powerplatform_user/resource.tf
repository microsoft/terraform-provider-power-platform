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

resource "powerplatform_environment" "dataverse_user_example" {
  display_name      = "user_example"
  location          = "europe"
  language_code     = "1033"
  currency_code     = "USD"
  environment_type  = "Sandbox"
  security_group_id = "00000000-0000-0000-0000-000000000000"
}

resource "powerplatform_user" "new_user" {
  environment_id = powerplatform_environment.dataverse_user_example.id
  security_roles = [
    "e0d2794e-82f3-e811-a951-000d3a1bcf17", // bot author
  ]
  aad_id         = "00000000-0000-0000-0000-000000000001"
  disable_delete = false
}
