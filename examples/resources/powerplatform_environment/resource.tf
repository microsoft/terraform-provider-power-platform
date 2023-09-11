terraform {
  required_providers {
    powerplatform = {
      version = "0.2"
      source  = "microsoft/power-platform"
    }
  }
}

provider "powerplatform" {
  username  = var.username
  password  = var.password
  client_id = var.client_id
  //secret    = var.secret
  tenant_id = var.tenant_id
}

resource "powerplatform_environment" "development" {
  display_name      = "DELETE - command with templates populated"
  location          = "canada"
  language_code     = "1033"
  currency_code     = "USD"
  environment_type  = "Sandbox"
  domain            = "tf-generated-fno-environment"
  templates         = ["D365_Finops_Finance"]
  template_metadata = "{\"PostProvisioningPackages\": [{ \"applicationUniqueName\": \"msdyn_FinanceAndOperationsProvisioningAppAnchor\",\n \"parameters\": \"DevToolsEnabled=true|DemoDataEnabled=true\"\n }\n ]\n }"
  security_group_id = "00000000-0000-0000-0000-000000000000"
}


