# Issue 2: Asymmetric Handling in PlanModifyBool

##

/workspaces/terraform-provider-power-platform/internal/modifiers/restore_original_value_modifier.go

## Problem

`PlanModifyBool` does not handle the resource destroy case (`if req.Plan.Raw.IsNull() {...}`), but `PlanModifyString` does. This is inconsistent and may lead to different behaviors for string and bool attributes regarding "restoring the original value" on destroy.

## Impact

This asymmetry could cause confusion or bugs, as the `PlanModifyBool` might not restore or log/handle the case where a resource is being destroyed, unlike its string counterpart. This is a **medium** severity maintainability and behavior consistency issue.

## Location

```go
func (d *restoreOriginalValueModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
    // Check if the resource is being created
    if req.State.Raw.IsNull() {
        if !req.ConfigValue.IsNull() {
            log.Default().Printf("Storing original value for attribute %s", req.PathExpression.String())
            resp.Private.SetKey(ctx, req.Path.String(), []byte{})
        }
    }
}
```

## Code Issue

```go
func (d *restoreOriginalValueModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
    // Check if the resource is being created
    if req.State.Raw.IsNull() {
        if !req.ConfigValue.IsNull() {
            log.Default().Printf("Storing original value for attribute %s", req.PathExpression.String())
            resp.Private.SetKey(ctx, req.Path.String(), []byte{})
        }
    }
}
```

## Fix

Add destroy handling logic to match what is done in `PlanModifyString`, or explain via comments why it is omitted.

```go
func (d *restoreOriginalValueModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
    // Check if the resource is being created
    if req.State.Raw.IsNull() {
        if !req.ConfigValue.IsNull() {
            // Remove log line
            resp.Private.SetKey(ctx, req.Path.String(), []byte{})
        }
    }
    // Optionally handle destroy case as in PlanModifyString
    if req.Plan.Raw.IsNull() {
        if !req.ConfigValue.IsNull() {
            // Restore logic or diagnostics if needed
            // e.g. resp.Private.GetKey(ctx, req.Path.String())
        }
    }
}
```
