terraform {
  required_version = "> 1.7.0"
  required_providers {
    powerplatform = {
      source  = "microsoft/power-platform"
      version = "~>4.0"
    }
    azapi = {
      source  = "azure/azapi"
      version = "~>2.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~>4.8"
    }
  }
}

variable "environment_id" {
  description = "The ID of the environment"
  type        = string
  validation {
    condition     = length(var.environment_id) > 0
    error_message = "The environment ID must not be empty"
  }
}

variable "should_register_provider" {
  description = "A flag to determine if the PowerPlatfomr provider should be registered in the subscription"
  type        = bool
  default     = true
}

variable "resource_group_name" {
  description = "The name of the resource group"
  type        = string
  validation {
    condition     = length(var.resource_group_name) > 0
    error_message = "The resource group name must not be empty"
  }
}

variable "resource_group_location" {
  description = "The location of the resource group"
  type        = string
  validation {
    condition     = length(var.resource_group_location) > 0
    error_message = "The resource group location must not be empty"
  }

}

variable "enterprise_policy_name" {
  description = "The name of the enterprise policy"
  type        = string
  validation {
    condition     = length(var.enterprise_policy_name) > 0
    error_message = "The enterprise policy name must not be empty"
  }
}

variable "enterprise_policy_location" {
  description = "The location of the enterprise policy"
  type        = string
  validation {
    condition     = length(var.enterprise_policy_location) > 0
    error_message = "The enterprise policy location must not be empty"
  }
}

variable "keyvault_name" {
  description = "The name of the key vault"
  type        = string
  validation {
    condition     = length(var.keyvault_name) > 0
    error_message = "The key vault name must not be empty"
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

data "azurerm_client_config" "current" {
}

resource "azurerm_key_vault" "key_vault" {
  name                        = var.keyvault_name
  location                    = azurerm_resource_group.resource_group.location
  resource_group_name         = azurerm_resource_group.resource_group.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = true

  sku_name = "standard"

  enable_rbac_authorization = true

  access_policy = []
}

resource "azurerm_role_assignment" "terraform_admin_access" {
  scope                = azurerm_key_vault.key_vault.id
  role_definition_name = "Key Vault Administrator"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azurerm_key_vault_key" "kv_ep_key" {
  name         = "generated-certificate"
  key_vault_id = azurerm_key_vault.key_vault.id
  key_type     = "RSA"
  key_size     = 2048

  key_opts = [
    "decrypt",
    "encrypt",
    "sign",
    "unwrapKey",
    "verify",
    "wrapKey",
  ]
  depends_on = [azurerm_role_assignment.terraform_admin_access]
}

resource "azapi_resource" "powerplatform_policy" {
  schema_validation_enabled = false

  type      = "Microsoft.PowerPlatform/enterprisePolicies@2020-10-30-preview"
  name      = var.enterprise_policy_name
  location  = var.enterprise_policy_location
  parent_id = azurerm_resource_group.resource_group.id
  body = {
    identity = {
      type = "SystemAssigned"
    }
    properties = {
      encryption = {
        keyVault = {
          id = azurerm_key_vault.key_vault.id
          key = {
            name    = azurerm_key_vault_key.kv_ep_key.name
            version = azurerm_key_vault_key.kv_ep_key.version
          }
        }
        state = "Enabled"
      }
    }
    kind = "Encryption"
  }
}

//we have to wait as managed identity created with enterprise policy is not available immediately
resource "time_sleep" "wait_90_seconds" {
  depends_on = [azapi_resource.powerplatform_policy]

  create_duration = "90s"
}

data "azapi_resource_action" "managed_identity_query" {
  type        = "Microsoft.ResourceGraph@2021-03-01"
  resource_id = "/providers/Microsoft.ResourceGraph"
  action      = "resources"
  body = {
    query = <<-KQL
resources
| where id == "${azapi_resource.powerplatform_policy.output.id}"
| take 1
    KQL
  }
  response_export_values = ["*"]

  depends_on = [time_sleep.wait_90_seconds]
}
output "o1" {
  value = data.azapi_resource_action.managed_identity_query.output.data[0].identity.principalId
}

resource "azurerm_role_assignment" "enterprise_policy_system_access" {
  scope                = azurerm_key_vault.key_vault.id
  role_definition_name = "Key Vault Crypto Service Encryption User"
  principal_id         = data.azapi_resource_action.managed_identity_query.output.data[0].identity.principalId
}

resource "powerplatform_enterprise_policy" "encryption" {
  environment_id = var.environment_id
  system_id      = azapi_resource.powerplatform_policy.output.properties.systemId
  policy_type    = "Encryption"
}

output "enterprise_policy_system_id" {
  value = azapi_resource.powerplatform_policy.output.properties.systemId
}

output "enterprise_policy_id" {
  value = azapi_resource.powerplatform_policy.output.id
}
