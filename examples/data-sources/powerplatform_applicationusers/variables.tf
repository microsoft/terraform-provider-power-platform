variable "tenant_id" {
  description = "The tenant id of the AAD tenant"
  type        = string
}

variable "environment_id" {
  description = "The name of the environment"
  type        = string
}
variable "client_id" {
  description = "The client ID of the of the service principal"
  type        = string

}
variable "secret" {
  description = "The client secret of the service principal"
  sensitive   = true
  type        = string
}
