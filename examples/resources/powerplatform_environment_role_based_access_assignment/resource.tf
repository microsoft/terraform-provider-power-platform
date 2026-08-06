terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  use_cli = true
}

resource "powerplatform_environment" "example_environment" {
  display_name     = "example_environment"
  location         = "europe"
  environment_type = "Sandbox"
}

# Fetch all available role definitions so we can look up the one we need by name
data "powerplatform_role_definitions" "all" {
}

variable "role_definition_name" {
  default     = "Power Platform Role Based Access Control Administrator"
  description = "Display name of the role definition to assign"
  type        = string
}

variable "enterprise_application_object_id" {
  default     = "00000000-0000-0000-0000-000000000000"
  description = "Object id of the enterprise application that will be granted the role"
  type        = string
}

locals {
  role_definition_id = [
    for role in data.powerplatform_role_definitions.all.role_definitions :
    role.role_definition_id if role.role_definition_name == var.role_definition_name
  ][0]
}

# Assign a role to a service principal at the environment level
resource "powerplatform_environment_role_based_access_assignment" "example" {
  environment_id                   = powerplatform_environment.example_environment.id
  enterprise_application_object_id = var.enterprise_application_object_id
  principal_type                   = "ApplicationUser"
  role_definition_id               = local.role_definition_id
}
