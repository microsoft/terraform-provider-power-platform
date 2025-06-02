# Issue: Test API—Missing Cleanup or Isolation for Shared State

## 
/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/resource_enterprise_policy_test.go

## Problem

The test sets up HTTP mocks globally and calls `mocks.ActivateEnvironmentHttpMocks()`, but may not restore all state between runs, risking cross-test pollution if more tests are added in the same package/file.

## Impact

This may cause flaky tests if state leaks between them, especially if more tests are written. Severity: medium.

## Location

Top of function:

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()
mocks.ActivateEnvironmentHttpMocks()
```

## Code Issue

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()

mocks.ActivateEnvironmentHttpMocks()
```

## Fix

Ensure all mocks and related state are cleanly set up and reset, including any global variables set by `mocks.ActivateEnvironmentHttpMocks()`. If `mocks` has a reset or teardown, defer it as well (example below, but adjust based on actual implementation):

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()

mocks.ActivateEnvironmentHttpMocks()
defer mocks.DeactivateEnvironmentHttpMocks() // If available
```

If no deactivate function, document the necessity for test isolation.
