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

data "powerplatform_environments" "all" {}

output "all_environments" {
  value = data.powerplatform_environments.all.environments
}