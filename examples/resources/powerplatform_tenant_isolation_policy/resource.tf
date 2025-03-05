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

resource "powerplatform_tenant_isolation_policy" "example" {
  is_disabled = false
  allowed_tenants = [
    # Allow all inbound and outbound connections with a specific tenant
    {
      tenant_id = "11111111-1111-1111-1111-111111111111"
      inbound   = true
      outbound  = true
    },
    # Allow only inbound connections from another tenant
    {
      tenant_id = "22222222-2222-2222-2222-222222222222"
      inbound   = true
      outbound  = false
    },
    # Allow only outbound connections to a third tenant
    {
      tenant_id = "33333333-3333-3333-3333-333333333333"
      inbound   = false
      outbound  = true
    },
    # Allow all tenants (wildcard) for inbound connections only
    {
      tenant_id = "*"
      inbound   = true
      outbound  = false
    }
  ]
}
