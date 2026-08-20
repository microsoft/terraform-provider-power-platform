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

variable "principal_id" {
  default     = "00000000-0000-0000-0000-000000000000"
  description = "Object id of the enterprise application that will be granted the role"
  type        = string
}

# The role is selected by exactly one of role_definition_name or role_definition_id.
# Names are matched case-insensitively and resolved to the id at create time;
# powerplatform_role_definitions lists the whole catalogue.

# Assign a role to a service principal at the tenant level, selecting the role by name
resource "powerplatform_role_assignment" "example" {
  scope_type           = "tenant"
  principal_id         = var.principal_id
  principal_type       = "ApplicationUser"
  role_definition_name = "Power Platform Reader"
}

# An id pins the role independently of its display name
locals {
  # Power Platform Reader
  role_definition_id = "c886ad2e-27f7-4874-8381-5849b8d8a090"
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
