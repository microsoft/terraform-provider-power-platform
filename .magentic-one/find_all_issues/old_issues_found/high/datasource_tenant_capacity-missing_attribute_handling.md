# Title

Missing Error Handling for `GetAttribute` in `Read` Method

##

`/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity.go`

## Problem

The `req.Config.GetAttribute(ctx, path.Root("tenant_id"), &tenantId)` call in the `Read` method assumes that the attribute 'tenant_id' will always exist without checking for errors. If the attribute is missing or invalid, it could lead to unexpected behavior or runtime errors.

## Impact

This omission could lead to unpredictable runtime issues, such as crashes or improper functioning of the `Read` method, potentially causing data corruption or failures. **Severity is high**, as error handling in configuration fetching is critical for robustness.

## Location

Line in the `Read` method where `GetAttribute` is invoked.

## Code Issue

```go
req.Config.GetAttribute(ctx, path.Root("tenant_id"), &tenantId)
```

## Fix

Add error handling for the `GetAttribute` function call to ensure robustness. If the attribute is missing or invalid, an appropriate error should be logged, and processing should halt gracefully.

```go
diags := req.Config.GetAttribute(ctx, path.Root("tenant_id"), &tenantId)
if diags.HasError() {
    resp.Diagnostics.Append(diags...)
    resp.Diagnostics.AddError(
        "Configuration Error",
        "Unable to fetch 'tenant_id' attribute from the configuration.",
    )
    return
}
```

This ensures that the method gracefully handles cases where the attribute is missing or incorrectly configured.
