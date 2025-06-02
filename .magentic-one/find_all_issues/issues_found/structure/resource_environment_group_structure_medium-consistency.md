# Title

Multiple Use of Pointers for ResourceModel in CRUD Functions Reduces Consistency and Clarity

##

/workspaces/terraform-provider-power-platform/internal/services/environment_groups/resource_environment_group.go

## Problem

In your CRUD operations, you sometimes use a value (`EnvironmentGroupResourceModel{}` in `Read`), and sometimes a pointer (`var plan *EnvironmentGroupResourceModel` in `Create`, `Update`). This inconsistency can cause unexpected nil pointer dereference panics or require unnecessary allocation and pointer indirection.

## Impact

- Potential for panic if the pointer is not properly set by the framework (medium impact).
- Reduces code readability and consistency.
- May complicate mocking and testing.

**Severity:** medium

## Location

Functions: `Read`, `Create`, `Update`, `Delete`

## Code Issue

Current pattern in `Create`:

```go
var plan *EnvironmentGroupResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
    return
}

// usage
environmentGroupToCreate := environmentGroupDto{
    DisplayName: plan.DisplayName.ValueString(),
    Description: plan.Description.ValueString(),
}
```

Pattern in `Read`:

```go
state := EnvironmentGroupResourceModel{}
resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
```

## Fix

Prefer consistent use—if the struct is not very large and ownership does not need to be transferred, use value types:

```go
plan := EnvironmentGroupResourceModel{}
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
    return
}
...
// use plan.X not plan->X
```

Apply this pattern across CRUD functions for clarity and safety.
