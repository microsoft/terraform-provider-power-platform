# Title

Repeated Activation and Deactivation of `httpmock` with Manual Cleanup

## Path

`/workspaces/terraform-provider-power-platform/internal/provider/provider_test.go`

## Problem

The section of the code activates `httpmock` and manually deactivates it after use. This approach is error-prone as it relies on developers to consistently use `defer httpmock.DeactivateAndReset()` to ensure cleanup is executed.

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()
```

## Impact

- **Potential Resource Leaks**: If cleanup with `Defer` is missed, test functions might result in a polluted HTTP mock environment, impacting subsequent tests.
- **Reduced Test Isolation**: Tests may fail unpredictably due to external interference.
- **Severity Critical**: Proper test resource isolation is critical for reliable unit testing.

## Location

- Line 102 in function `TestUnitPowerPlatformProvider_Validate_Telemetry_Optout_Is_False`
- Line 128 in function `TestUnitPowerPlatformProvider_Validate_Telemetry_Optout_Is_True`

## Code Issue

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()
```

## Fix

Use a helper function to automatically activate and safely deactivate `httpmock`. This ensures cleanup is consistently applied and eliminates any human error.

```go
func withHTTPMock(f func()) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    f()
}

// Usage example within tests:
func TestUnitPowerPlatformProvider_Validate_Telemetry_Optout_Is_False(t *testing.T) {
    withHTTPMock(func() {
        mocks.ActivateEnvironmentHttpMocks()
        httpmock.RegisterRegexpResponder(
            "GET",
            regexp.MustCompile(`^https://api\\.bap\\.microsoft\\.com/providers/Microsoft\\.BusinessAppPlatform/scopes/admin/environments.*api-version=2023-06-01$`),
            func(req *http.Request) (*http.Response, error) {
                return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../services/environment/tests/datasource/Validate_Read/get_environments.json").String()), nil
            },
        )
        test.Test(t, test.TestCase{/* ... */})
    })
}
```