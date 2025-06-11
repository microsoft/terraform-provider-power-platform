# Incorrect Use of Pointer in State Models (`*sourceModel`) for Terraform Framework

##

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/resource_enterprise_policy.go

## Problem

Throughout the resource methods (`Create`, `Read`, `Update`, `Delete`), the code uses `*sourceModel` pointers instead of the typical value type. The Terraform plugin framework expects value types rather than pointers to maintain compatibility, future-proofing, and to avoid nil dereferences or map assignment issues. Using pointers can cause unexpected behaviors when the framework marshals and unmarshals resource state.

## Impact

- Potential runtime panics (nil dereference on pointer).
- Framework integration issues; state updates may silently fail or behave unpredictably.
- Complicates refactoring and testing.
- **Severity:** Medium

## Location

Patterns in the following code:

```go
var plan *sourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
```
and

```go
var state *sourceModel
resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
```

## Code Issue

```go
var plan *sourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
```

## Fix

Use the value type, not a pointer. Change declarations and all related field access appropriately:

```go
var plan sourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
```

Same for `state` and any other state/plan model variable.

If you need to distinguish unset from set, use the Type's field accessors (e.g., `.IsNull()`/`.IsUnknown()`).

---

This markdown will be saved under:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/type_safety/resource_enterprise_policy_type_safety_medium.md`
