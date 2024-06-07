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

resource "powerplatform_data_record" "contact1" {
  environment_id     = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
  table_logical_name = "contact"

  columns = {
    contactid = "00000000-0000-0000-0000-000000000001"
    firstname = "contact1"
    lastname  = "contact1"

    contact_customer_contacts = [
      {
        table_logical_name = powerplatform_data_record.contact2.table_logical_name
        data_record_id     = powerplatform_data_record.contact2.columns.contactid
      },
      {
        table_logical_name = powerplatform_data_record.contact3.table_logical_name
        data_record_id     = powerplatform_data_record.contact3.columns.contactid
      }
    ]
  }
}

resource "powerplatform_data_record" "contact2" {
  environment_id     = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
  table_logical_name = "contact"
  columns = {
    contactid = "00000000-0000-0000-0000-000000000002"
    firstname = "contact2"
    lastname  = "contact2"

    contact_customer_contacts = [
      {
        table_logical_name = powerplatform_data_record.contact4.table_logical_name
        data_record_id     = powerplatform_data_record.contact4.columns.contactid
      }
    ]
  }
}

resource "powerplatform_data_record" "contact3" {
  environment_id     = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
  table_logical_name = "contact"
  columns = {
    contactid = "00000000-0000-0000-0000-000000000003"
    firstname = "contact3"
    lastname  = "contact3"
  }
}

resource "powerplatform_data_record" "contact4" {
  environment_id     = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
  table_logical_name = "contact"
  columns = {
    contactid = "00000000-0000-0000-0000-000000000004"
    firstname = "contact4"
    lastname  = "contact4"
    account_primary_contact = [
      {
        table_logical_name = powerplatform_data_record.account1.table_logical_name
        data_record_id     = powerplatform_data_record.account1.columns.accountid
      }
    ]
  }
}



resource "powerplatform_data_record" "account1" {
  environment_id     = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
  table_logical_name = "account"
  columns = {
    accountid = "00000000-0000-0000-0000-000000000010"
    name      = "account1"
    contact_customer_accounts = [
      {
        table_logical_name = powerplatform_data_record.contact5.table_logical_name
        data_record_id     = powerplatform_data_record.contact5.columns.contactid
      }
    ]
  }
}

resource "powerplatform_data_record" "contact5" {
  environment_id     = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
  table_logical_name = "contact"
  columns = {
    contactid = "00000000-0000-0000-0000-000000000005"
    firstname = "contact5"
    lastname  = "contact5"
  }
}

data "powerplatform_data_records" "data_query" {
  environment_id          = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
  entity_collection       = "contacts"
  select                  = ["fullname", "firstname", "lastname"]
  apply                   = "groupby((statuscode),aggregate($count as count))"
  return_total_rows_count = true

  depends_on = [
    powerplatform_data_record.contact1,
    powerplatform_data_record.contact2,
    powerplatform_data_record.contact3,
    powerplatform_data_record.contact4,
    powerplatform_data_record.contact5,
    powerplatform_data_record.account1,
  ]
}
