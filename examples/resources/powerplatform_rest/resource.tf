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

resource "powerplatform_environment" "env" {
  display_name     = "sample_data_environment"
  location         = "europe"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}

resource "powerplatform_rest" "install_sample_data" {
  environment_id = powerplatform_environment.env.id
  create = {
    url                  = "/api/data/v9.2/InstallSampleData"
    method               = "POST"
    expected_http_status = [204]
  }
  destroy = {
    url                  = "/api/data/v9.2/UninstallSampleData"
    method               = "POST"
    expected_http_status = [204]
  }
}
