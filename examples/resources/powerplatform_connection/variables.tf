variable "environment_id" {
  default     = "00000000-0000-0000-0000-000000000001"
  description = "Unique identifier of the environment"
  type        = string
}

variable "azure_openai_resource_name" {
  default     = "azureopenai"
  description = "Name of the Azure OpenAI resource"
  type        = string
}

variable "azure_openai_api_key" {
  default     = ""
  description = "API key of the Azure OpenAI resource"
  sensitive   = true
  type        = string
}

variable "azure_search_endpoint_url" {
  default     = "https://azuresearchendpoint.com"
  description = "URL of the Azure Search endpoint"
  type        = string
}

variable "azure_search_api_key" {
  default     = ""
  description = "API key of the Azure Search endpoint"
  sensitive   = true
  type        = string
}
