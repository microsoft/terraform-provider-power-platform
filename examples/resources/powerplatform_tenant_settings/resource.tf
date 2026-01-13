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

resource "powerplatform_tenant_settings" "settings" {
  walk_me_opt_out                                       = true
  disable_support_tickets_visible_by_all_users          = true
  disable_trial_environment_creation_by_non_admin_users = true
  disable_capacity_allocation_by_environment_admins     = true
  disable_environment_creation_by_non_admin_users       = true
  disable_portals_creation_by_non_admin_users           = true
  disable_newsletter_sendout                            = true
  enable_support_use_bing_search_solutions              = true

  power_platform = {
    product_feedback = {
      disable_microsoft_surveys_send = true
      disable_user_survey_feedback   = true
      disable_attachments            = true
      disable_microsoft_follow_up    = true
    }
    search = {
      disable_docs_search       = true
      disable_community_search  = true
      disable_bing_video_search = true
    }
    teams_integration = {
      share_with_colleagues_user_limit = 10001
    }
    power_apps = {
      disable_copilot                          = true
      disable_share_with_everyone              = true
      enable_guests_to_make                    = true
      disable_members_indicator                = true
      disable_maker_match                      = true
      disable_unused_license_assignment        = true
      disable_connection_sharing_with_everyone = true
      enable_canvas_app_insights               = true
    }
    power_automate = {
      disable_copilot                 = true
      disable_copilot_help_assistance = true
      allow_use_of_hosted_browser     = true
      disable_flow_resubmission       = true
    }
    environments = {
      disable_preferred_data_location_for_teams_environment = true
    }
    governance = {
      weekly_digest_email_recipients                            = toset(["test1@contoso.com", "test2@contoso.com"])
      disable_admin_digest                                      = true
      disable_developer_environment_creation_by_non_admin_users = true
      enable_default_environment_routing                        = false
      policy = {
        enable_desktop_flow_data_policy_management = true
      }
    }
    licensing = {
      disable_billing_policy_creation_by_non_admin_users           = true
      enable_tenant_capacity_report_for_environment_admins         = true
      storage_capacity_consumption_warning_threshold               = 88
      enable_tenant_licensing_report_for_environment_admins        = true
      enable_tenant_capacity_report_for_environment_admins         = true
      disable_use_of_unassigned_ai_builder_credits                 = true
      apply_auto_claim_power_apps_to_only_managed_environments     = true
      apply_auto_claim_power_automate_to_only_managed_environments = true
    }
    power_pages = {}
    champions = {
      disable_champions_invitation_reachout    = true
      disable_skills_match_invitation_reachout = true
    }
    intelligence = {
      disable_copilot                                            = true
      allow_copilot_authors_publish_when_ai_features_are_enabled = true
      basic_copilot_feedback                                     = true
      additional_copilot_feedback                                = true
      copilot_studio_authors_security_group_id                   = "00000000-0000-0000-0000-000000000000"
    }
    model_experimentation = {
      enable_model_data_sharing = true
      disable_data_logging      = true
    }
    catalog_settings = {
      power_catalog_audience_setting = "All"
    }
    user_management_settings = {
      enable_delete_disabled_user_in_all_environments = true
    }
  }
}
