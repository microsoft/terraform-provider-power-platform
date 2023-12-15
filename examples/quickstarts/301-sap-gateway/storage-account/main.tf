terraform {
  required_version = ">= 1.5"
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
    }
    azurecaf = {
      source  = "aztfmod/azurecaf"
      version = ">=1.2.26"
    }
  }
}

resource "azurecaf_name" "storage_account_name" {
  name          = var.base_name
  resource_type = "azurerm_storage_account"
  prefixes      = [var.prefix]
  random_length = 3
  clean_input   = true
}

resource "azurerm_storage_account" "storage_account" {
  name                     = azurecaf_name.storage_account_name.result
  resource_group_name      = var.resource_group_name
  location                 = var.region
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_container" "storage_container_installs" {
  name                  = "installs"
  storage_account_name  = azurerm_storage_account.storage_account.name
  container_access_type = "blob"
}

resource "azurerm_storage_blob" "storage_blob_ps7_setup" {
  name                   = "ps7-setup.ps1"
  storage_account_name   = azurerm_storage_account.storage_account.name
  storage_container_name = azurerm_storage_container.storage_container_installs.name
  type                   = "Block"
  source                 = "/workspaces/terraform-provider-power-platform/examples/quickstarts/301-sap-gateway/storage-account/scripts/ps7-setup.ps1"
}

resource "azurerm_storage_blob" "storage_blob_java_runtime" {
  name                   = "java-runtime-setup.ps1"
  storage_account_name   = azurerm_storage_account.storage_account.name
  storage_container_name = azurerm_storage_container.storage_container_installs.name
  type                   = "Block"
  source                 = "/workspaces/terraform-provider-power-platform/examples/quickstarts/301-sap-gateway/storage-account/scripts/java-setup.ps1"
}

resource "azurerm_storage_blob" "storage_blob_sapnco_install" {
  name                   = "sapnco.msi"
  storage_account_name   = azurerm_storage_account.storage_account.name
  storage_container_name = azurerm_storage_container.storage_container_installs.name
  type                   = "Block"
  source                 = "/workspaces/terraform-provider-power-platform/examples/quickstarts/301-sap-gateway/storage-account/sapnco-msi/sapnco.msi"
}

resource "azurerm_storage_blob" "storage_blob_runtime_setup" {
  name                   = "runtime-setup.ps1"
  storage_account_name   = azurerm_storage_account.storage_account.name
  storage_container_name = azurerm_storage_container.storage_container_installs.name
  type                   = "Block"
  source                 = "/workspaces/terraform-provider-power-platform/examples/quickstarts/301-sap-gateway/storage-account/scripts/runtime-setup.ps1"
}


