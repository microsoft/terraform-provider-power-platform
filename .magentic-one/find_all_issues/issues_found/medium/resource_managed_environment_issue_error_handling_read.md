# Title

Incomplete Error Handling in `Read` Method

##

`/workspaces/terraform-provider-power-platform/internal/services/managed_environment/resource_managed_environment.go`

## Problem

In the `Read` method, error handling is incomplete when calling the `GetEnvironment` method. Specifically, the logic assumes that a missing resource (`ERROR_OBJECT_NOT_FOUND`) should remove it from the state but does not log or track this action clearly.

## Impact

This incomplete error handling can result in silent failures, making debugging difficult, or it could lead to the resource being unintentionally removed from Terraform state.

Severity is assessed as `medium` because it does not break functionality but can complicate debugging.

## Location

**File:** `/workspaces/terraform-provider-power-platform/internal/services/managed_environment/resource_managed_environment.go`  
**Method:** `Read`

## Code Issue

```go
if err != nil {
    if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
        resp.State.RemoveResource(ctx)
        return
    }
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
    return
}
```

## Fix

Add explicit logging to track when a resource is removed intentionally from the state due to being not found:

```go
if err != nil {
    if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
        tflog.Warn(ctx, fmt.Sprintf("Resource not found: %s. Removing it from the state.", r.FullTypeName()))
        resp.State.RemoveResource(ctx)
        return
    }
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
    return
}
```