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

data "powerplatform_data_records" "example_data_records" {
  environment_id = "838f76c8-a192-e59c-a835-089ad8cfb047"
  # prefer = tomap({
  #   //respond-async             = null
  #   odata.maxpagesize         = "1"
  #   odata.include-annotations = "OData.Community.Display.V1.FormattedValue,Microsoft.PowerApps.CDS.ErrorDetails*"
  # })
  query = "systemusers?$select=fullname,systemuserid,createdon$expand=systemuserroles_association($select=roleid,roleidname)$top=2"
}
