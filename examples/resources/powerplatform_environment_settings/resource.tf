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

resource "powerplatform_environment" "example_environment_settings" {
  display_name      = "example_environment_settings"
  location          = "europe"
  language_code     = "1033"
  currency_code     = "USD"
  environment_type  = "Sandbox"
  security_group_id = "00000000-0000-0000-0000-000000000000"
}

resource "powerplatform_environment_settings" "settings" {
  environment_id                         = powerplatform_environment.example_environment_settings.id
  max_upload_file_size_in_bytes          = 100
  show_dashboard_cards_in_expanded_state = true
  plugin_trace_log_setting               = "Off"
  is_audit_enabled                       = true
  is_user_access_audit_enabled           = true
  is_read_audit_enabled                  = true
}
