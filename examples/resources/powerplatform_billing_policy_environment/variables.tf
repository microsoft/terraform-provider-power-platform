variable "username" {
  default     = "user@domain.onmicrosoft.com"
  description = "The username of the Power Platform API in user@domain format"
  type        = string

}
variable "password" {
  default     = "<my_passoword>"
  description = "The password of the Power Platform API user"
  sensitive   = true
  type        = string
}
variable "tenant_id" {
  default     = "<my_tenant_id>"
  description = "The tenant id of the AAD tenant"
  type        = string
}
