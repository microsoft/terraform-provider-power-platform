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

# Fetch all available role definitions so we can look up the one we need by name
data "powerplatform_role_definitions" "all" {
}

variable "role_definition_name" {
  default     = "Power Platform Role Based Access Control Administrator"
  description = "Display name of the role definition to assign"
  type        = string
}

variable "principal_id" {
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

# Assign a role to a service principal at the tenant level
resource "powerplatform_role_assignment" "example" {
  scope_type         = "tenant"
  principal_id       = var.principal_id
  principal_type     = "ApplicationUser"
  role_definition_id = local.role_definition_id
}

# The same resource scopes the assignment by which identifier you set.
# Set environment_id for an environment:
resource "powerplatform_role_assignment" "environment" {
  scope_type         = "environment"
  environment_id     = var.environment_id
  principal_id       = var.principal_id
  principal_type     = "ApplicationUser"
  role_definition_id = local.role_definition_id
}

# Set environment_group_id for an environment group:
resource "powerplatform_role_assignment" "environment_group" {
  scope_type           = "environment_group"
  environment_group_id = var.environment_group_id
  principal_id         = var.principal_id
  principal_type       = "ApplicationUser"
  role_definition_id   = local.role_definition_id
}

variable "environment_id" {
  description = "Id of the environment to scope an assignment to"
  type        = string
}

variable "environment_group_id" {
  description = "Id of the environment group to scope an assignment to"
  type        = string
}
