# Issue: Potential Leaked Context Resources with Defer

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy.go

## Problem

Each resource method wraps code with a call to `helpers.EnterRequestContext`, which returns a cleanup function that is deferred. However, there are multiple `return` statements before the logic ends (e.g., after error handling or diagnosis checks). In rare cases, resource cleanup may occur later than optimal (if the context carries large values or goroutines), as the function returns earlier. This can be a problem for resource leaks or lock contention in long-lived operations.

## Impact

Severity: **Low**

For most cases in Terraform, this may not be a practical issue (since the function scope is returned shortly afterwards), but for context-bound resources, prompt cleanup is a best practice and makes reasoning about program state easier.

## Location

- In every method that calls:
  ```go
  ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
  defer exitContext()
  ```

## Code Issue

```go
ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
defer exitContext()
if errorCondition {
    return
}
```

## Fix

Consider calling `exitContext()` directly before any `return` that happens before the end of the function. Example:

```go
if errorCondition {
    exitContext()
    return
}
```

Alternatively, review if `exitContext()` must always run at function return, or if the current usage is non-problematic. If so, document expected lifetimes in comments.
