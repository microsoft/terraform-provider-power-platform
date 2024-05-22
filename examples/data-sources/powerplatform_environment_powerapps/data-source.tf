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

data "powerplatform_environment_powerapps" "all" {}
