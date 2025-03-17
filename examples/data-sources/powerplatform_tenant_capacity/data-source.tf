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

data "powerplatform_tenant" "current" {}

data "powerplatform_tenant_capacity" "capacity" {
  tenant_id = data.powerplatform_tenant.current.tenant_id
}
