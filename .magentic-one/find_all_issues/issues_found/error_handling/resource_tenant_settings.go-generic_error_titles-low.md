# Use of generic error messages in AddError/AddWarning

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go

## Problem

Some error and warning diagnostics are generic and repetitive, especially around "Error converting tenant settings", "Error reading tenant", "Error applying corrections", etc. There’s potential for more actionable, contextual messages for the end user.

## Impact

Reduced effectiveness of errors/warnings: less actionable for users and more difficult for provider maintainers to distinguish between error sources during support or debugging. Severity: low.

## Location

Throughout, e.g.:
- `resp.Diagnostics.AddError("Error converting tenant settings", err.Error())`
- `resp.Diagnostics.AddWarning("Tenant Settings are not deleted", ...)`, etc.

## Code Issue

```go
if err != nil {
	resp.Diagnostics.AddError("Error converting tenant settings", err.Error())
	return
}
```

## Fix

Provide context-rich error titles and messages, e.g.:

```go
if err != nil {
	resp.Diagnostics.AddError(
		"Unable to Convert Tenant Settings in resource_tenant_settings Create",
		fmt.Sprintf("Could not convert planned tenant settings model to DTO: %s", err.Error()),
	)
	return
}
```

Use similar context in AddWarning messages as well.

---

**This feedback will be saved under:**  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/resource_tenant_settings.go-generic_error_titles-low.md`
