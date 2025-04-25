terraform {
  required_version = "> 1.7.0"
  required_providers {
    powerplatform = {
      source  = "microsoft/power-platform"
      version = "~>3.5.0"
    }
    azapi = {
      source  = "azure/azapi"
      version = "~>2.2.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~>4.16.0"
    }
  }
}
provider "powerplatform" {
  use_cli = true
}

provider "azurerm" {
  use_cli         = true
  subscription_id = var.subscription_id
  features {
    resource_group {
      prevent_deletion_if_contains_resources = false
    }
  }
}

// getting all locations available for the environment and their azure regions
data "powerplatform_locations" "all_locations" {
}

// getting european loocation details. Policy has to be in the same azure region as the environment
locals {
  europe_location = [for location in data.powerplatform_locations.all_locations.locations : location if location.name == "europe"]
}

// creating environment that will have the policies applied
resource "powerplatform_environment" "example_environment" {
  display_name     = "example_environment"
  location         = local.europe_location[0].name
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

// creating managed environment for the environment as this is required for encryption policy
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


// module that creates all azure resources required for the network injection (resource group, vnet's and subnet's) policy and the policy itself
module "network_injection" {
  source = "./network_injection"

  should_register_provider = false

  resource_group_name        = "rg_example_network_injection_policy"
  resource_group_location    = local.europe_location[0].azure_regions[0]
  vnet_locations             = local.europe_location[0].azure_regions
  enterprise_policy_name     = "ep_example_network_injection_policy"
  enterprise_policy_location = local.europe_location[0].name
}

// module that creates all azure resources required for the encryption policy and the policy itself
module "encryption" {
  source = "./encryption"

  should_register_provider = false

  environment_id = powerplatform_environment.example_environment.id

  resource_group_name        = "rg_example_encryption_policy8"
  resource_group_location    = local.europe_location[0].azure_regions[0]
  enterprise_policy_name     = "ep_example_encryption_policy8"
  enterprise_policy_location = "europe"
  keyvault_name              = "kv-ep-example8"

  // let's wait for first policy to be executed
  depends_on = [powerplatform_enterprise_policy.network_injection]
}

resource "powerplatform_enterprise_policy" "network_injection" {
  environment_id = powerplatform_environment.example_environment.id
  system_id      = module.network_injection.policy_system_id
  policy_type    = "NetworkInjection"

  depends_on = [powerplatform_managed_environment.managed_development]
}
