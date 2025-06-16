# No state refresh after Delete if restore operation fails

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go

## Problem

In the `Delete` method, if the attempt to restore settings (i.e., "delete" this resource by reverting the settings) fails, no attempt is made to refresh or clean up the state, nor is the state removed from Terraform. This could leave users with orphaned or inconsistent Terraform state.

## Impact

Terraform state may no longer accurately represent the underlying reality. On repeated apply/destroy operations, users may encounter confusing or inconsistent provider errors. Severity: medium.

## Location

Within `Delete` function, after this code:

```go
_, e := r.TenantSettingClient.UpdateTenantSettings(ctx, *correctedDto)
if e != nil {
	resp.Diagnostics.AddError(
		"Error deleting tenant settings", fmt.Sprintf("Error deleting tenant settings: %s", e.Error()),
	)
	return
}
```

## Code Issue

```go
if e != nil {
	resp.Diagnostics.AddError(
		"Error deleting tenant settings", fmt.Sprintf("Error deleting tenant settings: %s", e.Error()),
	)
	return
}
```
No `resp.State.RemoveResource(ctx)` or diagnostic/attempt to clear state.

## Fix

After the error diagnostic, consider also removing the resource from state if the restore failed:

```go
if e != nil {
	resp.Diagnostics.AddError(
		"Error deleting tenant settings", fmt.Sprintf("Error restoring (deleting) tenant settings: %s", e.Error()),
	)
	resp.State.RemoveResource(ctx)
	return
}
```

Alternatively, document the impact for users.

---

**This feedback will be saved under:**  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/resource_tenant_settings.go-delete_state_refresh-medium.md`
