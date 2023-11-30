terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  client_id = var.client_id
  secret    = var.secret
  tenant_id = var.tenant_id
}

data "powerplatform_environments" "all_environments" {}

data "powerplatform_applications" "onboarding_essentials_extension" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
  name           = "mscc-supplierportal"
  publisher_name = "Microsoft Dynamics 365"
}

data "powerplatform_applications" "all_applications" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
}

locals {
  onboarding_essential_application = toset([for each in data.powerplatform_applications.all_applications.applications : each if each.application_name == "Onboarding essentials"])
}

resource "powerplatform_application" "install_onboarding_essentials_extension" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
  unique_name = data.powerplatform_applications.onboarding_essentials_extension.applications[0].unique_name
}
