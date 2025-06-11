# Title

Error Handling: Ineffective Error Double-Check in Delete Function

##

/workspaces/terraform-provider-power-platform/internal/services/environment_groups/resource_environment_group.go

## Problem

The `Delete` function checks `err` immediately after assignment:

```go
err := r.EnvironmentGroupClient.DeleteEnvironmentGroup(ctx, state.Id.ValueString())
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
    return
}
```

But then, later in the same function after deleting rulesets, it repeats the deletion:

```go
err = r.EnvironmentGroupClient.DeleteEnvironmentGroup(ctx, state.Id.ValueString())
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
    return
}
```

This can lead to confusion, multiple API calls, and potentially inconsistent error handling or resource state if the first `DeleteEnvironmentGroup` succeeded but later code failed and retried. Also, error codes (like `customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND`) are checked on the first error – but that code path should only run if that specific error happens.

## Impact

- May cause double-deletion API calls, potentially triggering unwanted provider behavior (high impact).
- Confuses control flow, increasing maintenance difficulty and risk of subtle bugs.
- Could leave resources in an inconsistent state.

**Severity:** high

## Location

Function: `Delete`

## Code Issue

```go
err := r.EnvironmentGroupClient.DeleteEnvironmentGroup(ctx, state.Id.ValueString())
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
    return
}

if customerrors.Code(err) == customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND || customerrors.Code(err) == customerrors.ERROR_POLICY_ASSIGNED_TO_ENV_GROUP {
    // cleanup logic...
    err = r.EnvironmentGroupClient.DeleteEnvironmentGroup(ctx, state.Id.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
        return
    }
}
```

## Fix

Refactor the control flow to handle the error code cases cleanly. For example:

```go
err := r.EnvironmentGroupClient.DeleteEnvironmentGroup(ctx, state.Id.ValueString())
if err == nil {
    return
}

code := customerrors.Code(err)
if code == customerrors.ERROR_ENVIRONMENT_URL_NOT_FOUND || code == customerrors.ERROR_POLICY_ASSIGNED_TO_ENV_GROUP {
    // cleanup logic...
    // (then retry deletion AFTER cleanup)
    err = r.EnvironmentGroupClient.DeleteEnvironmentGroup(ctx, state.Id.ValueString())
    if err != nil {
        resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
    }
    return
}

// All other errors
resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
```
