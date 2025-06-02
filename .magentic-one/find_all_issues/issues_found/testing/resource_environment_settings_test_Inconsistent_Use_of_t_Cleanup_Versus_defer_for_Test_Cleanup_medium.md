# Inconsistent Use of `t.Cleanup` Versus `defer` for Test Cleanup

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/resource_environment_settings_test.go

## Problem

Some uses cleanup with `defer`, but Go 1.14+ has `t.Cleanup`, which executes even if the test fails or panics, ensuring cleanup runs reliably in subtests and future test composition.

## Impact

Medium; mostly makes tests more robust and future-proof for subtest expansion.

## Location

All uses like below:

## Code Issue

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()
```

## Fix

Switch to:

```go
httpmock.Activate()
t.Cleanup(httpmock.DeactivateAndReset)
```
