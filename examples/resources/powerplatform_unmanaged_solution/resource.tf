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
  display_name     = "Unmanaged Solution Example"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "powerplatform_unmanaged_solution" "solution" {
  environment_id = powerplatform_environment.environment.id
  uniquename     = "TerraformUnmanagedSolution"
  display_name   = "Terraform Unmanaged Solution"
  publisher_id   = "00000000-0000-0000-0000-000000000000"
  description    = "Unmanaged solution created directly through the Dataverse solutions table."
}
