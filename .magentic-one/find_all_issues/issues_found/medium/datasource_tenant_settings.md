# Title

Improper Error Handling in `Configure` Function

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/datasource_tenant_settings.go`

## Problem

In the `Configure` function, the error handling for the `ProviderData` type check is generic and does not allow developers to distinguish between an expected nil `ProviderData` and a completely unexpected type. Additionally, there is no logging associated with this error scenario.

## Impact

The improper error handling makes it harder to debug issues when `ProviderData` is provided but is of an unexpected type. This impacts maintainability and troubleshooting, qualifying as a **medium severity** issue because it involves diagnostics and potential operational oversight.

## Location

```go
func (d *TenantSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	if req.ProviderData == nil { // ProviderData is expected to be null occasionally here
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
	d.TenantSettingsClient = newTenantSettingsClient(client.Api)
}
```

---

## Fix

Enhance diagnostics and add logging to provide more context when such an issue arises. Example fix:

```go
func (d *TenantSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil { 
		tflog.Debug(ctx, "ProviderData is nil during Configure. This is expected in context of ValidateConfig.")
		return
	}

	client, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		errorMessage := fmt.Sprintf("Expected *api.ProviderClient, got: %T.", req.ProviderData)
		tflog.Error(ctx, "Unexpected ProviderData Type", map[string]interface{}{
			"error": errorMessage,
		})
		resp.Diagnostics.AddError(
			"Unexpected ProviderData Type",
			fmt.Sprintf("%s Please report this issue to the provider developers.", errorMessage),
		)
		return
	}

	d.TenantSettingsClient = newTenantSettingsClient(client.Api)
}
```

This solution:
1. Logs a more verbose message for debugging.
2. Adds specific error context for developers.