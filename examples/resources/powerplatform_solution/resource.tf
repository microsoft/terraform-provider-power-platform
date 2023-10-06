terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/power-platform"
    }
    local = {
      version = "2.4.0"
      source  = "hashicorp/local"
    }
  }
}

provider "powerplatform" {
  username  = var.username
  password  = var.password
  tenant_id = var.tenant_id
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
  display_name      = "Solution Import Test 1"
  location          = "europe"
  language_code     = "1033"
  currency_code     = "USD"
  environment_type  = "Sandbox"
  domain            = random_string.random_domain.result
  security_group_id = "00000000-0000-0000-0000-000000000000"
}

resource "powerplatform_solution" "solution" {
  environment_id = powerplatform_environment.environment.id
  solution_file    = "${path.module}/${var.solution_name}_Complex_1_1_0_0.zip"
  solution_name    = var.solution_name
  settings_file    = local_file.solution_settings_file.filename
  //settings_file  = "${path.module}/solution_settings_static.json"
}
