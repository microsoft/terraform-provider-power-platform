variable "environment_id" {
  default     = "00000000-0000-0000-0000-000000000001"
  description = "Unique identifier of the environment"
  type        = string
}

variable "azure_openai_connection_id" {
  default     = "00000000-0000-0000-0000-000000000002"
  description = "Unique identifier of the connection"
  type        = string
}
