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

//https://orgda1371a4.crm17.dynamics.com/api/data/v9.2/savedqueries?$select=name,savedqueryid,returnedtypecode&$filter=returnedtypecode%20eq%20%27systemuser%27%20and%20name%20eq%20%27Enabled%20Users%27
# data "powerplatform_data_records" "saved_view" {
#   environment_id    = "838f76c8-a192-e59c-a835-089ad8cfb047"
#   entity_collection = "userqueries"
#   select            = ["name", "returnedtypecode"]
#   //filter            = "returnedtypecode eq 'systemuser' and name eq 'Enabled Users'"
#   //top               = 1
# }

//https://orgda1371a4.crm17.dynamics.com/api/data/v9.2/systemusers?$select=fullname&$expand=systemuserroles_association($select=name),teammembership_association($select=name)
data "powerplatform_data_records" "example_data_records" {
  environment_id    = "838f76c8-a192-e59c-a835-089ad8cfb047"
  entity_collection = "systemusers"
  select            = ["fullname", "systemuserid"]
  expand = [
    {
      navigation_property = "systemuserroles_association"
      select              = ["roleid", "name"]
      // filter, select, orderby, top, expand
      //filter = "roleid eq 'f5e3d3d1-1ee7-4f3d-8f1d-4b4b1f7e3e3d'"
    },

  ]

}
