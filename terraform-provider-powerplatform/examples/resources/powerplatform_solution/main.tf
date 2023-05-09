terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/powerplatform"
    }
    local = {
      source = "hashicorp/local"
    }
  }
}

provider "powerplatform" {
  username = var.username
  password = var.password
  host     = var.host
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
      "Value": "{ \"value\": 123, \"text\": \"abc\" }"
    },
    {
      "SchemaName": "cra6e_SolutionVariableText",
      "Value": "${powerplatform_environment.environment.environment_name}"
    }
  ],
  "ConnectionReferences": [
    {
      "LogicalName": "cra6e_ConnectionReferenceSharePoint",
      "ConnectionId": "123",
      "ConnectorId": "/providers/Microsoft.PowerApps/apis/shared_sharepointonline"
    }
  ]
}
EOF
}

resource "powerplatform_environment" "environment" {
  display_name     = "Solution Import Test"
  location         = "europe"
  language_name    = "1033"
  currency_name    = "USD"
  environment_type = "Sandbox"
}

resource "powerplatform_solution" "solution" {
  environment_name = powerplatform_environment.environment.environment_name
  solution_file    = "${path.module}/${var.solution_name}_Complex_1_1_0_0.zip"
  settings_file    = local_file.solution_settings_file.filename
  solution_name    = var.solution_name
}
