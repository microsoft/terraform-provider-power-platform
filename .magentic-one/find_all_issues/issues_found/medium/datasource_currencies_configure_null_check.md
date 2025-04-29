# Title

Insufficient null check for `ProviderData` in `Configure` method

##

`/workspaces/terraform-provider-power-platform/internal/services/currencies/datasource_currencies.go`

## Problem

In the `Configure` method, the `ProviderData` is checked for `nil`, but there is no robust validation or additional logging to handle cases where the data may be malformed or contain unexpected values.

## Impact

This may lead to unexpected behavior or crashes if `ProviderData` holds invalid or corrupted entries. Severity is **Medium** since it could introduce a reliability issue.

## Location

The issue is within the `Configure` method, near the `if req.ProviderData == nil` block.

## Code Issue

```go
func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    ...
    if req.ProviderData == nil {
        // ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
        return
    }
    client, ok := req.ProviderData.(*api.ProviderClient)
    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected ProviderData Type",
            fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
        )
        return
    }
    ...
}
```

## Fix

Add more robust logging or handling to ensure unexpected values are identifiable during debugging.

```go
func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    ...
    if req.ProviderData == nil {
        tflog.Warn(ctx, "ProviderData is nil. This is expected during ValidateConfig.")
        return
    }
    client, ok := req.ProviderData.(*api.ProviderClient)
    if !ok {
        resp.Diagnostics.AddError(
            "Unexpected ProviderData Type",
            fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please check the initialization logic in ConfigureRequest.", req.ProviderData),
        )
        return
    }
    d.CurrenciesClient = newCurrenciesClient(client.Api)
    ...
}
```