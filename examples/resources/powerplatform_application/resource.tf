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

data "powerplatform_applications" "microsoft_flow_extension" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
  name           = "Onboarding essentials"
  publisher_name = "Microsoft"
}

resource "powerplatform_application" "install_microsoft_flow_extension" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
  unique_name = data.powerplatform_applications.microsoft_flow_extension.applications[0].unique_name
}
