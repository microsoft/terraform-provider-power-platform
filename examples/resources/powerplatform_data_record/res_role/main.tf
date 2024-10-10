terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

variable "environment_id" {
  description = "The unique identifier of the environment"
  type        = string
  validation {
    condition     = length(var.environment_id) > 0
    error_message = "The environment id must not be empty"
  }
}

variable "role_name" {
  description = "The name of the role"
  type        = string
  validation {
    condition     = length(var.role_name) > 0
    error_message = "The role name must not be empty"
  }
}

variable "business_unit_id" {
  description = "The unique identifier of the business unit"
  type        = string
  validation {
    condition     = length(var.business_unit_id) > 0
    error_message = "The business unit id must not be empty"
  }
}

resource "powerplatform_data_record" "role" {
  environment_id     = var.environment_id
  table_logical_name = "role"

  columns = {
    name = var.role_name

    businessunitid = {
      table_logical_name = "businessunit"
      data_record_id     = var.business_unit_id
    }
  }
}

output "role_id" {
  value = powerplatform_data_record.role.id
}
