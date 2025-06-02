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

resource "powerplatform_environment_application_admin" "application_user_import" {
  environment_id = "00000000-0000-0000-0000-000000000001"
  application_id = "00000000-0000-0000-0000-000000000002"
}
