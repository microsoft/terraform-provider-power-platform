# Inefficient Test Initialization and Teardown

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy_test.go

## Problem

Each test repeats HTTP mock activation, responder registration, and teardown logic. This could be abstracted via helper utilities or test setup methods to ensure DRY (Don’t Repeat Yourself) principles.

## Impact

- **Maintainability**: Extra boilerplate and higher chances to forget proper setup in new tests.
- **Test Duration**: Slightly slower due to repeated setup.

**Severity: Low**

## Location

Repeated blocks like:

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()
mocks.ActivateEnvironmentHttpMocks()
// register responders...
```
```

## Code Issue

Present in each `TestUnitTest...` function.

## Fix

Create helper functions for test initialization.

```go
func setupHTTPMock(t *testing.T) {
    httpmock.Activate()
    t.Cleanup(httpmock.DeactivateAndReset)
    mocks.ActivateEnvironmentHttpMocks()
}

func TestSomeAPICall(t *testing.T) {
    setupHTTPMock(t)
    // ...
}
```
