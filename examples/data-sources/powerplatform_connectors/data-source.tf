terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "github.com/microsoft/terraform-provider-power-platform"
    }
  }
}

provider "powerplatform" {
  username = var.username
  password = var.password
  tenant_id = var.tenant_id
}

data "powerplatform_connectors" "all_connectors" {}