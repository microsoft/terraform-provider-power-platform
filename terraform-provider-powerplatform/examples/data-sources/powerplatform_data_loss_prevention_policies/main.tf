terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/powerplatform"
    }
  }
}

provider "powerplatform" {
  username = "${var.username}"
  password = "${var.password}"
  host = "http://localhost:8080"
}

data "powerplatform_data_loss_prevention_policies" "all" {}

output "all_dlp_policies" {
  value = data.powerplatform_data_loss_prevention_policies.all
}