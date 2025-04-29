# Title

Handling error scenarios in `Read` method

## Path to the file

`/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/resource_admin_management_application.go`

## Problem

The `Read` method directly returns when an error occurs in the `r.AdminManagementApplicationClient.GetAdminApplication` call, but does not append diagnostic information or log the error. This reduces observability and makes debugging difficult.

## Impact

Without logging or diagnostics for errors, it's harder for developers or users of the Terraform provider to identify the specific reason for errors during the `Read` operation. Severity level is **high**, as lack of error handling and reporting significantly impacts usability and maintainability.

## Location

```go
adminApp, err := r.AdminManagementApplicationClient.GetAdminApplication(ctx, state.Id.ValueString())
if err != nil {
    return
}
```

## Fix

Append diagnostic information and log the error to improve observability and troubleshooting. For example:

### Suggested Fix:

```go
adminApp, err := r.AdminManagementApplicationClient.GetAdminApplication(ctx, state.Id.ValueString())
if err != nil {
    resp.Diagnostics.AddError(
        "Failed to fetch admin application",
        fmt.Sprintf("Error fetching admin application with ID %s: %v", state.Id.ValueString(), err),
    )
    tflog.Error(ctx, fmt.Sprintf("Read failed with error: %v", err))
    return
}
```
