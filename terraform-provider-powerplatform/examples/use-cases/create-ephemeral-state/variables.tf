variable "username" {
  default     = "user@domain.onmicrosoft.com"
  description = "The username of the Power Platform API in user@domain format"
  type        = string
}

variable "password" {
  description = "The password of the Power Platform API user"
  sensitive   = true
  type        = string
}

variable "host" {
  default     = "http://localhost:8080"
  description = "The host URL of the Power Platform API in https://<host> format"
  type        = string
}

variable "aad_client_id" {
  default     = "11111111-2222-3333-4444-5555555555555"
  description = "Client ID for the Azure AD application used to authenticate to the Azure API"
}

variable "aad_client_secret" {
  description = "Client secret for the Azure AD application used to authenticate to the Azure API"
  sensitive   = true
}

variable "aad_tenant_id" {
  default     = "11111111-2222-3333-4444-555555555555"
  description = "Tenant ID for the Azure AD application used to authenticate to the Azure API"
}

variable "aad_licensing_security_group" {
  description = "AAD Group ID for the security group that has the Power Platform licenses assigned"
  value       = "11111111-2222-3333-4444-555555555555"
}

variable "new_user_password" {
  description = "The password of the new AAD user"
  sensitive   = true
}