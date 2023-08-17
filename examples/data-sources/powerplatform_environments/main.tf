terraform {
  required_providers {
    powerplatform = {
      source  = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  username = var.username
  password = var.password
  tenant_id = var.tenant_id
}

data "powerplatform_environments" "all_environments" {}