variable "prefix" {
  description = "The prefix which should be used for all resources in this example"
  default     = "opdgw"
  type        = string
}

variable "base_name" {
  description = "The base name which should be used for all resources in this example"
  default     = "AzureSAPIntegration"
  type        = string
}

variable "resource_group_name" {
  description = "The name of the resource group where the resources in this example should be created"
  type        = string
}

variable "region" {
  description = "The Azure region where the resources in this example should be created"
  type        = string
}

variable "vm_pwd" {
  description = "The password for the VM"
  sensitive   = true
  type        = string
}

variable "client_id_pp" {
  description = "The client id of the service principal for Power Platform"
  type        = string
}

variable "tenant_id_pp" {
  description = "The tenant id of the service principal for Power Platform"
  type        = string
}

variable "key_vault_uri" {
  description = "The URI of the Key Vault"
  type        = string
}

variable "secret_pp_name" {
  description = "Value of the secret name for Power Platform"
  type        = string
}

variable "secret_name_irkey" {
  description = "Value of the secret name for Integration Runtime Key"
  type        = string
}

variable "user_id_admin_pp" {
  description = "The user id to be assigned as Admin role of the Power Platform"
  type        = string
}

variable "nic_id" {
  description = "The id of the network interface to attach to the VM"
  type        = string
}

variable "ps7_setup_link" {
  description = "The Blob link to the PowerShell 7 installation file"
  type        = string
}

variable "java_setup_link" {
  description = "The Blob link to the Java Runtime installation file"
  type        = string
}

variable "sapnco_install_link" {
  description = "The Blob link to the SAP NCo installation file"
  type        = string
}

variable "runtime_setup_link" {
  description = "The Blob link to the runtime setup script"
  type        = string
}

variable "gateway_name" {
  description = "The name of the gateway"
  type        = string
}

variable "secret_name_recover_key_gw" {
  description = "Value of the secret name for the recovery key of the gateway"
  type        = string
}
