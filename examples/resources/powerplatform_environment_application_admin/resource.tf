terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

resource "powerplatform_environment_application_admin" "import_fix" {
  environment_id = var.environment_id        # GUID of environment
  application_id = var.spn_application_id    # GUID (client ID) of the SP
}