# Title

Missing context cancellation handling in API calls

##

/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity.go

## Problem

The `Read` function receives a `context.Context` object, which is routinely passed to downstream calls (such as `GetTenantCapacity`). However, there is no handling for context cancellation or deadline exceeded scenarios. If the context is canceled, the function will still attempt to process results and set state, instead of properly checking for and responding to cancellation.

## Impact

Medium severity. Ignoring context cancellation can lead to unnecessary work, potential resource leaks, and inconsistent state updates, especially for long-running or expensive API calls.

## Location

Function: `Read`

## Code Issue

```go
tenantCapacityDto, err := d.CapacityClient.GetTenantCapacity(ctx, tenantId)
if err != nil {
    resp.Diagnostics.AddError(
        "error fetching tenant capacity",
        err.Error(),
    )
    return
}
// Proceed without checking if err is due to context cancellation
```

## Fix

Check for context cancellation explicitly after the API call:

```go
if err != nil {
    if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
        // Optionally log context cancellation, and return early
        return
    }
    resp.Diagnostics.AddError(
        "error fetching tenant capacity",
        err.Error(),
    )
    return
}
```

This prevents setting partially-computed state and improves robustness when requests are canceled or time out.
