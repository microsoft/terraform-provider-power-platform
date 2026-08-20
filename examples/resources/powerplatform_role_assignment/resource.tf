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

# Role definitions carry stable ids; display names have been recased before, so match on the id.
# powerplatform_role_definitions lists them all.
locals {
  # Power Platform Reader
  role_definition_id = "c886ad2e-27f7-4874-8381-5849b8d8a090"
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
