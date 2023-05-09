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

data "powerplatform_data_loss_prevention_policies" "all" {}
