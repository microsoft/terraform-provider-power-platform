# Error variable naming inconsistency (`erro`, `errj`) in Create

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go

## Problem

In the `Create` function, you use unconventional error variable names like `erro` and `errj`. For example:

- `originalSettings, erro := r.TenantSettingClient.GetTenantSettings(ctx)`
- `jsonSettings, errj := json.Marshal(originalSettings)`

This goes against Go naming conventions, which dictate simple, consistent use of `err` for error variables unless there is a strong necessity for disambiguation—typically only in more complex, nested scopes.

## Impact

Reduced code consistency and readability, especially for Go developers who expect the ubiquitous `err`. Severity: low.

## Location

Function `Create`.

## Code Issue

```go
originalSettings, erro := r.TenantSettingClient.GetTenantSettings(ctx)
if erro != nil {
	resp.Diagnostics.AddError(
		"Error reading tenant settings", fmt.Sprintf("Error reading tenant settings: %s", erro.Error()),
	)
	return
}

jsonSettings, errj := json.Marshal(originalSettings)
if errj != nil {
	resp.Diagnostics.AddError(
		"Error marshalling tenant settings", fmt.Sprintf("Error marshalling tenant settings: %s", errj.Error()),
	)
	return
}
```

## Fix

Use `err` for all error variables. Only use different names in special situations.

```go
originalSettings, err := r.TenantSettingClient.GetTenantSettings(ctx)
if err != nil {
	resp.Diagnostics.AddError(
		"Error reading tenant settings", fmt.Sprintf("Error reading tenant settings: %s", err.Error()),
	)
	return
}

jsonSettings, err := json.Marshal(originalSettings)
if err != nil {
	resp.Diagnostics.AddError(
		"Error marshalling tenant settings", fmt.Sprintf("Error marshalling tenant settings: %s", err.Error()),
	)
	return
}
```

---

**This feedback will be saved under:**  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/resource_tenant_settings.go-error_variable_naming-low.md`
