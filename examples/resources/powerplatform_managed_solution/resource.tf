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

  # Set exactly one of path or url. A url is typically a build artifact
  # location such as a blob SAS url and is treated as sensitive.
  source = {
    path = "${path.module}/TerraformSolutionExample_1_0_0_1_managed.zip"
    # url = "https://example.blob.core.windows.net/artifacts/TerraformSolutionExample_1_0_0_1_managed.zip?<sas>"
  }

  # Every connection reference declared by the package must be bound to a
  # connection in the target environment.
  connection_references = {
    terr_SolutionConnectionReference = powerplatform_connection.dataverse_connection.id
  }

  # Optional import behavior. Both default to false.
  skip_product_update_dependencies = false
  publish_all_customizations       = false

  depends_on = [powerplatform_connection.dataverse_connection]
}

# Environment variable definitions ship inside the solution package; their
# values are owned by powerplatform_environment_variable_value and set after
# the import that ships the definitions, ordered with an ordinary depends_on.
# A definition with no value is valid until a consumer needs one, so unset
# variables (like terr_SolutionVariableDataSource here) need no placeholder.
resource "powerplatform_environment_variable_value" "solution_variable_text" {
  environment_id = powerplatform_environment.environment.id
  schema_name    = "terr_SolutionVariableText"
  value          = "sample text value"

  depends_on = [powerplatform_managed_solution.solution]
}

resource "powerplatform_environment_variable_value" "solution_variable_json" {
  environment_id = powerplatform_environment.environment.id
  schema_name    = "terr_SolutionVariableJson"
  value = jsonencode({
    value = 1234
    text  = "abc"
  })

  depends_on = [powerplatform_managed_solution.solution]
}
