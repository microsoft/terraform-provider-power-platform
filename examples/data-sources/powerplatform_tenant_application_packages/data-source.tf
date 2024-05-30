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

data "powerplatform_tenant_application_packages" "all_applications" {
}

data "powerplatform_tenant_application_packages" "all_applications_from_publisher" {
  publisher_name = "Microsoft Dynamics SMB"
}

data "powerplatform_tenant_application_packages" "specific_application" {
  name = "Healthcare Home Health"
}
