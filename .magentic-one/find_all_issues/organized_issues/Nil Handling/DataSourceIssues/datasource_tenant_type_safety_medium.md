# Type assertion without error handling in Configure

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant.go

## Problem

In the `Configure` method, the code asserts that `req.ProviderData` is of type `*api.ProviderClient` via:

```go
client, ok := req.ProviderData.(*api.ProviderClient)
```

If the assertion fails, it adds an error to diagnostics, but does not actually return early or halt execution. Instead, the code proceeds, which may lead to further logic executing with a nil or incorrect client, resulting in nil reference errors or unexpected behaviors elsewhere.

## Impact

Medium: Adds error diagnostics but possible logic flow may continue with invalid state, risking subtle bugs or panics later in the lifecycle.

## Location

In `Configure`:

## Code Issue

```go
client, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
    resp.Diagnostics.AddError(
        "Unexpected ProviderData Type",
        fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )
    // Should return here to avoid further usage
}

d.TenantClient = NewTenantClient(client.Api)
```

## Fix

Return immediately after adding diagnostics for the assertion failure:

```go
if !ok {
    resp.Diagnostics.AddError(
        "Unexpected ProviderData Type",
        fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )
    return // <-- ensure no further logic executes
}
```
