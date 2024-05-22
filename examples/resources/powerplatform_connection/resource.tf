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

# resource "powerplatform_connection" "example" {
#   environment_id = "838f76c8-a192-e59c-a835-089ad8cfb047"
#   name           = "shared_servicebus"
#   display_name   = "blablabla"
#   parameters     = "{\"name\":\"connectionstringauth\",\"values\":{\"ConnectionString\":{\"value\":\"dadsadsadsds\"}}}"
# }

resource "powerplatform_connection" "example1" {
  environment_id = "838f76c8-a192-e59c-a835-089ad8cfb047"
  name           = "shared_sql"
  display_name   = "sql123"
}
