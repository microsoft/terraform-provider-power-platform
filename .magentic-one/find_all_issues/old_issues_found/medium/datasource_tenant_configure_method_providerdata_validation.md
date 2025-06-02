# Title

Lack of full validation for `req.ProviderData` type in `Configure` method

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant.go

## Problem

The `Configure` method assumes that `req.ProviderData` is of type `*api.ProviderClient` without verifying if `req.ProviderData` is either `nil` or an unexpected type before performing the conditional type assertion. While it does check if `req.ProviderData` can successfully be asserted to `*api.ProviderClient`, this happens only after a type assertion and does not handle cases comprehensively for erroneous configurations.

## Impact

Without proper validation, failure scenarios such as unexpected or invalid `ProviderData` values could lead to incorrect configuration behavior or errors subtly escaping proper reporting mechanisms. This is **medium** severity as it reduces the robustness of the code and leaves room for unexpected failure scenarios.

## Location

```go
func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
```

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

The validation can be improved to detect improperly configured `req.ProviderData`. Here's the updated code snippet:

```go
if req.ProviderData == nil {
    resp.Diagnostics.AddError(
        "Configure Error",
        "ProviderData is nil. Unable to configure the data source. This may be an unexpected issue.",
    )
    return
}

client, ok := req.ProviderData.(*api.ProviderClient)
if (!ok) {
    resp.Diagnostics.AddError(
        "ProviderData Type Mismatch",
        fmt.Sprintf("Expected *api.ProviderClient, found: %T. Ensure the provider has been properly initialized.", req.ProviderData),
    )
    return
}

// Proceed with configuration when this validation passes
d.TenantClient = NewTenantClient(client.Api)
```

This ensures that the code gracefully handles an improperly initialized `req.ProviderData` and provides useful diagnostics to aid end users and developers in troubleshooting the problem effectively.
