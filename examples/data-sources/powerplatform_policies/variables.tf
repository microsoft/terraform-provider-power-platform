variable "tenant_id" {
  default     = "<my_tenant_id>"
  description = "The tenant id of the AAD tenant"
  type        = string
}
variable "client_id" {
  default     = "<my_client_id>"
  description = "The client id of the AAD application"
  type        = string
}
variable "secret" {
  default     = "<my_secret>"
  description = "The secret of the AAD application"
  type        = string
}
