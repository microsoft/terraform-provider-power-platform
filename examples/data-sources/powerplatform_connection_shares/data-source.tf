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
  environment_id = var.environment_id
  connector_name = "shared_azureopenai"
  connection_id  = var.azure_openai_connection_id
}

