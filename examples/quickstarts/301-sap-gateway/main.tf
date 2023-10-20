terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">=3.74.0"
    }
    azurecaf = {
      source  = "aztfmod/azurecaf"
      version = ">=1.2.26"
    }
  }
}

provider "azurerm" {
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
    key_vault {
      purge_soft_delete_on_destroy    = true
      recover_soft_deleted_key_vaults = false
    }
  }
  client_id       = var.client_id_gw
  client_secret   = var.secret_gw
  tenant_id       = var.tenant_id_gw
  subscription_id = var.subscription_id_gw
}

data "azurerm_client_config" "current" {}

resource "azurecaf_name" "rg" {
  name          = var.base_name
  resource_type = "azurerm_resource_group"
  prefixes      = [var.prefix]
  random_length = 3
  clean_input   = true
}

resource "azurerm_resource_group" "rg" {
  name     = azurecaf_name.rg.result
  location = var.region_gw
}

resource "azurecaf_name" "vnet" {
  name          = var.base_name
  resource_type = "azurerm_virtual_network"
  prefixes      = [var.prefix]
  random_length = 3
  clean_input   = true
}

resource "azurerm_virtual_network" "vnet" {
  name                = azurecaf_name.vnet.result
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
}

resource "azurecaf_name" "subnet" {
  name          = var.base_name
  resource_type = "azurerm_subnet"
  prefixes      = [var.prefix]
  random_length = 3
  clean_input   = true
}

resource "azurerm_subnet" "subnet" {
  name                 = azurecaf_name.subnet.result
  resource_group_name  = azurerm_resource_group.rg.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azurecaf_name" "nsg" {
  name          = var.base_name
  resource_type = "azurerm_network_security_group"
  prefixes      = [var.prefix]
  random_length = 3
  clean_input   = true
}

resource "azurerm_network_security_group" "nsg" {
  name                = azurecaf_name.nsg.result
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  security_rule {
    name                       = "SSH"
    priority                   = 1001
    direction                  = "Inbound"
    access                     = "Deny"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "22"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }

  security_rule {
    name                       = "HTTP"
    priority                   = 1002
    direction                  = "Inbound"
    access                     = "Deny"
    protocol                   = "Tcp"
    source_port_range          = "*"
    destination_port_range     = "80"
    source_address_prefix      = "*"
    destination_address_prefix = "*"
  }
}

resource "azurecaf_name" "publicip" {
  name          = var.base_name
  resource_type = "azurerm_public_ip"
  prefixes      = [var.prefix]
  random_length = 3
  clean_input   = true
}

resource "azurerm_public_ip" "publicip" {
  name                = azurecaf_name.publicip.result
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  allocation_method   = "Dynamic"
}

resource "azurecaf_name" "nic" {
  name          = var.base_name
  resource_type = "azurerm_network_interface"
  prefixes      = [var.prefix]
  random_length = 3
  clean_input   = true
}

resource "azurerm_network_interface" "nic" {
  name                = azurecaf_name.nic.result
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name

  ip_configuration {
    name                          = "internal"
    subnet_id                     = azurerm_subnet.subnet.id
    private_ip_address_allocation = "Dynamic"
    public_ip_address_id          = azurerm_public_ip.publicip.id
  }

}

resource "azurerm_network_interface_security_group_association" "rgassociation" {
  network_interface_id      = azurerm_network_interface.nic.id
  network_security_group_id = azurerm_network_security_group.nsg.id
}

resource "random_string" "key_vault_suffix" {
  length  = 3
  upper   = false
  numeric = false
  special = false
}

resource "azurecaf_name" "key_vault" {
  name = "${var.prefix}-${random_string.key_vault_suffix.result}"
  #name          = var.base_name
  resource_type = "azurerm_key_vault"
  prefixes      = [var.prefix]
  random_length = 3
  clean_input   = true
}


output "name" {
  value = "${var.prefix}-${random_string.key_vault_suffix.result}"
}

resource "azurecaf_name" "key_vault2" {
  name = var.prefix
  #name          = var.base_name
  resource_type = "azurerm_key_vault"
  #prefixes      = [var.base_name]
  random_length = 3
  clean_input   = true
}


output "name2" {
  value = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault" "key_vault" {
  name                = azurecaf_name.key_vault2.result
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  tenant_id           = var.tenant_id_gw
  sku_name            = "standard"

  access_policy {
    tenant_id = var.tenant_id_gw
    object_id = data.azurerm_client_config.current.object_id

    secret_permissions = [
      "Get",
      "List",
      "Delete",
      "Set",
      "Purge",
    ]
  }
}

resource "azurecaf_name" "key_vault_secret" {
  name          = var.base_name
  resource_type = "azurerm_key_vault_secret"
  prefixes      = [var.prefix]
  random_length = 3
  clean_input   = true
}

resource "azurerm_key_vault_secret" "key_vault_secret" {
  name         = azurecaf_name.key_vault_secret.result
  value        = var.secret_pp
  key_vault_id = azurerm_key_vault.key_vault.id
}

module "gateway_vm" {
  source              = "./gateway-vm"
  resource_group_name = azurerm_resource_group.rg.name
  base_name           = var.base_name
  region              = var.region_gw
  vm_pwd              = var.vm_pwd_gw
  nic_id              = azurerm_network_interface.nic.id
  keyVaultUri         = azurerm_key_vault.key_vault.vault_uri
  secretPPName        = azurerm_key_vault_secret.key_vault_secret.name
  userIdAdmin_pp      = var.userIdAdmin_pp
  ps7_setup_link      = var.ps7_setup_link
  java_setup_link     = var.java_setup_link
  opdgw_install_link  = var.opdgw_install_link
  opdgw_setup_link    = var.opdgw_setup_link
  sapnco_install_link = var.sapnco_install_link
}

resource "azurerm_key_vault_access_policy" "key_vault_access_policy" {
  key_vault_id = azurerm_key_vault.key_vault.id
  tenant_id    = var.tenant_id_gw
  object_id    = module.gateway_vm.vm_opgw_id
  secret_permissions = [
    "Get",
    "List",
  ]
}

# It will be included in futures releases.
#module "gateway_principal" {
#  source    = "./gateway-principal"
#  client_id = var.client_id_pp
#  secret    = var.secret_pp
#  tenant_id = var.tenant_id_pp
#}

