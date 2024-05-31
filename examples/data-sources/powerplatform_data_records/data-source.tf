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
  environment_id    = "838f76c8-a192-e59c-a835-089ad8cfb047"
  entity_collection = "systemusers"
  //entity_collection = "systemusers(1f70a364-5019-ef11-840b-002248ca35c3)"
  select                     = ["firstname", "lastname", "createdon"]
  top                        = 2
  return_total_records_count = true
  //query          = "systemusers?$select=fullname,systemuserid,createdon&$top=3&$expand=systemuserroles_association($select=name)"
}
