terraform {
  required_providers {
    power-platform = {
      source = "microsoft/power-platform"
    }
    local = {
      version = "2.5.1"
      source  = "hashicorp/local"
    }
  }
}

provider "power-platform" {
  use_cli = true
}

provider "local" {}

resource "local_file" "solution_settings_file" {
  filename = "${path.module}/solution_settings.json"
  content  = <<EOF
{
  "EnvironmentVariables": [
    {
      "SchemaName": "cra6e_SolutionVariableDataSource",
      "Value": "/sites/Shared%20Documents1"
    },
    {
      "SchemaName": "cra6e_SolutionVariableJson",
      "Value": "{ \"value\": 1234, \"text\": \"abc\" }"
    },
    {
      "SchemaName": "cra6e_SolutionVariableText",
      "Value": "${powerplatform_environment.environment.id}"
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

resource "powerplatform_environment" "environment" {
  display_name     = "Solution Import Test 1"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "powerplatform_solution" "solution" {
  environment_id = powerplatform_environment.environment.id
  solution_file  = "${path.module}/${var.solution_name}_Complex_1_1_0_0.zip"
  solution_name  = var.solution_name
  settings_file  = local_file.solution_settings_file.filename
  //settings_file  = "${path.module}/solution_settings_static.json"
}
