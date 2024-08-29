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

data "powerplatform_tenant" "current_tenant" {}

output "current_config" {
  value = data.powerplatform_tenant.current_tenant
}
