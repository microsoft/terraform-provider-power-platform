terraform {
  required_version = ">= 1.5"
  required_providers {
    powerplatform = {
      source = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  client_id = var.client_id
  secret    = var.secret
  tenant_id = var.tenant_id
}


resource "powerplatform_environment" "development" {
  display_name     = var.d365_finance_environment_name
  location         = var.location
  language_code    = var.language_code
  currency_code    = var.currency_code
  environment_type = var.environment_type
  domain           = var.domain
  //There are many template options, including for other business applications.
  templates = ["D365_FinOps_Finance"]
  //This is a special JSON-formatted parameter specification that is currently required for D365 Finance deployments.
  template_metadata = "{\"PostProvisioningPackages\": [{ \"applicationUniqueName\": \"msdyn_FinanceAndOperationsProvisioningAppAnchor\",\n \"parameters\": \"DevToolsEnabled=true|DemoDataEnabled=true\"\n }\n ]\n }"
  security_group_id = var.security_group_id
}
