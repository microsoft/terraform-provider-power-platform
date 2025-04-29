# Title

Unstructured error handling in `Create` method

## Path to the file

`/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/resource_admin_management_application.go`

## Problem

In the `Create` method, when an error occurs during `r.AdminManagementApplicationClient.RegisterAdminApplication` execution, the error is directly appended to `Diagnostics` with a generic error message. However, it doesn't provide detailed diagnostic context for easier debugging, nor does it log the full error.

## Impact

Limited error diagnostics and lack of logging make troubleshooting difficult, especially in production environments. This affects both developers and users who may need detailed error contexts while creating resources. Severity is **high** due to the significant impact on maintainability and debugging.

## Location

```go
adminApp, err := r.AdminManagementApplicationClient.RegisterAdminApplication(ctx, plan.Id.ValueString())
if err != nil {
    resp.Diagnostics.AddError("Failed to register admin application", fmt.Sprintf("Failed to register admin application: %v", err))
    return
}
```

## Fix

Provide detailed diagnostic context and implement proper logging for the error. This ensures that users and developers have sufficient information to address the issue.

### Suggested Fix:

```go
adminApp, err := r.AdminManagementApplicationClient.RegisterAdminApplication(ctx, plan.Id.ValueString())
if err != nil {
    diagnosticMsg := fmt.Sprintf("Failed to register admin application with ID %s: %v", plan.Id.ValueString(), err)
    resp.Diagnostics.AddError("Failed to register admin application", diagnosticMsg)
    tflog.Error(ctx, diagnosticMsg)
    return
}
```