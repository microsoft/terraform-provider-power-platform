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

resource "random_string" "random_suffix" {
  length  = 5
  special = false
  upper   = false
}

data "powerplatform_environments" "all_environments" {}

resource "powerplatform_data_record" "data_record_sample_contact1" {
  environment_id     = data.powerplatform_environments.all_environments.environments[0].id
  table_logical_name = "contact"
  columns = {
    firstname          = "John"
    lastname           = "Doe ${random_string.random_suffix.result}"
    telephone1         = "555-555-5555"
    emailaddress1      = "johndoe@contoso.com"
    address1_composite = "123 Main St\nRedmond\nWA\n98052\nUS"
    anniversary        = "2024-04-10"
    annualincome       = 1234.56
    birthdate          = "2024-04-10"
    description        = "This is the description of the the terraform \n\nsample contact"
  }
}

resource "powerplatform_data_record" "data_record_sample_contact2" {
  environment_id     = data.powerplatform_environments.all_environments.environments[0].id
  table_logical_name = "contact"
  columns = {
    firstname          = "Jane"
    lastname           = "Doe ${random_string.random_suffix.result}"
    telephone1         = "555-555-5555"
    emailaddress1      = "janedoe@contoso.com"
    address1_composite = "123 Main St\nRedmond\nWA\n98052\nUS"
    anniversary        = "2024-04-11"
    annualincome       = 1234.56
    birthdate          = "2024-04-11"
    description        = "This is the description of the the terraform \n\nsample contact"
  }
}

resource "powerplatform_data_record" "data_record_accounts" {
  environment_id     = data.powerplatform_environments.all_environments.environments[0].id
  table_logical_name = "account"
  columns = {
    name                = "Sample Account ${random_string.random_suffix.result}"
    creditonhold        = false
    address1_latitude   = 47.639583
    description         = "This is the description of the sample account"
    revenue             = 5000000
    accountcategorycode = 1
  }
}

resource "powerplatform_data_record" "data_record_sample_account" {
  environment_id     = data.powerplatform_environments.all_environments.environments[0].id
  table_logical_name = "account"
  columns = {
    name                = "Terraform Sample Account ${random_string.random_suffix.result}"
    description         = "This is the description of the the terraform sample account"
    accountratingcode   = 1
    revenue             = 123456
    accountcategorycode = 1
    accountcategorycode = 1
    creditonhold        = true
    creditlimit         = 123456
    customersizecode    = 1
    donotbulkemail      = true
    donotfax            = true
    emailaddress1       = "johndoe@contoso.com"
    exchangerate        = 1.0
    ftpsiteurl          = "https://www.contoso.com"
    websiteurl          = "https://www.contoso.com"
    industrycode        = 8
    lastusedincampaign  = "2024-04-10T22:32:12Z"
    lastonholdtime      = "2024-04-10T22:32:12Z"
    telephone1          = "555-555-5555"

    primarycontactid = {
      entity_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
      data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
    }

    contact_customer_accounts = [
      {
        entity_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
        data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
      },
      {
        entity_logical_name = powerplatform_data_record.data_record_sample_contact2.table_logical_name
        data_record_id      = powerplatform_data_record.data_record_sample_contact2.id
      }
    ]

    address1_addresstypecode    = 1
    address1_city               = "Redmond"
    address1_country            = "US"
    address1_county             = "King"
    address1_fax                = "1234567890"
    address1_freighttermscode   = 1
    address1_latitude           = 47.639583
    address1_longitude          = -122.128868
    address1_name               = "Main Office"
    address1_postofficebox      = "123456"
    address1_primarycontactname = "John Doe"
    address1_shippingmethodcode = 1
    address1_stateorprovince    = "WA"
    address1_telephone1         = "1234567890"
    address1_telephone2         = "1234567890"
    address1_telephone3         = "1234567890"
    address1_upszone            = "1234"
    address1_utcoffset          = 5
    address1_line1              = "123 Main St"
    address1_line2              = "Suite 123"
    address1_line3              = "Building 123"

    entityimage = "iVBORw0KGgoAAAANSUhEUgAAAJAAAACQCAYAAADnRuK4AAAAAXNSR0IArs4c6QAAAARnQU1BAACxjwv8YQUAAAAJcEhZcwAADsMAAA7DAcdvqGQAAAMLSURBVHhe7dZBStQBGMbhWZg6EC66gWFHcDZ5gLbagcToJjF2Axc2RwikO4RthAajZU6O/Bl4C8K3TcH/eeC9wMdv8U0AAAD4H3w7Ofh4v89j3PJ4/3A4w8bpYjo7W+xcj3FvFrtXwxke7+GYr1+sxrjbk+cvhzNsnH7YPTpbbK/GuZ3r4QyPJ6AkoJKAkoBKAkoCKgkoCagkoCSgkoCSgEoCSgIqCSgJqCSgJKCSgJKASgJKAioJKAmoJKAkoJKAkoBKAkoCKgkoCagkoCSgkoCSgEoCSgIqCSgJqCSgJKCSgJKASgJKAioJKAmoJKAkoJKAkoBKAkoCKgkoCagkoCSgkoCSgEoCSgIqCSgJqCSgJKCSgJKASgJKAioJKAmoJKAkoJKAkoBKAkoCKgkoCagkoCSgkoCSgEoCSgIqCSgJqCSgJKCSgJKASgJKAioJKAmoJKAkoJKAkoBKAkoCKgkoCagkoCSgkoCSgEoCSgIqCSgJqCSgJKCSgJKASgJKAioJKAmoJKAkoNL34/3D9SHHuJtXB3vDGTZOLyZ764jGuLeL6Ww4AwAAfzKdL2db89ujMW7y7ua3J/ru4tne6v6hHOPu/uaJ3j7/ev3kfLka4x4i+sXDMS93V2Pcj8vpl+EMjyegJKCSgJKASgJKAioJKAmoJKAkoJKAkoBKAkoCKgkoCagkoCSgkoCSgEoCSgIqCSgJqCSgJKCSgJKASgJKAioJKAmoJKAkoJKAkoBKAkoCKgkoCagkoCSgkoCSgEoCSgIqCSgJqCSgJKCSgJKASgJKAioJKAmoJKAkoJKAkoBKAkoCKgkoCagkoCSgkoCSgEoCSgIqCSgJqCSgJKCSgJKASgJKAioJKAmoJKAkoJKAkoBKAkoCKgkoCagkoCSgkoCSgEoCSgIqCSgJqCSgJKCSgJKASgJKAioJKAmoJKAkoJKAkoBKAkoCKgkoCagkoCSgkoCSgEoCSgIqCSgJqCSgJKDSzvvl1TqiMW5rvpwNZ9i4W0xn60OOc08/DWcAAADgH5pMfgLTxq+1aamfnwAAAABJRU5ErkJggg=="
  }
}

resource "powerplatform_data_record" "data_record_tabletwos" {
  environment_id     = data.powerplatform_environments.all_environments.environments[0].id
  table_logical_name = "cr4d0_tabletwo"
  columns = {
    cr4d0_tabletwoid = "21715311-9ff6-ee11-a1fd-7c1e5217db96"
    cr4d0_name       = "Set ${random_string.random_suffix.result}"
  }
}
