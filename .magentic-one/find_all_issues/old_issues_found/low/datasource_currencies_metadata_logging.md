# Title

Lack of comprehensive logging in `Metadata` method

##

`/workspaces/terraform-provider-power-platform/internal/services/currencies/datasource_currencies.go`

## Problem

The logging in the `Metadata` method only provides minimal information using `tflog.Debug`. This lacks meaningful details that can assist in debugging, such as the exact content of `req.ProviderTypeName` or the eventual value of `resp.TypeName`.

## Impact

Diagnostic usefulness is low due to insufficient contextual details in log statements. This impacts operational tracing, but does not affect runtime functionality. Severity is **Low**.

## Location

The issue is in the `Metadata` method, at the point calling `tflog.Debug`.

## Code Issue

```go
func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    ...
    resp.TypeName = d.FullTypeName()
    tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}
```

## Fix

Enhance the log message to capture more contextual details. Example fix:

```go
func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    ...
    resp.TypeName = d.FullTypeName()
    tflog.Debug(ctx, "Setting metadata type name", map[string]interface{}{
        "ProviderTypeName": req.ProviderTypeName,
        "FullTypeName":     resp.TypeName,
    })
}
```