terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

variable "environment_id" {
  type = string
}

variable "name" {
  type = string
}

variable "costcenter" {
  type = string
}

variable "parent_business_unit_id" {
  type = string
}

resource "powerplatform_data_record" "business_unit" {
  environment_id     = var.environment_id
  table_logical_name = "businessunit"
  disable_on_destroy = true
  columns = {
    name       = var.name
    costcenter = var.costcenter
    parentbusinessunitid = {
      table_logical_name = "businessunit"
      data_record_id     = var.parent_business_unit_id
    }
  }
}

output "resource_id" {
  value = powerplatform_data_record.business_unit.id
}

output "resource" {
  value = powerplatform_data_record.business_unit
}
