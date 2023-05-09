terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/powerplatform"
    }
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 2.15.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.0.0"
    }
  }
}

provider "powerplatform" {
  username = var.username
  password = var.password
  host     = "http://localhost:8080"
}

provider "azuread" {
  tenant_id     = var.aad_tenant_id
  client_id     = var.aad_client_id
  client_secret = var.aad_client_secret
}

provider "azurerm" {

  subscription_id = var.azure_subscription_id
  client_id       = var.aad_client_id
  client_secret   = var.aad_client_secret
  tenant_id       = var.aad_tenant_id

  features {
    key_vault {
      purge_soft_delete_on_destroy    = false
      recover_soft_deleted_key_vaults = true
    }
  }
}

data "azuread_domains" "aad_domains" {
  only_initial = true
}

data "azurerm_client_config" "current" {}


module "azure" {
  terraform_service_principal_object_id = data.azurerm_client_config.current.object_id

  source                  = "./azure"
  tenant_id               = var.aad_tenant_id
  domain_name             = data.azuread_domains.aad_domains.domains[0].domain_name
  resource_group_name     = var.resource_group_name
  resource_group_location = var.resource_group_location
  key_vault_name          = var.key_vault_name
}

module "power_platform" {
  source = "./power_platform"

  automation_kit_application_id = module.azure.automation_kit_application_id

  key_vault_name                         = module.azure.key_vault_name
  key_vault_secret_client_id_name        = module.azure.key_vault_secret_client_id_name
  key_vault_secret_client_password_name  = module.azure.key_vault_secret_client_password_name
  key_vault_client_secret_tenant_id_name = module.azure.key_vault_client_secret_tenant_id_name

  creator_kit_solution_zip_path              = "${path.module}${var.creator_kit_solution_zip_path}"
  automation_coe_main_solution_zip_path      = "${path.module}${var.automation_coe_main_solution_zip_path}"
  automation_coe_satellite_solution_zip_path = "${path.module}${var.automation_coe_satellite_solution_zip_path}"

  main_conn_ref_shared_powerplatformforadmins   = var.main_conn_ref_shared_powerplatformforadmins
  main_conn_ref_shared_office365users           = var.main_conn_ref_shared_office365users
  main_conn_ref_shared_office365                = var.main_conn_ref_shared_office365
  main_conn_ref_shared_commondataserviceforapps = var.main_conn_ref_shared_commondataserviceforapps
  main_conn_ref_shared_approvals                = var.main_conn_ref_shared_approvals
  env_variable_autocoe_default_frequency_values = var.env_variable_autocoe_default_frequency_values

  satelite_conn_ref_shared_commondataserviceforapps     = var.satelite_conn_ref_shared_commondataserviceforapps
  satelite_conn_ref_shared_commondataservice            = var.satelite_conn_ref_shared_commondataservice
  satelite_conn_ref_shared_flowmanagement               = var.satelite_conn_ref_shared_flowmanagement
  satelite_conn_ref_shared_office365users               = var.satelite_conn_ref_shared_office365users
  satelite_conn_ref_shared_powerplatformforadmins       = var.satelite_conn_ref_shared_powerplatformforadmins
  satelite_conn_ref_shared_office365                    = var.satelite_conn_ref_shared_office365
  env_variable_autocoe_AutomationCoEAlertEmailRecipient = var.env_variable_autocoe_AutomationCoEAlertEmailRecipient
  env_variable_autocoe_StoreExtractedScript             = var.env_variable_autocoe_StoreExtractedScript
  env_variable_autocoe_FlowSessionTraceRecordOwnerId    = var.env_variable_autocoe_FlowSessionTraceRecordOwnerId

}


