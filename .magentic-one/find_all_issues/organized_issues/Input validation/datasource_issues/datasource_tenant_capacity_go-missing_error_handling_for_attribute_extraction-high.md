# Title

Missing error handling when extracting attribute from config

##

/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity.go

## Problem

In the `Read` function, the code calls `req.Config.GetAttribute(ctx, path.Root("tenant_id"), &tenantId)` to extract the tenant ID from the input configuration. However, it does not check or handle the returned diagnostics or errors from this method. If extraction fails, the function may proceed with an empty or invalid tenantId, resulting in unpredictable runtime errors or API rejections.

## Impact

High severity. Ignoring the result can produce misleading errors, cause API calls to fail, or make debugging difficult. Error handling of configuration extraction is crucial to ensure stable control flow.

## Location

Function: `Read`, during config attribute extraction

## Code Issue

```go
var tenantId string
req.Config.GetAttribute(ctx, path.Root("tenant_id"), &tenantId)
```

## Fix

Capture and handle the diagnostics appropriately:

```go
diags := req.Config.GetAttribute(ctx, path.Root("tenant_id"), &tenantId)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
    return
}
```

This ensures any extraction error halts further processing before making external calls.
