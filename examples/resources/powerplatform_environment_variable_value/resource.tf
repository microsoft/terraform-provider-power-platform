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

# This resource manages the current value for an existing definition.
# The definition itself must already exist in the environment, typically
# because it was deployed by a solution.
resource "powerplatform_environment_variable_value" "api_base_url" {
  environment_id = "00000000-0000-0000-0000-000000000000"
  schema_name    = "contoso_ApiBaseUrl"
  value          = "https://api.contoso.example"
}
