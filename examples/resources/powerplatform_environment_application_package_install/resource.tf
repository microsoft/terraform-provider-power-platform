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

resource "powerplatform_environment" "env" {
  display_name     = "example_environment"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

data "powerplatform_environment_application_packages" "application_to_install" {
  environment_id = powerplatform_environment.env.id
  name           = "Power Platform Pipelines"
  publisher_name = "Microsoft Dynamics 365"
}

data "powerplatform_environment_application_packages" "all_applications" {
  environment_id = powerplatform_environment.env.id
}

locals {
  onboarding_essential_application = toset([for each in data.powerplatform_environment_application_packages.all_applications.applications : each if each.application_name == "Onboarding essentials"])
}

resource "powerplatform_environment_application_package_install" "install_sample_application" {
  environment_id = powerplatform_environment.env.id
  unique_name    = data.powerplatform_environment_application_packages.application_to_install.applications[0].unique_name
}
