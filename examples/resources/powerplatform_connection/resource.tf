terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  use_cli = true
}


resource "powerplatform_connection" "azure_openai_connection" {
  environment_id = var.environment_id
  name           = "shared_azureopenai"
  display_name   = "OpenAI Connection"
  connection_parameters = jsonencode({
    "azureOpenAIResourceName" : "${var.azure_openai_resource_name}",
    "azureOpenAIApiKey" : "${var.azure_openai_api_key}"
    "azureSearchEndpointUrl" : "${var.azure_search_endpoint_url}",
    "azureSearchApiKey" : "${var.azure_search_api_key}"
  })

  lifecycle {
    ignore_changes = [
      connection_parameters_set
    ]
  }

}
