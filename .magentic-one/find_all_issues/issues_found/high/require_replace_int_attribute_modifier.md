# Title

Missing validation for `req.PlanValue` in `PlanModifyInt64`

##

`/workspaces/terraform-provider-power-platform/internal/modifiers/require_replace_int_attribute_modifier.go`

## Problem

`req.PlanValue` is being compared directly with `req.StateValue`, but it is not being validated for null or unknown states before use. Directly comparing potentially invalid or unhandled values could result in unexpected behavior or runtime errors.

## Impact

If `req.PlanValue` is null or unknown and no checks are in place, the application may exhibit unpredictable behaviors. The severity of the issue is considered **high**, as it could lead to bugs in the functionality during plan modification in real-world scenarios.

## Location

Line in the following function:

```go
func (d *requireReplaceIntAttributePlanModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if req.PlanValue != req.StateValue && (!req.StateValue.IsNull() && !req.StateValue.IsUnknown() && req.StateValue.ValueInt64() != 0) {
		resp.RequiresReplace = true
	}
}
```

## Code Issue

```go
	if req.PlanValue != req.StateValue && (!req.StateValue.IsNull() && !req.StateValue.IsUnknown() && req.StateValue.ValueInt64() != 0) {
		resp.RequiresReplace = true
	}
```

## Fix

Add additional validation to ensure `req.PlanValue` is not null or unknown before comparing it with `req.StateValue`. This ensures all cases are handled, making the code resilient against invalid states.

```go
func (d *requireReplaceIntAttributePlanModifier) PlanModifyInt64(ctx context.Context, req planmodifier.Int64Request, resp *planmodifier.Int64Response) {
	if !req.PlanValue.IsNull() && !req.PlanValue.IsUnknown() && 
		req.PlanValue != req.StateValue && 
		(!req.StateValue.IsNull() && !req.StateValue.IsUnknown() && req.StateValue.ValueInt64() != 0) {
		resp.RequiresReplace = true
	}
}
```

This modification adds checks for `req.PlanValue` to ensure it is not null or unknown before comparison. This will prevent unexpected behavior and potential errors.
