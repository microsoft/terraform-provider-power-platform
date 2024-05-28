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

resource "powerplatform_environment" "data_record_example_env" {
  display_name     = "powerplatform_data_record_example"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}


resource "powerplatform_data_record" "data_record_sample_contact1" {
  environment_id     = powerplatform_environment.data_record_example_env.id
  table_logical_name = "contact"
  columns = {
    firstname          = "John"
    lastname           = "Doe"
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
  environment_id     = powerplatform_environment.data_record_example_env.id
  table_logical_name = "contact"
  columns = {
    firstname          = "Jane"
    lastname           = "Doe"
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
  environment_id     = powerplatform_environment.data_record_example_env.id
  table_logical_name = "account"
  columns = {
    name                = "Sample Account"
    creditonhold        = false
    address1_latitude   = 47.639583
    description         = "This is the description of the sample account"
    revenue             = 5000000
    accountcategorycode = 1

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

    # new_Account_Contact_Contact = [
    #   {
    #     entity_logical_name = powerplatform_data_record.data_record_sample_contact1.table_logical_name
    #     data_record_id      = powerplatform_data_record.data_record_sample_contact1.id
    #   },
    #   # {
    #   #   entity_logical_name = powerplatform_data_record.data_record_sample_contact2.table_logical_name
    #   #   data_record_id      = powerplatform_data_record.data_record_sample_contact2.id
    #   # }
    # ]
  }
}
