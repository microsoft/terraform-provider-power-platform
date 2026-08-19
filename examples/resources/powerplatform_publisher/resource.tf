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

resource "powerplatform_environment" "example_publisher" {
  display_name     = "example_publisher"
  location         = "unitedstates"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "powerplatform_publisher" "example" {
  environment_id         = powerplatform_environment.example_publisher.id
  uniquename             = "contoso"
  friendly_name          = "Contoso Publisher"
  customization_prefix   = "cts"
  description            = "Terraform-managed Dataverse publisher"
  email_address          = "publisher@contoso.example"
  supporting_website_url = "https://contoso.example"

  address = [
    {
      slot        = 1
      line1       = "1 Collins Street"
      city        = "Melbourne"
      country     = "Australia"
      postal_code = "3000"
      telephone1  = "+61-3-5555-0101"
    },
    {
      slot        = 2
      line1       = "100 Queen Street"
      city        = "Auckland"
      country     = "New Zealand"
      postal_code = "1010"
      telephone1  = "+64-9-555-0102"
    }
  ]
}
