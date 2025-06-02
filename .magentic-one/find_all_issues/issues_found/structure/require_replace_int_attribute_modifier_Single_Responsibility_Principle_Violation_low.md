# Single Responsibility Principle Violation in Struct Declaration

##

/workspaces/terraform-provider-power-platform/internal/modifiers/require_replace_int_attribute_modifier.go

## Problem

The `requireReplaceIntAttributePlanModifier` struct is declared as an empty struct, which doesn't clarify its purpose or intent. Additionally, the implementation of its methods ties description, markdown description, and the core plan-modifying logic all to the same struct, violating Single Responsibility Principle and possibly confusing future maintainers.

## Impact

This is a **low** severity issue as it does not break functionality, but it does decrease code clarity and maintainability. A new developer may not immediately realize that this struct is solely a marker for the modifier, not holding any state.

## Location

Line where the struct is declared:

```go
type requireReplaceIntAttributePlanModifier struct {
}
```

## Code Issue

```go
type requireReplaceIntAttributePlanModifier struct {
}
```

## Fix

Add a comment explaining the reason for the empty struct or consider using `struct{}` directly if there's truly no internal state now or in the future. At minimum, provide a comment for maintainers:

```go
// requireReplaceIntAttributePlanModifier is a stateless plan modifier to force attribute replacement on change.
// It intentionally holds no state.
type requireReplaceIntAttributePlanModifier struct{}
```
