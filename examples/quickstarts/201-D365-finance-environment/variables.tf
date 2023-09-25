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
  description = "The name of the D365 Finance environment"
  type        = string
  validation {
    condition     = length(var.d365_finance_environment_name) <= 20
    error_message = "The length of the d365_finance_environment_name property cannot exceed 20 characters for D365 Finance environment deployments."
  }
  default = "d365fin-environment"
}
#Available locations are listed at //https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/locations?api-version=2023-06-01
variable "location" {
  description = "The Azure region where the environment will be deployed"
  type        = string
  #This default will eventually be removed when other regions become supported.
  default = "canada"
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
  default     = "sample-d365f-environment"
  validation {
    condition     = length(var.domain) <= 32
    error_message = "The length of the domain property cannot exceed 32 characters for D365 Finance environment deployments."
  }
}
variable "templates" {
  description = "The list of application templates to use when deploying the environment."
  type        = list(string)
  default     = ["D365_FinOps_Finance"]
}
variable "template_metadata" {
  description = "Any additional JSON-formatted metadata required to augment the selected templates."
  type        = string
  default     = "{\"PostProvisioningPackages\": [{ \"applicationUniqueName\": \"msdyn_FinanceAndOperationsProvisioningAppAnchor\",\n \"parameters\": \"DevToolsEnabled=true|DemoDataEnabled=true\"\n }\n ]\n }"
}
