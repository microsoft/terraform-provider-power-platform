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

variable "secret_pp" {
  description = "The secret of the service principal with Power Platform admin permissions"
  sensitive   = true
  type        = string
}

variable "userIdAdmin_pp" {
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

variable "opdgw_install_link" {
  description = "The Blob link to the GatewayInstall.exe file"
  type        = string
}

variable "opdgw_setup_link" {
  description = "The Blob link to the opdgw-setup.ps1 script"
  type        = string
}

variable "sapnco_install_link" {
  description = "The Blob link to the SAP NCo installation file"
  type        = string
}
