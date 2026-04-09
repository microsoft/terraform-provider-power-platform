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

resource "powerplatform_managed_solution" "solution" {
  environment_id = powerplatform_environment.environment.id
  unique_name    = "TerraformSimpleTestSolution"
  version        = "1.0.0.1"

  source = {
    path = "${path.module}/TerraformSimpleTestSolution_1_0_0_1_managed.zip"
  }
}
