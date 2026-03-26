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

# Use `scope = "Environment"` to mirror the maker UI environment-level binding.
# In this mode the provider manages the root Dataverse binding and proactively
# enables eligible visible unmanaged solutions in the environment. Built-in
# platform solutions are excluded automatically.
resource "powerplatform_environment" "example" {
  display_name     = var.environment_display_name
  description      = "Example environment for validating Dataverse Git integration."
  location         = var.location
  azure_region     = var.azure_region
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = var.security_group_id
  }
}

resource "powerplatform_environment_git_integration" "example" {
  environment_id    = powerplatform_environment.example.id
  git_provider      = var.git_provider
  scope             = var.scope
  organization_name = var.organization_name
  project_name      = var.project_name
  repository_name   = var.repository_name
}
