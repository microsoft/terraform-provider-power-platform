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
  username  = var.username
  password  = var.password
}

data "powerplatform_environments" "all_environments" {}

data "powerplatform_applications" "all_applications" {
  environment_id = data.powerplatform_environments.all_environments.environments[0].id
}

/*
data "powerplatform_applications" "microsoft_flow_extension" {
  name           = "Microsoft Flow Extensions"
  publisher_name = "Microsoft Dynamic 365"
}
*/
