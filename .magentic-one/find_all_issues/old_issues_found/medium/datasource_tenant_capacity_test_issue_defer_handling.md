# Title

Improper Defer Call Without Explicit Error Handling

##

/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity_test.go

## Problem

The test function `TestUnitTenantCapacityDataSource_Validate_Read` activates HTTP mocking (`httpmock.Activate()`) and defers its cleanup using `httpmock.DeactivateAndReset()`. However, there is no validation to ensure the `Activate` and `DeactivateAndReset` calls execute successfully, nor does it handle potential errors in activation.

## Impact

If `httpmock.Activate` encounters an issue during initialization, it may silently fail without a clear indication to the developer. Additionally, if `DeactivateAndReset` isn't properly executed, it could leave mocking artifacts in subsequent tests, leading to unreliable test results. The severity of this issue is **medium**, as it impacts test consistency and transparency.

## Location

Within the function `TestUnitTenantCapacityDataSource_Validate_Read`:

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()
```

## Code Issue

The code block does not validate the success of the `httpmock.Activate` function call and lacks recovery mechanisms.

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()
```

## Fix

Introduce error checking for `httpmock.Activate` and adjust the deferred call to use conditional cleanup, ensuring robust error handling:

```go
if err := httpmock.Activate(); err != nil {
    t.Fatalf("Failed to activate HTTP mocking: %v", err)
}
defer func() {
    if !httpmock.IsActivated() {
        t.Log("HTTP mocking was not activated successfully.")
        return
    }
    httpmock.DeactivateAndReset()
}()
```

- This fix ensures the activation success is validated and any failure is promptly reported, improving the test’s reliability.
- The deferred cleanup is now context-aware, avoiding redundant or risky deactivation calls.
