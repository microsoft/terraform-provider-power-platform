terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  tenant_id = var.tenant_id
  client_id = var.client_id
  secret    = var.secret
}

data "powerplatform_data_loss_prevention_policies" "all_policies" {}
