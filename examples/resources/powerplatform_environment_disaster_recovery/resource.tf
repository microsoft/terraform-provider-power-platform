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


resource "powerplatform_billing_policy" "prod" {
  name     = "prodBillingPolicy"
  location = "unitedstates"
  status   = "Enabled"
  billing_instrument = {
    resource_group  = "your-resource-group"
    subscription_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "powerplatform_environment" "prod" {
  display_name      = "Production Environment"
  location          = "unitedstates"
  environment_type  = "Production"
  billing_policy_id = powerplatform_billing_policy.prod.id
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "powerplatform_managed_environment" "prod" {
  environment_id             = powerplatform_environment.prod.id
  is_usage_insights_disabled = true
  is_group_sharing_disabled  = false
  limit_sharing_mode         = "NoLimit"
  max_limit_user_sharing     = -1
  solution_checker_mode      = "None"
  suppress_validation_emails = true
}

resource "powerplatform_environment_disaster_recovery" "prod" {
  environment_id = powerplatform_environment.prod.id
  enabled        = true

  depends_on = [powerplatform_managed_environment.prod]
}
