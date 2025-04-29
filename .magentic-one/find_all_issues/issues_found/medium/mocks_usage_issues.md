# Title

Inconsistent Usage of Mock Activation

##

/workspaces/terraform-provider-power-platform/internal/services/application/datasource_environment_application_packages_test.go

## Problem

The activation of mocks for testing interactions relies on the `mocks.ActivateEnvironmentHttpMocks()` call, which appears scattered across multiple functions without centralized control. If the mock activation fails or changes, updating multiple locations in the test code will be required. This practice increases the risk of manual errors and inconsistency.

## Impact

Scattered mock activation logic reduces maintainability and increases the possibility of introducing issues when changes occur. Improper mock usage could lead to unpredictable test results. Severity of the issue: **Medium**.

## Location

File: datasource_environment_application_packages_test.go
Functions:
1. `TestUnitEnvironmentApplicationPackagesDataSource_Validate_Read`
2. `TestUnitEnvironmentApplicationPackagesDataSource_Validate_No_Dataverse`

## Code Issue

```go
mocks.ActivateEnvironmentHttpMocks()
```

## Fix

Encapsulate the mock activation logic into a helper function or utilize setup and teardown functionality to manage it more reliably across multiple test instances.

```go
func setUpMocks() {
    httpmock.Activate()
    mocks.ActivateEnvironmentHttpMocks()
}
func tearDownMocks() {
    httpmock.DeactivateAndReset()
}

// Usage
func TestUnitEnvironmentApplicationPackagesDataSource_Validate_Read(t *testing.T) {
    setUpMocks()
    defer tearDownMocks()
    ...
}
```