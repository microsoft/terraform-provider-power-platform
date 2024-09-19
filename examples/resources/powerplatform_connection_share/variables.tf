variable "environment_id" {
  default     = "00000000-0000-0000-0000-000000000001"
  description = "Unique identifier of the environment"
  type        = string
}

variable "tenant_id" {
  default     = "00000000-0000-0000-0000-000000000002"
  description = "Unique identifier of the tenant where the environment is located"
  type        = string
}

variable "client_id" {
  default     = ""
  description = "Unique identifier of the client application that will be used for authentication"
  type        = string
}

variable "client_secret" {
  default     = ""
  description = "Secret of the client application that will be used for authentication"
  sensitive   = true
  type        = string
}

variable "user_object_id" {
  default     = "00000000-0000-0000-0000-000000000003"
  description = "Entra Object Id identifier of the user that will be granted access to the connection"
  type        = string
}
