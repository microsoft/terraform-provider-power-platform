
##terraform variables

#service principal object id of the terraform service principal
#this service principal has to exist prio to running the terraform script has following pre-requisites:
#https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/service_principal_client_secret
#https://registry.terraform.io/providers/hashicorp/azuread/latest/docs/guides/service_principal_configuration
variable "aad_client_id" {
    default = "11111111-2222-3333-4444-555555555555"
}

#service principal secret of the terraform service principal
variable "aad_client_secret" {
    default = ""
}

#terraform user that will have access to the Power Platform
#this user has to exist prio to running the terraform script has following pre-requisites:
#1. user has to have a valid license for the Power Platform
#2. user has to be Dyamics 365 admin or Power Platform admin
variable "username" {
    default = "user@domain.onmicrosoft.com"
}

#terraform user password
variable "password" {
    default = ""
}

##azure environment variables

variable "resource_group_name" {
  default = "rg_automation_kit"
}

variable "resource_group_location" {
  default = "westeurope"
}

variable "key_vault_name" {
  default = "kv-automation-kit"
}

#tenant id where the azure subscription is located
variable "aad_tenant_id" {
    default = "11111111-2222-3333-4444-555555555555"
}

#subscription id where thre azure resources will be created
variable "azure_subscription_id" {
    default = "11111111-2222-3333-4444-555555555555"
}

##main environment variables

variable "env_variable_autocoe_default_frequency_values" {
  default = "{}"
}

variable "creator_kit_solution_zip_path" {
  default = "/power_platform/data/CreatorKitCore_1.0.20230118.1_managed.zip"
}

variable "automation_coe_main_solution_zip_path" {
  default = "/power_platform/data/AutomationCoEMain_1.0.20230308.1_managed.zip"
}

variable "main_conn_ref_shared_powerplatformforadmins" {
  default = "replace_with_your_connection_reference"
}

variable "main_conn_ref_shared_office365users" {
  default = "replace_with_your_connection_reference"
}

variable "main_conn_ref_shared_office365" {
  default = "replace_with_your_connection_reference"
}

variable "main_conn_ref_shared_commondataserviceforapps" { 
  default = "replace_with_your_connection_reference"
}

variable "main_conn_ref_shared_approvals" {
  default = "replace_with_your_connection_reference"
}


#satelite environment variables

variable "automation_coe_satellite_solution_zip_path" {
  default = "/power_platform/data/AutomationCoESatellite_1.0.20230308.2_managed.zip"
}

variable "env_variable_autocoe_AutomationCoEAlertEmailRecipient" {
  default = "null"
}

variable "env_variable_autocoe_StoreExtractedScript" {
  default = "yes"
}

variable "env_variable_autocoe_FlowSessionTraceRecordOwnerId" {
  default = "null"
}

variable satelite_conn_ref_shared_commondataserviceforapps{
  default = "replace_with_your_connection_reference"
}

variable "satelite_conn_ref_shared_commondataservice" {
  default = "replace_with_your_connection_reference"
}

variable "satelite_conn_ref_shared_flowmanagement" {
  default = "replace_with_your_connection_reference"
}

variable "satelite_conn_ref_shared_office365users" {
  default = "replace_with_your_connection_reference"
}

variable "satelite_conn_ref_shared_powerplatformforadmins" {
  default = "replace_with_your_connection_reference"
}

variable "satelite_conn_ref_shared_office365" {
  default = "replace_with_your_connection_reference"
}