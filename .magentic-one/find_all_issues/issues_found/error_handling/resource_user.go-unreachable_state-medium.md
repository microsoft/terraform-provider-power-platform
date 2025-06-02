# Title

Unreachable code or inconsistent state after removing all security roles in Read (environment user)

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user.go

## Problem

In the `Read` method (when handling an environment user), if all security roles are removed, a special case triggers assignment of `user.AadObjectId` and `user.DomainName` based on state, but the code then continues with marshalling and converting the (potentially incomplete) `user`.

## Impact

This can result in setting state with partially populated user data, leading to data drift between Terraform and the actual environment, possibly breaking subsequent operations or plan consistency. Severity: **Medium**.

## Location

From `Read` in "environment" branch:

```go
user, err := r.UserClient.GetEnvironmentUserByAadObjectId(ctx, state.EnvironmentId.ValueString(), state.AadId.ValueString())
// if all the security roles are removed, the user will not be found
if user.AadObjectId == "" {
	user.AadObjectId = state.AadId.ValueString()
	user.DomainName = state.UserPrincipalName.ValueString()
}
if err != nil {
...
}
rolesBytes, err := json.Marshal(user.SecurityRoles)
...
updateUser = *user
```

## Code Issue

```go
if user.AadObjectId == "" {
	user.AadObjectId = state.AadId.ValueString()
	user.DomainName = state.UserPrincipalName.ValueString()
}
// (error handling for user == nil and err follows, possibly too late)
```

## Fix

Move error handling and nil checks before attempting to operate on the `user` pointer.  
If the result indicates a removed user (all roles removed / not found), return early and clear state:

```go
user, err := r.UserClient.GetEnvironmentUserByAadObjectId(ctx, state.EnvironmentId.ValueString(), state.AadId.ValueString())
if err != nil {
    if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
        resp.State.RemoveResource(ctx)
        return
    }
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", r.FullTypeName()), err.Error())
    return
}

if user == nil || user.AadObjectId == "" {
    // User not present, remove from state and exit
    resp.State.RemoveResource(ctx)
    return
}

// Only now safe to use user safely
rolesBytes, err := json.Marshal(user.SecurityRoles)
// ... etc
```

---

This issue should be saved in:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/resource_user.go-unreachable_state-medium.md.
