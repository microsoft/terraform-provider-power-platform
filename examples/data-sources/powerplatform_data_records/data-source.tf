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



# resource "powerplatform_data_record" "sub_business_unit_1" {
#   environment_id     = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
#   table_logical_name = "businessunit"
#   columns = {
#     name       = "sub buissness unit 1"
#     costcenter = "cost center 1"

#     parent_business_unit = {
#       data_record_id     = powerplatform_data_record.main_business_unit.id
#       table_logical_name = "businessunit"
#     }
#   }
# }

# resource "powerplatform_data_record" "sub_business_unit_2" {
#   environment_id     = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
#   table_logical_name = "businessunit"
#   columns = {
#     name       = "sub buissness unit 1"
#     costcenter = "cost center 1"

#     parent_business_unit = {
#       data_record_id     = powerplatform_data_record.main_business_unit.id
#       table_logical_name = "businessunit"
#     }
#   }
# }

# resource "powerplatform_data_record" "contact1" {
#   environment_id     = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
#   table_logical_name = "contact"
#   columns = {
#     contactid = "00000000-0000-0000-0000-000000000001"

#   }
# }

# resource "powerplatform_data_record" "contact2" {
#   environment_id     = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
#   table_logical_name = "contact"
#   columns = {
#     contactid = "00000000-0000-0000-0000-000000000002"
#     contact_customer_contacts = [
#       {
#         table_logical_name = "contact"
#         data_record_id     = "00000000-0000-0000-0000-000000000001"
#       }
#     ]
#   }
# }

# resource "powerplatform_data_record" "contact3" {
#   environment_id     = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
#   table_logical_name = "contact"
#   columns = {
#     contactid = "00000000-0000-0000-0000-000000000003"
#     firstname = "contact3"
#     lastname  = "contact3"
#   }
# }

# data "powerplatform_data_records" "data_query" {
#   environment_id    = "a1e605fb-80ad-e1b2-bae0-f046efc0e641" //powerplatform_environment.data_env.id
#   entity_collection = "contacts(00000000-0000-0000-0000-000000000001)"
#   select            = ["contactid", "firstname", "lastname"]
#   expand = toset([
#     {
#       navigation_property = "contact_customer_contacts"
#       select              = ["contactid", "firstname", "lastname"]
#     },
#   ])

#   depends_on = [
#     powerplatform_data_record.contact1,
#     powerplatform_data_record.contact2,
#     powerplatform_data_record.contact3
#   ]
# }
