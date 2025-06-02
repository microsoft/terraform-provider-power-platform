# Title

Potential Resource Leak in HTTP Mock Activation/Deactivation

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/datasource_environment_settings_test.go

## Problem

In the unit tests, `httpmock.Activate()` and `defer httpmock.DeactivateAndReset()` are used to enable HTTP mocking. If an early return, panic, or failed test setup occurs before the `defer` can execute, it's possible that the mock will not be correctly deactivated — especially if future refactors introduce returns before the defer.

## Impact

While currently unlikely due to test structure, if this pattern is copied elsewhere with less strict control flow, it may result in global state leak from shared test resources, causing flaky tests or confusion in other packages. Severity: Low, since current implementation guards against it, but the pattern needs care.

## Location

In both unit tests:
- `TestUnitTestEnvironmentSettingsDataSource_Validate_Read`
- `TestUnitTestEnvironmentSettingsDataSource_Validate_No_Dataverse`

## Code Issue

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()
```

## Fix

Consider using a test-helper function that wraps the test logic, explicitly handling setup/teardown, or make sure to confirm that panics or early-exit scenarios cannot occur before the `defer`.

No code change is strictly required, but updating with a helper for clarity:

```go
func withHTTPMock(t *testing.T, testFunc func()) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    testFunc()
}

// Usage:

func TestUnitTestEnvironmentSettingsDataSource_Validate_Read(t *testing.T) {
    withHTTPMock(t, func() {
        ...
    })
}
```

Alternatively, validate at the start of each test that setup succeeded and catching panics for assertable recovery.

---

This file will be saved as:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/datasource_environment_settings_test.go_httpmock_resource_leak_low.md`
