# Title

Improper CapacityClient Initialization in the `Configure` Method

##

`/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity.go`

## Problem

The `newCapacityClient(client.Api)` invocation in the `Configure` method assumes that the `Api` property on `client` is correctly set without validating its state. If `client.Api` is nil or in an invalid state, it will cause runtime errors during execution.

## Impact

This issue leads to fragility in the data source configuration process. If `client.Api` is not correctly initialized, the provider might crash or fail to properly configure. **Severity is high**, as this impacts initial setup and configuration.

## Location

Line in the `Configure` method where `newCapacityClient(client.Api)` is invoked.

## Code Issue

```go
d.CapacityClient = newCapacityClient(client.Api)
```

## Fix

Validate the state of `client.Api` before passing it to the `newCapacityClient` function to avoid potential runtime errors.

```go
if client.Api == nil {
    resp.Diagnostics.AddError(
        "Invalid API Client",
        "The API client is uninitialized. Ensure the provider is correctly configured.",
    )
    return
}
d.CapacityClient = newCapacityClient(client.Api)
```

This ensures the API client is always valid and prevents potential crashes due to nil or invalid references.
