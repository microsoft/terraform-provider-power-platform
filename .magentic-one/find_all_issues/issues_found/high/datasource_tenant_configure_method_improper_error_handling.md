# Title

Improper error handling in `Configure` method when `ProviderData` is `nil`

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant.go

## Problem

When `req.ProviderData` is `nil`, the `Configure` method of the data source simply returns without providing any diagnostics feedback or logging. This silent failure can confuse users or developers attempting to troubleshoot configuration issues.

## Impact

The lack of appropriate feedback or error messaging lowers code maintainability and could lead to debugging difficulties. Severity is **high**, as it impacts usability and the ability to diagnose issues effectively.

## Location

```go
func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse)
```

## Code Issue

```go
if req.ProviderData == nil {
    // ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
    return
}
```

## Fix

Provide proper diagnostics or logging when `req.ProviderData` is `nil`, to ensure developers or users can understand the cause of silent configuration failure:

```go
if req.ProviderData == nil {
    resp.Diagnostics.AddWarning(
        "ProviderData Missing",
        "The ProviderData is nil. Configuration might not proceed as expected. If this is unexpected, please check your setup.",
    )
    return
}
```

This fix improves visibility into system behavior and aids debugging efforts.
