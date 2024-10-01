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
  columns = {
    name       = var.name
    costcenter = var.costcenter
    parentbusinessunitid = {
      table_logical_name = "businessunit"
      data_record_id     = var.parent_business_unit_id
    }
  }
}

data "powerplatform_environments" "envs" {
}

locals {
  dataverse_url = tostring(one([for e in data.powerplatform_environments.envs.environments : e.dataverse.url if e.id == var.environment_id]))
}

resource "powerplatform_rest" "business_unit_disable_on_destroy" {
  destroy = {
    scope                = "${local.dataverse_url}.default"
    method               = "PATCH"
    url                  = "${local.dataverse_url}api/data/v9.2/businessunits(${powerplatform_data_record.business_unit.id})"
    expected_http_status = [200, 204]
    body = jsonencode({
      isdisabled    = true
    })
  }
  depends_on = [powerplatform_data_record.business_unit]
}

output "resource_id" {
  value = powerplatform_data_record.business_unit.id
}

output "resource" {
  value = powerplatform_data_record.business_unit
}
