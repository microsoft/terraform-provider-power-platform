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

# Known limitation: Dataverse Git integration currently works only with delegated
# user principal authentication that also has Azure DevOps repository access.
# Service principal, app-only, and OIDC pipeline identities are not supported.

# Use `scope = "Environment"` to mirror the maker UI environment-level binding.
# In this mode the provider manages the root Dataverse binding and proactively
# enables eligible visible unmanaged solutions in the environment. Built-in
# platform solutions are excluded automatically.
resource "powerplatform_environment" "example" {
  display_name     = "example-git-integration-environment"
  description      = "Example environment for validating Dataverse Git integration."
  location         = "europe"
  azure_region     = "northeurope"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "powerplatform_environment_git_integration" "example" {
  environment_id    = powerplatform_environment.example.id
  scope             = "Environment"
  organization_name = "contoso-org"
  project_name      = "PowerPlatform Solutions"
  repository_name   = "power-platform-solutions"
}
