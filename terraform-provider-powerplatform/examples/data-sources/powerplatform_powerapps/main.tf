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

#return all apps in tenant
data "powerplatform_powerapps" "all" {
}

output "all_powerapps" {
  value = data.powerplatform_powerapps.all.apps
}

#return all apps in specific environment that have a specific display name
data "powerplatform_powerapps" "environment_app_filter" {
  environment_name = "11111111-2222-3333-4444-555555555555"
}

locals {
  only_specific_name = toset(
    [
      for each in data.powerplatform_powerapps.environment_app_filter.apps : 
          each if each.display_name == "App Display Name"
    ])
}

output "environment_powerapp_filter" {
  value = one(local.only_specific_name).name
}