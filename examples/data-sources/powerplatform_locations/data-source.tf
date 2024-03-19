terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  client_id     = var.client_id
  client_secret = var.secret
  tenant_id     = var.tenant_id
}

data "powerplatform_locations" "all_locations" {
}