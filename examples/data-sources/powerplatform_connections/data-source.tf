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

data "powerplatform_environments" "all_environments" {}

data "powerplatform_connections" "all_connections" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
}
