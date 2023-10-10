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
  description = "The client id / app id of the service principal with Power Platform admin permissions"
  type        = string
}

variable "secret_pp" {
  description = "The secret of the service principal with Power Platform admin permissions"
  sensitive   = true
  type        = string
}
variable "tenant_id_pp" {
  description = "The AAD tenant id of service principal or user at Power Platform"
  type        = string
}

variable "client_id_gw" {
  description = "The client id / app id of the service principal where the on-premise data gateway admin permissions"
  type        = string
}

variable "secret_gw" {
  description = "The secret of the service principal with on-premise data gateway admin permissions"
  sensitive   = true
  type        = string
}
variable "tenant_id_gw" {
  description = "The AAD tenant id of service principal or user"
  type        = string
}

variable "subscription_id_gw" {
  description = "The subscription id of the service principal with on-premise data gateway admin permissions"
  type        = string
}

variable "region_gw" {
  description = "The Azure region where the resources in this example should be created"
  type        = string
}

variable "vm_pwd_gw" {
  description = "The password for the VM of the on-premise data gateway"
  sensitive   = true
  type        = string
}

variable "installps7_link" {
  description = "The Blob link to the PowerShell 7 installation file"
  type        = string
}
