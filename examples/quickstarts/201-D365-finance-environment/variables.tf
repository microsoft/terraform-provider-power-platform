variable "client_id" {
  description = "The client ID of the service principal (app registration)"
  type        = string
}
variable "secret" {
  description = "The client secret of the service principal (app registration)"
  sensitive   = true
  type        = string
}
variable "tenant_id" {
  description = "The Entra (AAD) tenant id of service principal or user"
  type        = string
}
variable "d365_finance_environment_name" {
  description = "The name of the D365 Finance environment, such as 'd365fin-dev1"
  type        = string
  validation {
    condition     = length(var.d365_finance_environment_name) <= 20
    error_message = "The length of the d365_finance_environment_name property cannot exceed 20 characters for D365 Finance environment deployments."
  }
}
#Available locations are listed at //https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01
variable "location" {
  description = "The region where the environment will be deployed, such as 'unitedstates'"
  type        = string
  #This default will eventually be removed when other regions become supported.
}
#Available langauge codes are listed at https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentLanguages?api-version=2023-06-01
variable "language_code" {
  description = "The desired Language Code for the environment, such as '1033' (U.S. english)"
  type        = string
}
#Available currency codes are listed at https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentCurrencies?api-version=2023-06-01
variable "currency_code" {
  description = "The desired Currency Code for the environment, such as 'USD'"
  type        = string
}
#Options are "Sandbox", "Production", "Trial", "Developer"
variable "environment_type" {
  description = "The type of environment to deploy, such as 'Sandbox'"
  type        = string
  validation {
    condition     = contains(["Sandbox", "Production", "Trial", "Developer"], var.environment_type)
    error_message = "The selected value for environment_type is not in the list of allowed values in variables.tf"
  }
}
variable "security_group_id" {
  description = "The security group the environment will be associated with, a GUID. Can be set to 00000000-0000-0000-0000-000000000000 to indicate no security group restricting Dataverse access."
  type        = string
}
variable "domain" {
  description = "The domain of the environment, such as 'd365fin-dev1'"
  type        = string
  validation {
    condition     = length(var.domain) <= 32
    error_message = "The length of the domain property cannot exceed 32 characters for D365 Finance environment deployments."
  }
}
