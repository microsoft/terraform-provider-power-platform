# Title

Lack of Error Recovery in `Read` Method for Nonexistent Resources

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy.go

## Problem

In the `Read` method, the code checks if the error is of type `customerrors.ERROR_OBJECT_NOT_FOUND` and calls `resp.State.RemoveResource(ctx)` without returning additional diagnostics or recovery information. This omits clarity about the nature of the error or steps that developers might take.

## Impact

While functional, this lacks an informative diagnostic description for users and developers. Enhanced error messages and structured diagnostics can make the resource's behavior clearer in scenarios where reading fails. Severity: **medium**.

## Location

```go
if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
    resp.State.RemoveResource(ctx)
    return
}
```

## Fix

Add detailed diagnostics explaining the removal of the resource:

```go
if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
    resp.Diagnostics.AddWarning(
        "Resource Not Found",
        fmt.Sprintf("No Billing Policy found with ID: %s. Removing resource state.", state.Id.ValueString()),
    )
    resp.State.RemoveResource(ctx)
    return
}
```

This improves transparency and ensures users understand the consequences of the missing resource.