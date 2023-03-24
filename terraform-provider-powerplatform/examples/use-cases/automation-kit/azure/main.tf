terraform {
  required_providers {
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 2.15.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0.0"
    }
  }
}

#register resource provider for a subscription 
#if this is already done, it will thorw an error
# resource "azurerm_resource_provider_registration" "azure_provider_registration" {
#   name = "Microsoft.PowerPlatform"
# }


#Add service principal to AAD
data "azuread_application_published_app_ids" "well_known" {}

resource "azuread_application" "automation_kit_app_user" {
  display_name = "CoE_Automation_Kit_App_User"

  required_resource_access {
      resource_app_id = data.azuread_application_published_app_ids.well_known.result.DynamicsCrm
      resource_access {
          id = "78ce3f0f-a1ce-49c2-8cde-64b5c0896db4" #"user_impersonation"
          type = "Scope"
      }
  }

  web {
    redirect_uris = ["https://localhost/"]
  }

}

resource "azuread_service_principal" "app_user_principal" {
  application_id = azuread_application.automation_kit_app_user.application_id
}


#Create resource group for all resources
resource "azurerm_resource_group" "resource_group" {
  name     = "${var.resource_group_name}"
  location = "${var.resource_group_location}"
}

resource "azurerm_key_vault" "kv_dev" {
  name                = "${var.key_vault_name}"
  location            = azurerm_resource_group.resource_group.location
  resource_group_name = azurerm_resource_group.resource_group.name
  tenant_id           = var.tenant_id
  sku_name            = "premium"
}

resource "azurerm_key_vault_access_policy" "kv_access_policy_for_dataverse" {
  key_vault_id = azurerm_key_vault.kv_dev.id
  tenant_id    = var.tenant_id
  object_id    = "899abecf-d28a-4f3b-8a1a-4d05c2b907fb" #dataverse service principal

  secret_permissions = [
    "List",
    "Get",
  ]
}

resource "azurerm_key_vault_access_policy" "kv_access_policy_for_terraform" {
  key_vault_id = azurerm_key_vault.kv_dev.id

  tenant_id = var.tenant_id
  object_id = var.terraform_service_principal_object_id

  secret_permissions = [
    "Delete",
    "List",
    "Get",
    "Set",
    "Purge",
    "Recover",
  ]
}

resource "azurerm_key_vault_secret" "client_id" {
  name         = "KVS-AutomationCoE-ClientID"
  value        = "${azuread_application.automation_kit_app_user.application_id}"
  key_vault_id = azurerm_key_vault.kv_dev.id

  depends_on = [azurerm_key_vault_access_policy.kv_access_policy_for_terraform]
}

resource "azurerm_key_vault_secret" "tenant_id" {
  name         = "KVS-AutomationCoE-TenantID"
  value        = "${var.tenant_id}"
  key_vault_id = azurerm_key_vault.kv_dev.id

  depends_on = [azurerm_key_vault_access_policy.kv_access_policy_for_terraform]
}

resource "azurerm_key_vault_secret" "app_user_secret" {
  name         = "KVS-AutomationCoE-Secret"
  value        = ""
  key_vault_id = azurerm_key_vault.kv_dev.id

  depends_on = [azurerm_key_vault_access_policy.kv_access_policy_for_terraform]
}