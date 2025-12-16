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

resource "powerplatform_environment" "development" {
  display_name     = "example_managed_environment"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "powerplatform_managed_environment" "managed_development" {
  environment_id                                     = powerplatform_environment.development.id
  is_usage_insights_disabled                         = true
  is_group_sharing_disabled                          = true
  limit_sharing_mode                                 = "ExcludeSharingToSecurityGroups"
  max_limit_user_sharing                             = 10
  solution_checker_mode                              = "Warn"
  suppress_validation_emails                         = true
  solution_checker_rule_overrides                    = toset(["meta-avoid-reg-no-attribute", "meta-avoid-reg-retrieve", "app-use-delayoutput-text-input"])
  power_automate_is_sharing_disabled                 = true
  copilot_allow_grant_editor_permissions_when_shared = false
  copilot_limit_sharing_mode                         = "ExcludeSharingToSecurityGroups"
  copilot_max_limit_user_sharing                     = 55
}

