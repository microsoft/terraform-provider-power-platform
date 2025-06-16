# Lack of Blank Lines Between Major Function Definitions

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/datasource_tenant_settings.go

## Problem

Top-level exported methods in this file are not all separated by a blank line. In Go, it is a common style convention to separate function/method definitions by at least one blank line for better readability and code navigation. The lack of visual separation can make the file hard to scan, especially as it grows or if functions themselves are large.

## Impact

Reduces code readability and slows down developers trying to navigate the codebase. **Severity: Low**

## Location

E.g. (snippet):
```go
func NewTenantSettingsDataSource() datasource.DataSource {
    ...
}
func (d *TenantSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    ...
}
```

## Fix

Insert a blank line between each top-level function/method definition, as in:

```go
func NewTenantSettingsDataSource() datasource.DataSource {
    ...
}

func (d *TenantSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    ...
}
```
