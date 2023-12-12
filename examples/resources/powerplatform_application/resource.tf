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

data "powerplatform_applications" "application_to_install" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
  name           = "Power Platform Pipelines"
  publisher_name = "Microsoft Dynamics 365"
}

data "powerplatform_applications" "all_applications" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
}

locals {
  onboarding_essential_application = toset([for each in data.powerplatform_applications.all_applications.applications : each if each.application_name == "Onboarding essentials"])
}

resource "powerplatform_application" "install_sample_application" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
  unique_name = data.powerplatform_applications.application_to_install.applications[0].unique_name
}
