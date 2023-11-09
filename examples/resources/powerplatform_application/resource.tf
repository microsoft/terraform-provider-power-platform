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

data "powerplatform_applications" "all_applications" {}

data "powerplatform_applications" "microsoft_flow_extension" {
  name           = "Microsoft Flow Extensions"
  publisher_name = "Microsoft Dynamic 365"
}

resource "powerplatform_application" "development" {
  id = data.powerplatform_applications.microsoft_flow_extension.id
}
