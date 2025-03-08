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

data "powerplatform_environment_settings" "example" {
  environment_id = "00000000-0000-0000-0000-000000000001"
}
