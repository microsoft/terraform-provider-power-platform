variable "client_id" {
  description = "The client id / app id of the service principal with Power Platform admin permissions"
  type        = string
}

variable "secret" {
  description = "The secret of the service principal with Power Platform admin permissions"
  sensitive   = true
  type        = string
}
variable "tenant_id" {
  description = "The AAD tenant id of service principal or user at Power Platform"
  type        = string
}
