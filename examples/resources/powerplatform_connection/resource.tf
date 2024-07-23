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


resource "powerplatform_connection" "new_connection" {
  environment_id = "00000000-0000-0000-0000-000000000000"
  name           = "shared_azureopenai"
  display_name   = "OpenAI Connection 123"
  connection_parameters = jsonencode({
    "azureOpenAIResourceName" : "aaa",
    "azureOpenAIApiKey" : "bbb",
    "azureSearchEndpointUrl" : "ccc",
    "azureSearchApiKey" : "dddd"
  })
}
