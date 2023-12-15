variable "prefix" {
  description = "The prefix which should be used for all resources name"
  default     = "opdgw"
  type        = string
}

variable "base_name" {
  description = "The base name which should be used for all resources name"
  default     = "AzureSAPIntegration"
  type        = string
}

variable "client_id_pp" {
  description = "The client ID / app ID of the service principal with Power Platform admin permissions"
  type        = string
}

variable "secret_pp" {
  description = "The secret of the service principal with Power Platform admin permissions"
  sensitive   = true
  type        = string
}

variable "tenant_id_pp" {
  description = "The tenant ID of service principal or user at Power Platform"
  type        = string
}

variable "client_id_gw" {
  description = "The client ID / app ID of the service principal where the on-premise data gateway admin permissions"
  type        = string
}

variable "secret_gw" {
  description = "The secret of the service principal with on-premise data gateway admin permissions"
  sensitive   = true
  type        = string
}
variable "tenant_id_gw" {
  description = "The tenant ID of service principal or user"
  type        = string
}

variable "subscription_id_gw" {
  description = "The subscription ID of the service principal with on-premise data gateway admin permissions"
  type        = string
}

variable "region_gw" {
  description = "The Azure region where the resources in this example should be created"
  type        = string
}

variable "sap_subnet_id" {
  description = "The SAP system subnet ID"
  type        = string
}

variable "user_id_admin_pp" {
  description = "The user ID to be assigned as Admin role of the Power Platform"
  type        = string
}

variable "ir_key" {
  description = "Value of the secret name for the IR key"
  type        = string
}

variable "gateway_name" {
  description = "The name of the gateway to be created on Power Platform"
  type        = string
}

variable "recover_key_gw" {
  description = "The recovery key of the gateway"
  type        = string
}
