# Function Naming: Non-Idiomatic Name ‘RequireReplaceStringFromNonEmptyPlanModifier’

##

/workspaces/terraform-provider-power-platform/internal/modifiers/requires_replace_string_from_non_empty_modifier.go

## Problem

The function name `RequireReplaceStringFromNonEmptyPlanModifier` does not follow Go idioms (like PascalCase for exported functions), and is excessively verbose and difficult to quickly parse. Good function names should be concise but meaningful and generally avoid repeating information unnecessarily. This can make function usages more readable and maintainable.

## Impact

**Severity: Low**

Verbose or unclear naming can reduce readability and make the code slightly harder to maintain or consume, especially as your codebase or team grows. Go recommends being brief but accurate in naming, relying on package context.

## Location

```go
func RequireReplaceStringFromNonEmptyPlanModifier() planmodifier.String {
	return &requireReplaceStringFromNonEmptyPlanModifier{}
}
```

## Code Issue

```go
func RequireReplaceStringFromNonEmptyPlanModifier() planmodifier.String {
	return &requireReplaceStringFromNonEmptyPlanModifier{}
}
```

## Fix

Rename the function and corresponding struct to follow Go best practices. For example:

```go
func ReplaceOnNonEmptyStringChange() planmodifier.String {
	return &replaceOnNonEmptyStringChange{}
}
```

This aligns with Go naming conventions (concise and descriptive), and pairs well with revised struct naming. Update usages in the project as needed.
