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
  display_name     = "example_environment_settings"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "powerplatform_managed_environment" "managed_environment" {
  environment_id                  = powerplatform_environment.example_environment_settings.id
  is_usage_insights_disabled      = true
  is_group_sharing_disabled       = true
  limit_sharing_mode              = "ExcludeSharingToSecurityGroups"
  max_limit_user_sharing          = 10
  solution_checker_mode           = "Warn"
  suppress_validation_emails      = true
  solution_checker_rule_overrides = toset(["meta-remove-dup-reg"])
}

resource "powerplatform_environment_settings" "settings" {
  depends_on = [powerplatform_managed_environment.managed_environment]

  environment_id = powerplatform_environment.example_environment_settings.id

  audit_and_logs = {
    plugin_trace_log_setting = "Exception"
    audit_settings = {
      is_audit_enabled             = true
      is_user_access_audit_enabled = true
      is_read_audit_enabled        = true
      log_retention_period_in_days = -1 //Forever
    }
  }
  email = {
    email_settings = {
      max_upload_file_size_in_bytes = 123456
    }
  }
  product = {
    behavior_settings = {
      show_dashboard_cards_in_expanded_state = true
    }
    features = {
      power_apps_component_framework_for_canvas_apps                = false
      enable_access_to_session_transcripts_for_copilot_studio       = false
      enable_transcript_recording_for_copilot_studio                = false
      enable_ai_prompts                                             = false
      enable_copilot_studio_share_data_with_viva_insights           = false
      enable_copilot_studio_cross_geo_share_data_with_viva_insights = false
      enable_preview_and_experimental_ai_models                     = false
      enable_powerapps_maker_bot                                    = false

      enable_ai_powered_chat                        = "Off"
      ai_form_fill_automatic_suggestions            = "Off"
      ai_form_fill_smart_paste_and_file_suggestions = "Off"
      ai_form_fill_toolbar                          = "Off"
      natural_language_grid_and_view_search         = "NoOne"
      allow_ai_to_generate_charts                   = "Off"
    }
    security = {
      allow_application_user_access               = true
      allow_microsoft_trusted_service_tags        = true
      allowed_ip_range_for_firewall               = toset(["10.10.0.0/16", "192.168.0.0/24"])
      allowed_service_tags_for_firewall           = toset(["ApiManagement", "AppService"])
      enable_ip_based_firewall_rule               = true
      enable_ip_based_firewall_rule_in_audit_mode = true
      reverse_proxy_ip_addresses                  = toset(["10.10.1.1", "192.168.1.1"])
    }
  }
}
