terraform {
  required_providers {
    powerplatform = {
      source = "local/power-platform"
    }
  }
}

provider "powerplatform" {
  use_cli = true
}

resource "powerplatform_copilot_studio_application_insights" "cps_app_insights_config" {
  environment_id                         = var.environment_id
  bot_id                                 = var.bot_id
  application_insights_connection_string = var.application_insights_connection_string
  include_sensitive_information          = false
  include_activities                     = true
  include_actions                        = true
}
