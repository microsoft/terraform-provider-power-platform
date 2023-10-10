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

variable "nic_id" {
  description = "The id of the network interface to attach to the VM"
  type        = string
}

variable "installps7_link" {
  description = "The Blob link to the PowerShell 7 installation file"
  type        = string
}

