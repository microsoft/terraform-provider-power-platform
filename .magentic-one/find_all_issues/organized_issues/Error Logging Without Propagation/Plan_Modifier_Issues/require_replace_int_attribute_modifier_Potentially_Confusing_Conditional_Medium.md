# Title

Potentially Confusing Conditional in PlanModifyInt64

##

/workspaces/terraform-provider-power-platform/internal/modifiers/require_replace_int_attribute_modifier.go

## Problem

The conditional in `PlanModifyInt64` is currently implemented as a single dense line:

```go
if req.PlanValue != req.StateValue && (!req.StateValue.IsNull() && !req.StateValue.IsUnknown() && req.StateValue.ValueInt64() != 0) {
    resp.RequiresReplace = true
}
```

This logic is not self-explanatory and could easily cause confusion or lead to errors if maintained in the future without clear documentation or refactoring. It can be made more readable by splitting the condition into well-named variables and adding proper documentation.

## Impact

Medium: Reduced code readability and maintainability, as future changes to the condition might introduce mistakes. It's also difficult to debug or audit the business logic.

## Location

Within the method `PlanModifyInt64`:

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

Refactor the conditional to use well-named intermediate variables, improving readability and making logic modification safer.

```go
isValueChanged := req.PlanValue != req.StateValue
hasValidPreviousState := !req.StateValue.IsNull() && !req.StateValue.IsUnknown() && req.StateValue.ValueInt64() != 0

if isValueChanged && hasValidPreviousState {
    resp.RequiresReplace = true
}
```

Add a comment to document why this logic is in place to prevent confusion for future maintainers.
