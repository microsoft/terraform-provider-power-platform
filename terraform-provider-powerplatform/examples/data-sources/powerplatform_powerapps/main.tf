terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/powerplatform"
    }
  }
}

provider "powerplatform" {
  username = var.username
  password = var.password
  host     = var.host
}

data "powerplatform_powerapps" "all_environments" {}

data "powerplatform_powerapps" "specific_environment" {
  environment_name = "11111111-2222-3333-4444-555555555555"
}