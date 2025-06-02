# Title

Ambiguous return value handling for validation and CRUD path

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user.go

## Problem

In "environment" user mode in both `Create` and `Update`, the code validates the list of security roles via `validateEnvironmentSecurityRoles`, but the validation is potentially run multiple times (once per invocation of the function, e.g., both in `Create` and again in `Update`). There's logic drift on how to ensure only valid roles are sent to API or in plan/state. Additionally, after some validation failures, code paths for CRUD operations still attempt to marshal or update user data, which should not occur for invalid input.

## Impact

This introduces potential for inconsistent runtime/resource states, repeated API calls for known-invalid operations, and blurs separation of validation versus mutation. It also increases cognitive load for maintainers. Severity: **Medium**.

## Location

Example from `Update`:

```go
err := validateEnvironmentSecurityRoles(plan.SecurityRoles)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
}
if len(addedSecurityRoles) > 0 {
    userDto, err := r.UserClient.AddEnvironmentUserSecurityRoles(ctx, plan.EnvironmentId.ValueString(), plan.AadId.ValueString(), addedSecurityRoles)
    ...
}
```

No explicit short-circuit/return after error, and logic continues as if validation passed.

## Code Issue

```go
err := validateEnvironmentSecurityRoles(plan.SecurityRoles)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
}
// further update/CRUD logic follows, potentially using invalid state
```

## Fix

Always return immediately after a validation error to prevent invalid state propagating to CRUD:

```go
err := validateEnvironmentSecurityRoles(plan.SecurityRoles)
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s", r.FullTypeName()), err.Error())
    return
}
// safe to continue with CRUD logic here
```

---

This issue should be saved in:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_user.go-validation_vs_crud-medium.md.
