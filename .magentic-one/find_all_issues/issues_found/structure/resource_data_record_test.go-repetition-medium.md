# Repetition of Test Setup and Resource Configuration

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record_test.go

## Problem

There is a large amount of duplicated test setup, especially repeated HCL resource configuration strings and HTTP mock responders among acceptance (`TestAcc*`) and unit (`TestUnit*`) tests as well as across different test functions. For instance, the configuration for `powerplatform_environment`, `powerplatform_data_record`, and HTTP mock registrations are all mostly duplicated, differing only by resource names or minor details. 

This duplication makes it hard to update tests, introduces the risk of inconsistencies, and bloats the test file.

## Impact

- **Maintainability**: Changes in resource shape or expected network interactions require updates in many places.
- **DRY Principle Violation**: Higher risk of inconsistency and accidental drift between similar tests.
- **Readability**: Adds significant noise, making it harder to spot the actual unique intent of a test.
- **Severity**: Low to Medium

## Location

Repeated blocks like:

```go
resource "powerplatform_environment" "test_env" {
  display_name     = "` + mocks.TestName() + `"
  location         = "unitedstates"
  environment_type = "Sandbox"
  // ...
}

resource "powerplatform_data_record" "data_record_sample_contact1" {
  // ...
}
```

and similar repeated HTTP mock responder setups.

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

Move repeated config and responder code into constants and helper functions. For example:

```go
const envConfig = `
resource "powerplatform_environment" "test_env" {
  display_name     = "%s"
  location         = "unitedstates"
  environment_type = "Sandbox"
  dataverse = {
    language_code     = "1033"
    currency_code     = "USD"
    security_group_id = "00000000-0000-0000-0000-000000000000"
  }
}
`

// Usage in test:
Config: fmt.Sprintf(envConfig, mocks.TestName()),

// And for HTTP mock setup:
func registerCommonResponders() {
  httpmock.RegisterResponder("GET", "...", ...)
  // etc.
}
```

--- 

This reduces duplication, makes tests easier to update, and makes intent clearer.

---
