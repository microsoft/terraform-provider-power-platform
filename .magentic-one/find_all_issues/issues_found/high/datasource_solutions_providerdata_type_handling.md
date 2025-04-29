# Title

Incorrectly Handled Pointer Type Assertions in `Configure`

##

Path: `/workspaces/terraform-provider-power-platform/internal/services/solution/datasource_solutions.go`

## Problem

In the `Configure` method, there is a type assertion for `ProviderData` to a pointer type (`*api.ProviderClient`). However, the assertion does not verify whether the asserted value is `nil`. This omission could lead to unexplained failures when `ProviderData` is non-`nil` but doesn't hold the expected data.

## Impact

While rare, this issue can lead to a panic during runtime if someone provides an invalid but non-nil value for `ProviderData`. Severity is high as this can cause a crash.

## Location

Function `Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)`.

## Code Issue

```go
client, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
    resp.Diagnostics.AddError(
        "Unexpected ProviderData Type",
        fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )
    return
}
```

## Fix

Before performing the type assertion, confirm that the value is non-`nil` and add a check to handle cases where the pointer is invalid.

```go
if req.ProviderData == nil {
    resp.Diagnostics.AddError(
        "Invalid ProviderData",
        "Expected non-nil ProviderData, but got nil. Ensure provider configuration is correctly set.",
    )
    return
}

client, ok := req.ProviderData.(*api.ProviderClient)
if !ok || client == nil {
    resp.Diagnostics.AddError(
        "Unexpected ProviderData Type",
        fmt.Sprintf("Expected non-nil *api.ProviderClient, got: %T. Ensure provider configuration is correctly set.", req.ProviderData),
    )
    return
}
```