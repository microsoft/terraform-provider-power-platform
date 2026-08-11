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

data "powerplatform_publishers" "example" {
  environment_id = "e9f7a826-3dc2-e0ca-944e-630e3acead06"
}
