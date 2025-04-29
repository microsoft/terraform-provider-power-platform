# Title

Redundant Code in Duplicate Test Steps

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record_test.go

## Problem

Certain test cases repeat almost identical setups and validation logic across multiple steps, leading to redundant code blocks in the file (e.g., `TestAccDataRecordResource_Validate_Create` and `TestUnitDataRecordResource_Validate_Create`).

## Impact

Redundant code increases the file size and leads to maintenance difficulties. Making modifications to common configurations requires updates in multiple places, raising the likelihood of introducing errors. Severity: High.

## Location

Examples include `TestAccDataRecordResource_Validate_Create` steps vs. `TestUnitDataRecordResource_Validate_Create`.

## Code Issue

```go
resource "powerplatform_environment" "test_env" {
    display_name     = "` + mocks.TestName() + `"
    location         = "unitedstates"
    environment_type = "Sandbox"
    dataverse = {
        language_code     = "1033"
        currency_code     = "USD"
        security_group_id = "00000000-0000-0000-0000-000000000000"
    }
}
```

## Fix

Refactor the tests to use shared utilities or functions to define common configurations and validation logic, reducing redundancy. Example:

```go
func generateTestEnvConfig() string {
    return `resource "powerplatform_environment" "test_env" {
        display_name     = "` + mocks.TestName() + `"
        location         = "unitedstates"
        environment_type = "Sandbox"
        dataverse = {
          language_code     = "1033"
          currency_code     = "USD"
          security_group_id = "00000000-0000-0000-0000-000000000000"
        }
    }`
}

func testAccDataRecordValidateCreateConfig() string {
    envConfig := generateTestEnvConfig()
    return envConfig + `
resource "powerplatform_data_record" "data_record_sample_contact1" { ... }
resource "powerplatform_data_record" "data_record_account" { ... }
    `
}
```