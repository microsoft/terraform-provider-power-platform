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

data "powerplatform_applications" "all_applications" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
}

data "powerplatform_applications" "all_applications_from_publisher" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
  publisher_name = "Power Platform Host Service"
}

data "powerplatform_applications" "specific_application" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
  publisher_name = "Microsoft Dynamics 365"
  name           = "Virtual connectors in Dataverse"
}
