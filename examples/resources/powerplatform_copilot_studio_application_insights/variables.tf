variable "environment_id" {
  description = "Unique identifier GUID of the Power Platform environment"
  type        = string
}

variable "bot_id" {
  description = "Unique identifier GUID of the Copilot"
  type        = string
}

variable "application_insights_connection_string" {
  description = "The connection string for the target Application Insights resource"
  type        = string
}