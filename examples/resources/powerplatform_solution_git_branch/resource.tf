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

# The environment Git integration must use `scope = "Solution"` so that
# individual solutions can be bound to Git branches.
resource "powerplatform_environment_git_integration" "example" {
  environment_id    = "00000000-0000-0000-0000-000000000001"
  scope             = "Solution"
  organization_name = "contoso-org"
  project_name      = "PowerPlatform Solutions"
  repository_name   = "power-platform-solutions"
}

# `solution_id` uses the provider solution ID format
# (`<environment_id>_<dataverse_solution_id>`) of an existing unmanaged
# solution in the same environment, such as the `id` of a
# `powerplatform_solution` resource.
resource "powerplatform_solution_git_branch" "example" {
  environment_id     = "00000000-0000-0000-0000-000000000001"
  git_integration_id = powerplatform_environment_git_integration.example.id
  solution_id        = "00000000-0000-0000-0000-000000000001_22222222-2222-2222-2222-222222222222"
  branch_name        = "main"
  root_folder_path   = "solutions/sample-solution"
}
