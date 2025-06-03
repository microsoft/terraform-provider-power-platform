# Test with Side-Effects Not Cleaned Beyond httpmock Deactivation

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/resource_environment_settings_test.go

## Problem

Some tests activate httpmock (which intercepts HTTP calls) but do not guarantee cleanup if the test panics before reaching `defer httpmock.DeactivateAndReset()`. Additionally, the `mocks.ActivateEnvironmentHttpMocks()` call in `TestUnitTestEnvironmentSettingsResource_Validate_No_Dataverse` isn’t cleaned up after.

## Impact

Medium severity: If test panics or other execution paths bypass the cleanup, other tests may fail in unpredictable ways because global hooks remain in effect.

## Location

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()
// ...tests...
// also: mocks.ActivateEnvironmentHttpMocks() // no corresponding cleanup
```

## Code Issue

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()

mocks.ActivateEnvironmentHttpMocks() // not cleaned up
```

## Fix

Ensure all side-effects are cleaned up and, where necessary, provide a cleanup method for `mocks.ActivateEnvironmentHttpMocks` if it changes global state.

```go
httpmock.Activate()
t.Cleanup(httpmock.DeactivateAndReset)

mocks.ActivateEnvironmentHttpMocks()
t.Cleanup(mocks.DeactivateEnvironmentHttpMocks) // if such a function exists
```

This guarantees cleanup on test abort.
