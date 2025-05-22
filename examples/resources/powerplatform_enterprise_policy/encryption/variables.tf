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
