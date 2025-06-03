# No check for nil `r.TenantSettingClient` in resource lifecycle methods

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go

## Problem

Resource lifecycle functions (`Create`, `Read`, `Update`, `Delete`, `ModifyPlan`) assume that `r.TenantSettingClient` is set without checking for nil after configuration. If for any reason `Configure` isn't called or fails to initialize `TenantSettingClient`, using this nil pointer will cause a panic at runtime.

## Impact

Unexpected panics at runtime if `Configure` method does not set up the client, causing a poor user experience and difficult debugging. Severity: high.

## Location

Any place in the resource's methods where `r.TenantSettingClient` is dereferenced, e.g.:

```go
originalSettings, erro := r.TenantSettingClient.GetTenantSettings(ctx)
```
(No nil-check in the first lines of the function.)

## Code Issue

```go
func (r *TenantSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
	defer exitContext()
	var plan TenantSettingsResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save the original tenant settings in private state
	originalSettings, erro := r.TenantSettingClient.GetTenantSettings(ctx) // <- nil pointer panic risk
	...
}
```

## Fix

At the start of each lifecycle method, check for a nil client and add a diagnostic error if so:

```go
if r.TenantSettingClient == nil {
	resp.Diagnostics.AddError(
		"Tenant Setting Client Not Configured",
		"Provider client was not properly configured. Please report this issue or check provider initialization.",
	)
	return
}
```

Add this block to start of `Create`, `Read`, `Update`, `Delete`, and (where relevant) `ModifyPlan`.

---

**This feedback will be saved under:**  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/resource_tenant_settings.go-nil_client_check-high.md`
