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

data "powerplatform_client_config" "current" {}

data "powerplatform_tenant_capacity" "capacity" {
  tenant_id = data.powerplatform_client_config.current.tenant_id
}

output "tenant_capacity" {
  value = data.powerplatform_tenant_capacity.capacity
}
