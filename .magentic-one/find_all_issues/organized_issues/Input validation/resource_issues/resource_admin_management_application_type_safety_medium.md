# Inadequate Handling of nil/Zero State in Read, Delete

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/resource_admin_management_application.go

## Problem

In both the `Read` and `Delete` methods, the code assumes that the state object is successfully populated by:

```go
resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
if resp.Diagnostics.HasError() {
    return
}
```

However, beyond checking for `Diagnostics.HasError()`, there is no validation that `state.Id` is actually set or that it passes UUID format requirements. This could enable runtime errors (nil dereference, bad API call) if the state is malformed or incomplete.

## Impact

Severity: **medium**

This could result in crashes or unpredictable calls to the API client, especially if Terraform’s state is externally modified, corrupted, or not correctly managed.

## Location

Read and Delete methods:

```go
var state AdminManagementApplicationResourceModel
resp.Diagnostics.Append(req.State.Get(ctx, &state)...) 
if resp.Diagnostics.HasError() { return }

// no further check of state.Id
adminApp, err := r.AdminManagementApplicationClient.GetAdminApplication(ctx, state.Id.ValueString())
err := r.AdminManagementApplicationClient.UnregisterAdminApplication(ctx, state.Id.ValueString())
```

## Fix

After retrieving state, validate that `state.Id` is set and is non-zero (and optionally validate UUID correctness):

```go
if state.Id.IsNull() || state.Id.ValueString() == "" {
    resp.Diagnostics.AddError("Missing or invalid ID in state", "Cannot perform operation: resource ID is not set or invalid.")
    return
}
```

This prevents awkward downstream errors and provides clearer diagnostics.
