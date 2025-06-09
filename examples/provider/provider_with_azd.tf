terraform {
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

# Provider configuration with Azure Developer CLI (azd) enabled
# azd simplifies the setup and management of Azure resources, including Power Platform
# This configuration assumes you have already set up azd in your environment
provider "powerplatform" {
  use_dev_cli = true
}
