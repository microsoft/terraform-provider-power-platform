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

variable "team_name" {
  description = "The name of the team"
  type        = string
  validation {
    condition     = length(var.team_name) > 0
    error_message = "The team name must not be empty"
  }
}

variable "team_description" {
  description = "The description of the team"
  type        = string
}

variable "role_ids" {
  type        = set(string)
  description = "The role ids that are granted to this team"
  
}

resource "powerplatform_data_record" "team" {
  environment_id     = var.environment_id
  table_logical_name = "team"
  columns = {
    name        = var.team_name
    description = var.team_description

    teamroles_association = tolist([for rid in var.role_ids : { table_logical_name = "role", data_record_id = tostring(rid) }])
  }
}

output "team_id" {
  value = powerplatform_data_record.team.id
}
