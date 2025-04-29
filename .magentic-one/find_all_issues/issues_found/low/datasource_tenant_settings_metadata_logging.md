# Title

Potential Misleading Logging in `Metadata` Function

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/datasource_tenant_settings.go`

## Problem

The `Metadata` function uses logging (`METADATA`) for debugging purposes but does not provide sufficient context about the operation or the data structure being processed. The log statement is generic and does not specify what metadata was successfully processed.

## Impact

This can lead to confusion during debugging as the log message does not provide detailed insights into the metadata operation. The severity is **low**, as it only impacts debugging and usability for developers.

## Location

```go
func (d *TenantSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    d.ProviderTypeName = req.ProviderTypeName

    ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
    defer exitContext()

    resp.TypeName = d.FullTypeName()
    tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}
```

---

## Fix

Refactor the log statement to provide detailed information about the metadata being processed. Example fix:

```go
func (d *TenantSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    d.ProviderTypeName = req.ProviderTypeName

    ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
    defer exitContext()

    resp.TypeName = d.FullTypeName()
    tflog.Debug(ctx, "Metadata processed successfully", map[string]interface{}{
        "ProviderTypeName": req.ProviderTypeName,
        "TypeName": resp.TypeName,
    })
}
```

This solution:
1. Adds structured logging with relevant fields like `ProviderTypeName` and `TypeName`.
2. Improves the debugging experience by providing more granular information.