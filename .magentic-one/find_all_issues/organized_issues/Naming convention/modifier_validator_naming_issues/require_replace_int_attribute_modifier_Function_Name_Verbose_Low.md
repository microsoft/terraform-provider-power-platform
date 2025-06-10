# Title

Function Name is Verbose and Not Self-Documenting

##

/workspaces/terraform-provider-power-platform/internal/modifiers/require_replace_int_attribute_modifier.go

## Problem

The function `RequireReplaceIntAttributePlanModifier` serves as a factory, but its name is lengthy and a bit awkward. While descriptive, its purpose and return contract aren't as clear as they could be. For maintainability, a well-placed docstring or a cleaner, project-aligned naming convention should be applied.

## Impact

Low: Does not affect functionality, but can make code less ergonomic and clear for users and maintainers.

## Location

Factory function near the top:

```go
func RequireReplaceIntAttributePlanModifier() planmodifier.Int64 {
    return &requireReplaceIntAttributePlanModifier{}
}
```

## Code Issue

```go
func RequireReplaceIntAttributePlanModifier() planmodifier.Int64 {
    return &requireReplaceIntAttributePlanModifier{}
}
```

## Fix

Add a short documentation string, and (if allowed by project conventions) consider shortening the function name:

```go
// RequireReplaceIntAttributePlanModifier returns a plan modifier that forces resource replacement when the int64 attribute changes from a non-zero, non-null, known state.
func RequireReplaceIntAttributePlanModifier() planmodifier.Int64 {
    return &requireReplaceIntAttributePlanModifier{}
}
```

Or, as an alternative:

```go
// NewReplaceInt64Modifier returns a plan modifier for int64 attributes that requires resource replacement on change from an existing value.
func NewReplaceInt64Modifier() planmodifier.Int64 {
    return &requireReplaceIntAttributePlanModifier{}
}
```
