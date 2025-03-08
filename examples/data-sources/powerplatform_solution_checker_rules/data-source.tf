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

# Create an environment to use for solution checker rules
resource "powerplatform_environment" "example" {
  display_name     = "Solution Checker Example"
  location         = "unitedstates"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

# Output to ensure Dataverse URL exists (indicates Dataverse is ready)
output "dataverse_url" {
  value = powerplatform_environment.example.dataverse.url
}

# Use the created environment's ID for solution checker rules
data "powerplatform_solution_checker_rules" "example" {
  # Only proceed after both environment ID and Dataverse URL are available
  depends_on = [powerplatform_environment.example]
  environment_id = powerplatform_environment.example.id
}
