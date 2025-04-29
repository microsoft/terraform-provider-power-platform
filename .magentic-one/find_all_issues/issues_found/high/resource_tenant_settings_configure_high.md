# Title

Improper Error Handling in `Configure` Function

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go`

## Problem

The `Configure` function performs a type assertion on `req.ProviderData` to derive an API client. While it checks if the assertion fails, it doesn't clarify why this might happen or provide diagnostic messages that could assist debugging efforts. The error message provided is generic and doesn't offer guidance on how to resolve issues arising from improperly initialized or incompatible provider data.

## Impact

This lack of transparency impacts debugging and user experience significantly, especially if `ProviderData` is incompatible or improperly initialized. Critical reliance on type assertions without proper diagnostics greatly increases ambiguity when issues arise. Severity: **High**.

## Location

Line 256: Inside the `Configure` method.

## Code Issue

```go
client := req.ProviderData.(*api.ProviderClient).Api

if client == nil {
    resp.Diagnostics.AddError(
        "Unexpected Resource Configure Type",
        fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )
    return
}
```

## Fix

Provide a more detailed debugging message that includes a potential solution or sets diagnostic expectations. Additionally, refactor `ProviderData` validation to ensure compatible initialization before the type assertion.

```go
providerClient, ok := req.ProviderData.(*api.ProviderClient)
if !ok || providerClient.Api == nil {
    resp.Diagnostics.AddError(
        "Unexpected ProviderData Type or Uninitialized API Client",
        fmt.Sprintf("Expected *api.ProviderClient with a valid API client, got type: %T. Ensure the provider is properly initialized in configuration.", req.ProviderData),
    )
    return
}
client := providerClient.Api
```