# Title

Missing Context Management in HTTP Mock Activation

##

/workspaces/terraform-provider-power-platform/internal/services/application/resource_environment_application_package_install_test.go

## Problem

The code activates the `httpmock` library for mocking HTTP responses but relies solely on manual activation and resetting within test functions:

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()
```

This approach may lead to unintentional state persistence between tests, especially if a test fails prematurely. Testing environments require robust context management to ensure proper initialization and teardown.

## Impact

- **Medium Severity**: Persistent mock states may interfere with subsequent tests, leading to unreliable results and false positives/negatives.
- Improper resource cleanup can cause instability during parallel test execution.

## Location

This pattern exists in multiple locations, such as:

```go
func TestUnitEnvironmentApplicationPackageInstallResource_Validate_Install(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    ...
}
```

Another instance:
```go
func TestUnitEnvironmentApplicationPackageInstallResource_Validate_No_Dataverse(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    ...
}
```

## Fix

Utilize a testing framework's `setup` and `teardown` capabilities or create utility functions to manage HTTP mock activation and deactivation consistently across tests.

A possible utility function:

```go
func withHttpMock(t *testing.T, testFunc func()) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    testFunc()
}
```

Refactor tests to utilize the function:

```go
func TestUnitEnvironmentApplicationPackageInstallResource_Validate_Install(t *testing.T) {
    withHttpMock(t, func() {
        mocks.ActivateEnvironmentHttpMocks()
        ...
    })
}
```

This approach maintains mock isolation, ensures clean context switching, and promotes code reuse.