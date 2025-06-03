# Title

Incorrect boolean logic for disabling delete in `Delete` method

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user.go

## Problem

In the `Delete` method, the conditional logic for handling deletion on Dataverse environments is inverted or misleading. The method currently checks for `if state.DisableDelete.ValueBool()` and then proceeds to delete only if it is true, which contradicts the intent and the documentation ("Disable delete" should prevent deletion).

## Impact

This could lead to actual deletion of Dataverse users when the "disable delete" flag is meant to prevent deletion—potentially causing unexpected and irreversible resource loss for end-users. Severity: **High** (risk of data/resource loss).

## Location

The problematic section:

```go
if hasEnvDataverse {
    if state.DisableDelete.ValueBool() { // this should be !state.DisableDelete.ValueBool()
        err := r.UserClient.DeleteDataverseUser(ctx, state.EnvironmentId.ValueString(), state.Id.ValueString())
        if err != nil {
            resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
            return
        }
    } else {
        tflog.Debug(ctx, fmt.Sprintf("Disable delete is set to false. Skipping delete of systemuser with id %s", state.Id.ValueString()))
    }
}
```

## Code Issue

```go
if state.DisableDelete.ValueBool() {
    err := r.UserClient.DeleteDataverseUser(ctx, state.EnvironmentId.ValueString(), state.Id.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
        return
    }
} else {
    tflog.Debug(ctx, fmt.Sprintf("Disable delete is set to false. Skipping delete of systemuser with id %s", state.Id.ValueString()))
}
```

## Fix

Invert the logic so the user is only deleted if `DisableDelete` is set to `false`:

```go
if hasEnvDataverse {
    if !state.DisableDelete.ValueBool() {
        err := r.UserClient.DeleteDataverseUser(ctx, state.EnvironmentId.ValueString(), state.Id.ValueString())
        if err != nil {
            resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
            return
        }
    } else {
        tflog.Debug(ctx, fmt.Sprintf("Disable delete is set to true. Skipping delete of systemuser with id %s", state.Id.ValueString()))
    }
}
```

---

This issue should be saved in:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/resource_user.go-disable_delete_logic-high.md.
