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

data "powerplatform_tenant_capacity" "capacity" {
  tenant_id = "4481d6dc-0f72-4841-a3a9-0c8f9798d2d6"
}

output "tenant_capacity" {
  value = data.powerplatform_tenant_capacity.capacity
}
