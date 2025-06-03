# Use of Pointers for ResourceModel in Plan and State May Cause Nil Dereferences

##
/workspaces/terraform-provider-power-platform/internal/services/solution/resource_solution.go

## Problem
Throughout the code (e.g., in Create, Update, Read), the `plan` and `state` variables are declared as pointers (`*ResourceModel`), and populated via `req.Plan.Get()` or `req.State.Get()`, but code later dereferences fields without checking if the pointer is non-nil. If, for any reason, the plan or state is not properly assigned (e.g., decoding failure or data bug), this will cause a runtime panic with nil pointer dereference.

Example code:
```go
var plan *ResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...) // Could leave plan nil
...
if !plan.SettingsFile.IsNull() ... // possible nil dereference
```
Similar issue exists for `state`.

## Impact
- **Severity:** High
- Potential for runtime panics if Terraform sends an invalid or unexpected state or plan, leading to provider crashes.
- Makes the code unsafe to refactor or test.
- Reduces type safety and defensiveness, especially in a plugin context.

## Location
Functions like Create, Update, Read at lines resembling:
```go
var plan *ResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
// ...fields accessed via plan.<Field>
```

## Code Issue
```go
var plan *ResourceModel
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

// ...

if !plan.SettingsFile.IsNull() && !plan.SettingsFile.IsUnknown() {
    value, err := helpers.CalculateSHA256(plan.SettingsFile.ValueString())
// ...
```

## Fix
Immediately after extracting the plan or state, validate that it is non-nil before use, and emit a diagnostic (or handle gracefully) if it is nil:

```go
if plan == nil {
    resp.Diagnostics.AddError("Invalid plan received", "Resource plan is nil after decoding. This is likely an internal bug or provider incompatibility.")
    return
}

// Similar for 'state' as used in other functions
```
