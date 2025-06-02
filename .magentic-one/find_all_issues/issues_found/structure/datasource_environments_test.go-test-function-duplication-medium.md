# Title

Test Function Duplication Between Acceptance and Unit Tests

##

/workspaces/terraform-provider-power-platform/internal/services/environment/datasource_environments_test.go

## Problem

There is overlap between what the acceptance test `TestAccEnvironmentsDataSource_Basic` and the unit test `TestUnitEnvironmentsDataSource_Validate_Read` do: both execute the data source, and check for a variety of returned fields. This duplication increases maintenance overhead, and makes it harder to update tests since business logic changes and field expectations must be synchronized between both tests. Ideally, validations of field mapping and value-extraction logic should mostly be covered in one place, with error-path and infrastructure integration in another.

## Impact

Severity: **Medium**

Maintaining two near-duplicate test functions with long assertion lists creates maintainability/friction challenges and risks inconsistencies creeping in. Test logic drift could cause future bugs to slip through, as only one test may receive updates for new/changed fields or behavior.

## Location

- TestAccEnvironmentsDataSource_Basic
- TestUnitEnvironmentsDataSource_Validate_Read

## Code Issue

```go
func TestAccEnvironmentsDataSource_Basic(t *testing.T) {
    // lots of field assertions, maps to resource and data source outputs
}

func TestUnitEnvironmentsDataSource_Validate_Read(t *testing.T) {
    // very similar: asserts lots of individual fields for returned environments
}
```

## Fix

To reduce maintenance overhead:
- Extract common result-checking helpers for asserting expected outputs/fields, and use those helpers in both test functions.
- Consider limiting the acceptance test to a high-level integration check, and keeping detailed field-level assertions in the unit tests (or vice versa)
- Review duplication regularly as more test scenarios are added

Example helper:

```go
func AssertExpectedEnvironmentFields(t *testing.T, envIndex int, expected map[string]string) resource.TestCheckFunc {
    // build and return a list of TestCheckResourceAttr for each key/value
}

// usage:
Check: resource.ComposeTestCheckFunc(
    AssertExpectedEnvironmentFields(t, 0, map[string]string{ ... }),
),
```
