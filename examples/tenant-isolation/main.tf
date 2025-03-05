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

resource "powerplatform_tenant_isolation_policy" "example" {
  allowed_tenants = toset([
    {
      tenant_id = "11111111-1111-1111-1111-111111111111"
      inbound   = true
      outbound  = true
    }
  ])
  is_disabled = false
}
