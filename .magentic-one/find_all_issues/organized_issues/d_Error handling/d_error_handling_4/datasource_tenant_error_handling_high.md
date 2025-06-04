# Error handling missing when setting state

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant.go

## Problem

In the `Read` method, `resp.State.Set(ctx, &state)` is called, but its returned error is not checked. Failing to check this error could result in silent failures where state is not set properly, causing undetected problems for Terraform users and potentially confusing diagnostics.

## Impact

The impact is high since Terraform would not surface clear errors to the user if state marshalling/set fails. This could result in unreliable provider behavior and unpredictable failures.

## Location

In the `Read` method of the file `/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant.go`, near the end:

## Code Issue

```go
resp.State.Set(ctx, &state)
```

## Fix

Check the error returned by `Set` and add diagnostics if it fails:

```go
if err := resp.State.Set(ctx, &state); err != nil {
    resp.Diagnostics.AddError("Failed to set state", err.Error())
}
```
