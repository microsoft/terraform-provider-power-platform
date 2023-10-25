terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
}

# resource "powerplatform_environment" "development" {
#   display_name     = "example_managed_environment"
#   location         = "europe"
#   language_code    = "1033"
#   currency_code    = "USD"
#   environment_type = "Sandbox"
#   domain           = "mydomain_managed_environment"
# }

resource "powerplatform_managed_environment" "managed_development" {
  environment_id             = "ae6407ff-270b-e60c-b82b-e4e7e4b36154" //powerplatform_environment.development.id
  is_usage_insights_disabled = true
  is_group_sharing_disabled  = true
  limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
  max_limit_user_sharing     = 141
  solution_checker_mode      = "None"
  suppress_validation_emails = true
  maker_onboarding_markdown  = "adadsasdasd111"
  maker_onboarding_url       = "aaa111"
}


