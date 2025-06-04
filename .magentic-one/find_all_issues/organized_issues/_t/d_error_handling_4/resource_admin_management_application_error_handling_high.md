# Missing Diagnostic Error Reporting on Get Failure in Read Method

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/resource_admin_management_application.go

## Problem

In the `Read` method, if `r.AdminManagementApplicationClient.GetAdminApplication` returns an error, the error is simply ignored and the function returns without adding diagnostics or user-facing errors. This leads to silent failures where users or tools won't get meaningful feedback when a read operation fails (for example, if the resource is not found or a network error occurs).

## Impact

Severity: **high**

This makes debugging difficult for end-users and can hide underlying API or data issues. The lack of diagnostic reporting contradicts Terraform plugin best practices and leads to a poor user experience.

## Location

```go
func (r *AdminManagementApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	...
	adminApp, err := r.AdminManagementApplicationClient.GetAdminApplication(ctx, state.Id.ValueString())
	if err != nil {
		return // ISSUE: fails silently without diagnostics
	}
	...
}
```

## Fix

Add a user-visible error message to diagnostics when the Get operation fails:

```go
adminApp, err := r.AdminManagementApplicationClient.GetAdminApplication(ctx, state.Id.ValueString())
if err != nil {
    resp.Diagnostics.AddError("Failed to read admin management application", err.Error())
    return
}
```

This ensures that failures are communicated to the user and surface in Terraform diagnostics.
