# Title

Readability Issues in Terraform Configuration within Tests

##

/workspaces/terraform-provider-power-platform/internal/services/application/datasource_environment_application_packages_test.go

## Problem

The inline Terraform configuration used within test steps is highly nested and lacks proper formatting. This makes it less readable and harder to debug when tests fail. While functional for small configurations, larger configurations like these benefit from being split into multiple lines for improved clarity.

## Impact

Poor readability reduces developer efficiency when debugging or updating tests. The severity is **Low**, as this does not impact functionality but slightly affects maintainability.

## Location

File: datasource_environment_application_packages_test.go
Functions:
1. `TestAccEnvironmentApplicationPackagesDataSource_Validate_Read`
2. `TestUnitEnvironmentApplicationPackagesDataSource_Validate_Read`

## Code Issue

```go
TestAccEnvironmentApplicationPackagesDataSource_Validate_Read:
Config: `resource "powerplatform_environment" "env" { display_name = "Test Name" location = "US" environment_type = "Sandbox" ... }

TestUnitEnvironmentApplicationPackagesDataSource_Validate_Read:
Config: `resource "powerplatform_environment" "env" { display_name ... }
```

## Fix

Reformat the Terraform configurations into multiline blocks for improved readability:

```go
Config: `
resource "powerplatform_environment" "env" {
    display_name      = "Test Name"
    location          = "US"
    environment_type  = "Sandbox"
    dataverse = {
        language_code     = "1033"
        currency_code     = "USD"
        security_group_id = "00000000-0000-0000-0000-000000000000"
    }
}

data "powerplatform_environment_application_packages" "all_applications" {
    environment_id = powerplatform_environment.env.id
}`
```

This formatting makes it easier to read and update the test cases in the future.