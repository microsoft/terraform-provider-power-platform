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

variable "policy_reader_object_id" {
  description = "The object ID of the user to assign the Reader role to (needed to finish the setup of the Azure Synapse Link to Dataverse)"
  type        = string
  validation {
    condition     = length(var.policy_reader_object_id) > 0
    error_message = "The policy reader object ID must not be empty"
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

resource "azurerm_role_assignment" "policy_reader" {
  scope                = azapi_resource.powerplatform_policy.id
  role_definition_name = "Reader"
  principal_id         = var.policy_reader_object_id
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
    kind = "Identity"
  }
}

resource "powerplatform_enterprise_policy" "identity_policy" {
  environment_id = var.environment_id
  system_id      = azapi_resource.powerplatform_policy.output.properties.systemId
  policy_type    = "Identity"
}

output "enterprise_policy_system_id" {
  value = azapi_resource.powerplatform_policy.output.properties.systemId
}

output "enterprise_policy_id" {
  value = azapi_resource.powerplatform_policy.output.id
}

