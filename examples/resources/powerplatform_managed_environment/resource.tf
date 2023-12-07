terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  username  = var.username
  password  = var.password
  tenant_id = var.tenant_id
}

resource "powerplatform_environment" "development" {
  display_name      = "example_managed_environment"
  location          = "europe"
  language_code     = "1033"
  currency_code     = "USD"
  environment_type  = "Sandbox"
  domain            = "mydomainmanagedenvironment"
  security_group_id = "00000000-0000-0000-0000-000000000000"
}

resource "powerplatform_managed_environment" "managed_development" {
  environment_id             = powerplatform_environment.development.id
  is_usage_insights_disabled = true
  is_group_sharing_disabled  = true
  limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
  max_limit_user_sharing     = 10
  solution_checker_mode      = "None"
  suppress_validation_emails = true
  maker_onboarding_markdown  = "this is example markdown"
  maker_onboarding_url       = "https://www.microsoft.com"
}

