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

data "powerplatform_connections" "all_connections" {
  environment_id = "469aeedd-f3d1-ee95-bafa-3e5364302246"
}
