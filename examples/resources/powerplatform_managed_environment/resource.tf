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
  display_name     = "example_managed_env_sol"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
    domain            = "examplemanagedenvsol-sandbox"
  }
}

resource "powerplatform_managed_environment" "managed_development" {
  environment_id                  = powerplatform_environment.development.id
  is_usage_insights_disabled      = true
  is_group_sharing_disabled       = true
  limit_sharing_mode              = "ExcludeSharingToSecurityGroups"
  max_limit_user_sharing          = 10
  solution_checker_mode           = "Warn"
  suppress_validation_emails      = true
  solution_checker_rule_overrides = ["meta-remove-dup-reg", "meta-avoid-reg-no-attribute"]
  maker_onboarding_markdown       = "this is example markdown"
  maker_onboarding_url            = "https://www.microsoft.com"
}

