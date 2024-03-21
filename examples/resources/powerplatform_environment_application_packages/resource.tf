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

data "powerplatform_environments" "all_environments" {}

data "powerplatform_environment_application_packages" "application_to_install" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
  name           = "Power Platform Pipelines"
  publisher_name = "Microsoft Dynamics 365"
}

data "powerplatform_environment_application_packages" "all_applications" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
}

locals {
  onboarding_essential_application = toset([for each in data.powerplatform_environment_application_packages.all_applications.applications : each if each.application_name == "Onboarding essentials"])
}

resource "powerplatform_environment_application_package_install" "install_sample_application" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
  unique_name    = data.powerplatform_environment_application_packages.application_to_install.applications[0].unique_name
}
