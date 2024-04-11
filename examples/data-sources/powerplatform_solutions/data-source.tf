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
  display_name     = "displayname"
  location         = "europe"
  environment_type = "Sandbox"

}

data "powerplatform_solutions" "all" {
  environment_id = powerplatform_environment.env.id
}
