# Title

Ambiguity and potential misuse in plan/state assignment after `convertDataverseFromUserDto` call

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user.go

## Problem

When updating the Terraform plan or state after a CRUD operation, the code sets various fields from the result of `convertDataverseFromUserDto`. However, not all attributes are consistently assigned (sometimes state and plan are handled differently), and use of `req.Plan.SetAttribute` for `security_roles` can be confusing if plan merging is not needed.

## Impact

This could cause confusion in Terraform plan output, separate state from the model, or inadvertently drift between real resource and Terraform's model, especially if fields are optional or computed. Severity: **Medium**.

## Location

A snippet from `Create`:

```go
model := convertDataverseFromUserDto(&newUser, plan.DisableDelete.ValueBool())
plan.Id = model.Id
plan.AadId = model.AadId
req.Plan.SetAttribute(ctx, path.Root("security_roles"), model.SecurityRoles)
plan.UserPrincipalName = model.UserPrincipalName
plan.FirstName = model.FirstName
plan.LastName = model.LastName
plan.DisableDelete = model.DisableDelete
plan.BusinessUnitId = model.BusinessUnitId
```

## Code Issue

```go
plan.Id = model.Id
plan.AadId = model.AadId
req.Plan.SetAttribute(ctx, path.Root("security_roles"), model.SecurityRoles)
plan.UserPrincipalName = model.UserPrincipalName
plan.FirstName = model.FirstName
plan.LastName = model.LastName
plan.DisableDelete = model.DisableDelete
plan.BusinessUnitId = model.BusinessUnitId
```

## Fix

Assign values directly and consistently on the plan or state struct, ensuring uniformity and avoiding the mix of explicit assignment with `SetAttribute`. For example:

```go
plan.Id = model.Id
plan.AadId = model.AadId
plan.SecurityRoles = model.SecurityRoles  // Direct assignment
plan.UserPrincipalName = model.UserPrincipalName
plan.FirstName = model.FirstName
plan.LastName = model.LastName
plan.DisableDelete = model.DisableDelete
plan.BusinessUnitId = model.BusinessUnitId
```

Only use `SetAttribute` when explicit plan merging or non-struct fields are required.

---

This issue should be saved in:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/resource_user.go-plan_assignment-medium.md.
