terraform {
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "2.6.2"
    }
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

provider "local" {}

provider "powerplatform" {
  use_cli = true
}

resource "local_file" "solution_settings_file" {
  filename = "${path.module}/solution_settings.json"
  content  = <<EOF
{
  "EnvironmentVariables": [
    {
      "SchemaName": "cra6e_SolutionVariableDataSource",
      "Value": "/sites/Shared Documents"
    },
    {
      "SchemaName": "cra6e_SolutionVariableJson",
      "Value": "{ \"value\": 1234, \"text\": \"abc\" }"
    },
    {
      "SchemaName": "cra6e_SolutionVariableText",
      "Value": "${powerplatform_environment.example.id}"
    }
  ],
  "ConnectionReferences": [
    {
      "LogicalName": "cra6e_ConnectionReferenceSharePoint",
      "ConnectionId": "00000000-0000-0000-0000-000000000000",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_sharepointonline"
    }
  ]
}
EOF
}

resource "powerplatform_environment" "example" {
  display_name     = var.environment_display_name
  description      = "Example environment for validating Dataverse Git branch bindings."
  location         = var.location
  azure_region     = var.azure_region
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = var.security_group_id
  }
}

resource "powerplatform_solution" "example" {
  environment_id = powerplatform_environment.example.id
  solution_file  = coalesce(var.solution_file, "${path.module}/../powerplatform_solution/TerraformTestSolution_Complex_1_1_0_0.zip")
  settings_file  = local_file.solution_settings_file.filename
}

resource "powerplatform_environment_git_integration" "example" {
  count = var.enable_git_binding ? 1 : 0

  environment_id    = powerplatform_environment.example.id
  git_provider      = var.git_provider
  scope             = var.scope
  organization_name = var.organization_name
  project_name      = var.project_name
  repository_name   = var.repository_name
}

resource "powerplatform_solution_git_branch" "example" {
  count = var.enable_git_binding ? 1 : 0

  environment_id       = powerplatform_environment.example.id
  git_integration_id   = powerplatform_environment_git_integration.example[0].id
  solution_id          = powerplatform_solution.example.id
  branch_name          = var.branch_name
  upstream_branch_name = var.upstream_branch_name
  root_folder_path     = var.root_folder_path
}
