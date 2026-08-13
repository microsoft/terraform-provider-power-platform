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

resource "powerplatform_environment" "environment" {
  display_name     = "Managed Solution Import Test"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "powerplatform_connection" "dataverse_connection" {
  environment_id = powerplatform_environment.environment.id
  name           = "shared_commondataserviceforapps"
  display_name   = "Dataverse Connection"

  connection_parameters = jsonencode({
  })

  lifecycle {
    ignore_changes = [
      connection_parameters
    ]
  }
}

resource "powerplatform_managed_solution" "solution" {
  environment_id = powerplatform_environment.environment.id
  unique_name    = "TerraformSolutionExample"
  version        = "1.0.0.1"

  source = {
    path = "${path.module}/TerraformSolutionExample_1_0_0_1_managed.zip"
  }

  connection_references = {
    terr_SolutionConnectionReference = powerplatform_connection.dataverse_connection.id
  }

  # Environment variable values are owned by their own resource and set after the
  # import that ships the definitions, ordered with an ordinary depends_on.

  depends_on = [powerplatform_connection.dataverse_connection]
}
