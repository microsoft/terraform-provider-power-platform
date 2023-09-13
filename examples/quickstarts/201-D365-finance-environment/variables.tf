variable "client_id" {
  description = "The client ID of the of the service principal"
  type        = string
}
variable "secret" {
  description = "The client secret of the service principal"
  sensitive   = true
  type        = string
  # validation {
  #   condition     = (var.username != null && var.password != null) || (var.username == null && var.password == null)
  #   error_message = "Both username and password must be provided, or both should be null."
  # }
}
variable "tenant_id" {
  description = "The AAD tenant id of service principal or user"
  type        = string
}
variable "d365_finance_environment_name" {
  description = "The name of the D365 Finance environment"
  type        = string
  validation {
    condition     = (len(var.d365_finance_environment_name) < 20)
    error_message = "The length of the d365_finance_environment_name property cannot exceed 20 characters for F&O environment deployments."
  }
}
#Available locations are listed at //https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01
variable "location" {
  description = "The Azure region where the environment will be deployed"
  type        = string
  #This default will eventually be removed when other regions become supported.
  default = "Canada"
}
#Available langauge codes are listed at https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentLanguages?api-version=2023-06-01
variable "language_code" {
  description = "The desired Language Code for the environment"
  type        = string
  default     = "1033"
}
#Available currency codes are listed at https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations/unitedstates/environmentCurrencies?api-version=2023-06-01
variable "currency_code" {
  description = "The desired Currency Code for the environment"
  type        = string
  default     = "USD"
}
#Options are "Sandbox", "Production", "Trial", "Developer"
variable "environment_type" {
  description = "The type of environment to deploy"
  type        = string
  default     = "Sandbox"
}
variable "security_group_id" {
  description = "The security group the environment will be associated with"
  type        = string
  default     = "00000000-0000-0000-0000-000000000000"
}
variable "domain" {
  description = "The domain of the environment"
  type        = string
  default     = "sample-d365-finance-environment"
}
