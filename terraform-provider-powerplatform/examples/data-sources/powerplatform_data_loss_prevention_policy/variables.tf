variable "username" {
  default     = "user@domain.onmicrosoft.com"
  description = "The username of the Power Platform API in user@domain format"
  type        = string
}

variable "password" {
  description = "The password of the Power Platform API user"
  sensitive   = true
  type        = string
}

variable "host" {
  default     = "http://localhost:8080"
  description = "The host URL of the Power Platform API in https://<host> format"
  type        = string
}
