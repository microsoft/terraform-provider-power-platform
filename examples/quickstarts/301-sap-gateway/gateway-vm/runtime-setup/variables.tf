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

variable "sig_id" {
  description = "The id of the shared image gallery where the image should be created"
  type        = string
}

variable "runtime_setup_link" {
  description = "The Blob link to the runtime-setup.ps1 script"
  type        = string
}
