terraform {
  required_version = "> 1.7.0"
  required_providers {
    powerplatform = {
      source  = "microsoft/power-platform"
      version = "~>3.7.2"
    }
    azapi = {
      source  = "azure/azapi"
      version = "~>2.2.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~>4.16.0"
    }
  }
}


resource "azurerm_resource_group" "resource_group" {
  name     = var.resource_group_name
  location = var.resource_group_location
}

resource "azurerm_resource_provider_registration" "provider_registration" {
  count = var.should_register_provider ? 1 : 0
  name  = "Microsoft.PowerPlatform"
}

locals {
  vnet_names = tolist(["vnet_primary", "vnet_secondary"])
}

resource "azurerm_virtual_network" "vnet" {
  count               = 2
  name                = local.vnet_names[count.index]
  location            = var.vnet_locations[count.index]
  resource_group_name = azurerm_resource_group.resource_group.name
  address_space       = ["10.0.0.0/16"]
}

resource "azurerm_subnet" "subnet" {
  count                = 2
  name                 = "enterprise_policy_subnet"
  resource_group_name  = var.resource_group_name
  virtual_network_name = azurerm_virtual_network.vnet[count.index].name
  address_prefixes     = ["10.0.1.0/24"]

  delegation {
    name = "delegation"

    service_delegation {
      name    = "Microsoft.PowerPlatform/enterprisePolicies"
      actions = ["Microsoft.Network/virtualNetworks/subnets/join/action"]
    }
  }
}

resource "azapi_resource" "powerplatform_policy" {
  schema_validation_enabled = false

  type      = "Microsoft.PowerPlatform/enterprisePolicies@2020-10-30-preview"
  name      = var.enterprise_policy_name
  location  = var.enterprise_policy_location
  parent_id = azurerm_resource_group.resource_group.id
  body = {
    properties = {
      networkInjection = {
        virtualNetworks = [
          {
            id = azurerm_virtual_network.vnet[0].id
            subnet = {
              name = azurerm_subnet.subnet[0].name
            }
          },
          {
            id = azurerm_virtual_network.vnet[1].id

            subnet = {
              name = azurerm_subnet.subnet[1].name
            }
          }
        ]
      }
    }
    kind = "NetworkInjection"
  }
}

resource "powerplatform_enterprise_policy" "network_injection" {
  environment_id = var.environment_id
  system_id      = azapi_resource.powerplatform_policy.output.properties.systemId
  policy_type    = "NetworkInjection"

  #depends_on = [powerplatform_managed_environment.managed_development]
}

output "enterprise_policy_system_id" {
  value = azapi_resource.powerplatform_policy.output.properties.systemId
}

output "enterprise_policy_id" {
  value = azapi_resource.powerplatform_policy.output.id
}
