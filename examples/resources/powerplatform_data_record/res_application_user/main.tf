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

variable "application_id" {
  description = "EntraId clientid of the application"
  type        = string
  validation {
    condition     = length(var.application_id) > 0
    error_message = "The application id must not be empty"
  }
}

variable "business_unit_id" {
  description = "Unique identifier of the business unit"
  type        = string
  validation {
    condition     = length(var.business_unit_id) > 0
    error_message = "The business unit id must not be empty"
  }
}

variable "role_ids" {
  type        = set(string)
  description = "The role ids that are granted to this application user"
}


resource "powerplatform_data_record" "app_user" {
  table_logical_name = "systemuser"
  environment_id     = var.environment_id
  disable_on_destroy = true # Application Users cannot be deleted without being disabled first
  columns = {
    applicationid = var.application_id
    businessunitid = {
      table_logical_name = "businessunit"
      data_record_id     = var.business_unit_id
    }
    systemuserroles_association =toset([for rid in var.role_ids : { table_logical_name = "role", data_record_id = tostring(rid) }])
  }
}

output "application_user_id" {
  value = powerplatform_data_record.app_user.id
}
