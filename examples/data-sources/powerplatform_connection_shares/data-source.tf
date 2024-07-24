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

data "powerplatform_connection_shares" "all_shares" {
  environment_id = "00000000-0000-0000-0000-000000000000"
  connector_name = "shared_azureopenai"
  connection_id  = "11111111-1111-1111-1111-111111111111"
}

