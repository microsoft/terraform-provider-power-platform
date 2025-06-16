# Issue 3: Unnecessary Empty Struct Declaration

##

/workspaces/terraform-provider-power-platform/internal/modifiers/restore_original_value_modifier.go

## Problem

The struct `restoreOriginalValueModifier` is declared as an empty struct with no fields and used only as a receiver for interface implementation. While this is a common Go idiom, if no state is tracked, consider using a singleton or implement as a package-level variable to avoid unnecessary allocations.

## Impact

Minor performance and code clarity impact, as using instances of empty structs as function receivers is unnecessary if no state is managed. **Severity: low**

## Location

```go
type restoreOriginalValueModifier struct {
}
```

## Code Issue

```go
type restoreOriginalValueModifier struct {
}
```

## Fix

Use a single instance or a package-level variable if that improves clarity, or add a comment explaining its necessity. If nothing is stored and these are stateless, a singleton pattern is preferable.

```go
var restoreOriginalValueModifierInstance = &restoreOriginalValueModifier{}

func RestoreOriginalStringModifier() planmodifier.String {
    return restoreOriginalValueModifierInstance
}
func RestoreOriginalBoolModifier() planmodifier.Bool {
    return restoreOriginalValueModifierInstance
}
```

Or, if this idiom is required by the framework, add a comment to clarify.

