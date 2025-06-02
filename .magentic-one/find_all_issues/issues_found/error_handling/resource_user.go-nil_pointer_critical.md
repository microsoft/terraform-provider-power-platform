# Title

Potential nil pointer dereference when unwrapping userDto in environment path

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user.go

## Problem

In several code paths (notably in else branches for "environment" users in `Create`, `Read`, and `Update`), the code dereferences pointers returned from helper functions without nil checks, e.g. `user, err := r.UserClient.CreateEnvironmentUser(...)` followed immediately by `newUser = *user`. If the implementation of the helper methods (`CreateEnvironmentUser`, `GetEnvironmentUserByAadObjectId`, etc.) can return a `nil` pointer on error, this will result in a runtime panic.

## Impact

If the underlying API returns nil user objects, this will cause a panic, crashing the provider and bringing down the whole Terraform operation. Severity: **Critical** (reliability and crash safety concern).

## Location

Multiple locations, such as (from `Create`):

```go
user, err := r.UserClient.CreateEnvironmentUser(ctx, plan.EnvironmentId.ValueString(), plan.AadId.ValueString(), plan.SecurityRoles)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
    return
}

// No nil check on user
rolesBytes, err := json.Marshal(user.SecurityRoles)
```

And then:

```go
newUser = *user
```

## Code Issue

```go
user, err := r.UserClient.CreateEnvironmentUser(ctx, plan.EnvironmentId.ValueString(), plan.AadId.ValueString(), plan.SecurityRoles)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
    return
}

// dereferencing user without checking for nil
rolesBytes, err := json.Marshal(user.SecurityRoles)
if err != nil {
    ...
}
resp.Private.SetKey(ctx, "role", rolesBytes)

newUser = *user
```

## Fix

Before dereferencing the user pointer, always check for nil and handle accordingly (add a diagnostic if user is nil):

```go
user, err := r.UserClient.CreateEnvironmentUser(ctx, plan.EnvironmentId.ValueString(), plan.AadId.ValueString(), plan.SecurityRoles)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
    return
}

if user == nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Unexpected nil user returned when creating %s", r.FullTypeName()), "API returned nil user object")
    return
}

// Now safe to use
rolesBytes, err := json.Marshal(user.SecurityRoles)
// ...
newUser = *user
```

Perform similar nil-checks in `Read`, `Update`, etc. wherever API helper returns a pointer to avoid panics.

---

This issue should be saved in:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/resource_user.go-nil_pointer_critical.md.
