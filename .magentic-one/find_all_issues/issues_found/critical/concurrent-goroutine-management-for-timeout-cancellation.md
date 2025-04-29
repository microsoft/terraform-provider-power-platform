# Title

Concurrent Goroutine Management for Timeout Cancellation

##

`/workspaces/terraform-provider-power-platform/internal/helpers/contexts.go`

## Problem

The returned function for `EnterRequestContext` blindly calls `(*cancel)()` without first checking whether `cancel` is `nil`. This may result in panics if the `cancel` function pointer hasn't been initialized.

## Impact

This has the potential to crash the application with a runtime panic if `(*cancel)()` is executed without proper checks. Severity: **Critical**.

## Location

Function `EnterRequestContext`, located in `/workspaces/terraform-provider-power-platform/internal/helpers/contexts.go`.

## Code Issue

```go
if cancel != nil {
    (*cancel)()
}
```

## Fix

Exactly verify the value of `cancel` before making the function call to ensure its safety.

**Example Fix:**
```go
return ctx, func() {
    tflog.Debug(ctx, fmt.Sprintf("%s END: %s", reqType, name))
    if cancel != nil {
        (*cancel)()
    } else {
        tflog.Warn(ctx, "Cancel function reference is nil - skipping cancel execution.")
    }
}
```