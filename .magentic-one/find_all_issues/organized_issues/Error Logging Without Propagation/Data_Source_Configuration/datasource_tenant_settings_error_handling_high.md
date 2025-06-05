# Lack of Error Handling After Config Get in Read Method

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/datasource_tenant_settings.go

## Problem

Within the `Read` method of `TenantSettingsDataSource`, the result of `req.Config.Get(ctx, &configuredSettings)` is not checked for diagnostics/errors. This violates robust error handling principles, since if the configuration cannot be decoded (due to, for instance, a schema mismatch or context error), the method continues execution, potentially leading to incorrect logic or even panics downstream.

## Impact

Failing to handle diagnostics can result in unexpected behavior, invalid state propagation, and complicates troubleshooting for users. **Severity: High**

## Location

Line inside `Read` method:
```go
var configuredSettings TenantSettingsDataSourceModel
req.Config.Get(ctx, &configuredSettings)
state, _, err = convertFromTenantSettingsDto[TenantSettingsDataSourceModel](*tenantSettings, state.Timeouts)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Error converting tenant settings: %s", d.FullTypeName()), err.Error())
    return
}
```

## Fix

Check and append diagnostics immediately after calling `req.Config.Get`, and abort execution if any errors are present, consistent with prior error checks in this method:

```go
var configuredSettings TenantSettingsDataSourceModel
diags := req.Config.Get(ctx, &configuredSettings)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
    return
}
state, _, err = convertFromTenantSettingsDto[TenantSettingsDataSourceModel](*tenantSettings, state.Timeouts)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Error converting tenant settings: %s", d.FullTypeName()), err.Error())
    return
}
```

This ensures proper error handling, improves code robustness, and aligns with best practices used elsewhere in the provider.
