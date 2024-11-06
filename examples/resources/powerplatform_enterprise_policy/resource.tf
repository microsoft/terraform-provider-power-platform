terraform {
  required_version = "> 1.7.0"
  required_providers {
    powerplatform = {
      source  = "microsoft/power-platform"
      version = "~>3.0.0"
    }
    azapi = {
      source  = "azure/azapi"
      version = "2.0.1"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~>4.8.0"
    }
  }
}

provider "powerplatform" {
  alias   = "pp"
  use_cli = true
}

provider "azurerm" {
  alias           = "azrm"
  use_cli         = true
  subscription_id = var.subscription_id
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}


resource "powerplatform_environment" "example_environment" {
  display_name     = "example_environment"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "powerplatform_managed_environment" "managed_development" {
  environment_id             = powerplatform_environment.example_environment.id
  is_usage_insights_disabled = true
  is_group_sharing_disabled  = true
  limit_sharing_mode         = "ExcludeSharingToSecurityGroups"
  max_limit_user_sharing     = 10
  solution_checker_mode      = "None"
  suppress_validation_emails = true
  maker_onboarding_markdown  = "this is example markdown"
  maker_onboarding_url       = "https://www.microsoft.com"
}

module "network_injection" {
  providers = {
    azurerm.azrm = azurerm.azrm
  }

  source = "./network_injection"

  should_register_provider = false

  resource_group_name        = "rg_example_network_injection_policy"
  resource_group_location    = "westeurope"
  vnet_locations             = ["westeurope", "northeurope"]
  enterprise_policy_name     = "ep_example_network_injection_policy"
  enterprise_policy_location = "europe"
}

resource "powerplatform_enterprise_policy" "network_injection" {
  environment_id = powerplatform_environment.example_environment.id
  system_id      = module.network_injection.enterprise_policy_system_id
  policy_type    = "NetworkInjection"

  depends_on = [powerplatform_managed_environment.managed_development]
}

module "encryption" {
  providers = {
    powerplatform.pp = powerplatform.pp
    azurerm.azrm     = azurerm.azrm
  }

  source = "./encryption"

  should_register_provider = false

  resource_group_name        = "rg_example_encryption_policy8"
  resource_group_location    = "westeurope"
  enterprise_policy_name     = "ep_example_encryption_policy8"
  enterprise_policy_location = "europe"
  keyvault_name              = "kv-ep-example8"
}

resource "powerplatform_enterprise_policy" "encryption" {
  environment_id = powerplatform_environment.example_environment.id
  system_id      = module.encryption.enterprise_policy_system_id
  policy_type    = "Encryption"

  depends_on = [powerplatform_enterprise_policy.network_injection]
}
