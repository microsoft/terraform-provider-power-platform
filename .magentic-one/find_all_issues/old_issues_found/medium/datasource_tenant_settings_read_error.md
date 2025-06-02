# Title

User-Friendly Error Messages in `Read` Function

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/datasource_tenant_settings.go`

## Problem

In the `Read` function, the error messages provided do not contain actionable information for end-users or developers. The current error message uses generic language and does not guide users on potential resolutions. Additionally, the error handling lacks logging, which could be useful for debugging purposes.

## Impact

This impacts the overall maintainability and end-user experience, as vague error messages make debugging more difficult. The severity is **medium** because it complicates error recovery and increases support costs but does not directly lead to code instability.

## Location

```go
func (d *TenantSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state TenantSettingsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantSettings, err := d.TenantSettingsClient.GetTenantSettings(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
		return
	}

	var configuredSettings TenantSettingsDataSourceModel
	req.Config.Get(ctx, &configuredSettings)
	state, _ = convertFromTenantSettingsDto[TenantSettingsDataSourceModel](*tenantSettings, state.Timeouts)
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
```

---

## Fix

Update error handling to include actionable messages and logging for better debug capability. Example fix:

```go
func (d *TenantSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state TenantSettingsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to retrieve current state.", map[string]interface{}{
			"state": state,
			"diagnostics": resp.Diagnostics,
		})
		return
	}

	tenantSettings, err := d.TenantSettingsClient.GetTenantSettings(ctx)
	if err != nil {
		errorMessage := fmt.Sprintf("Error occurred while reading tenant settings of type '%s'. Verify permissions and API connection settings.", d.FullTypeName())
		tflog.Error(ctx, errorMessage, map[string]interface{}{
			"error": err.Error(),
		})
		resp.Diagnostics.AddError("Failed to fetch tenant settings", errorMessage)
		return
	}

	var configuredSettings TenantSettingsDataSourceModel
	req.Config.Get(ctx, &configuredSettings)
	state, _ = convertFromTenantSettingsDto[TenantSettingsDataSourceModel](*tenantSettings, state.Timeouts)
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Error occurred while setting state.", map[string]interface{}{
			"state": state,
			"diagnostics": diags,
		})
		return
	}
}
```

This solution:
1. Logs the errors and includes relevant data for debugging.
2. Provides more actionable and user-friendly error messages.