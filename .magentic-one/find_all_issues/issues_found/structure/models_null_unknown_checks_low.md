# Overly Complex and Repetitive Null/Unknown Checks

##

/workspaces/terraform-provider-power-platform/internal/services/environment/models.go

## Problem

Throughout the code, there are multiple verbose checks against `IsNull()` and `IsUnknown()` to determine if certain values (primarily from the `types` package) are valid. This leads to repetitive, error-prone code and makes it difficult to determine intent, especially when used in many locations for every optional attribute.

## Impact

- **Severity:** Low
- Makes the codebase more verbose and harder to maintain.
- May obscure which values are *required* vs. *optional*.
- Risk of subtle bugs if a check is missed or made incorrectly.

## Location

Widely used pattern such as:

```go
if !environmentSource.EnvironmentGroupId.IsNull() && !environmentSource.EnvironmentGroupId.IsUnknown() {
	environmentDto.Properties.ParentEnvironmentGroup = &ParentEnvironmentGroupDto{Id: environmentSource.EnvironmentGroupId.ValueString()}
}
...
if !environmentSource.AllowBingSearch.IsNull() && !environmentSource.AllowBingSearch.IsUnknown() {
	environmentDto.Properties.BingChatEnabled = environmentSource.AllowBingSearch.ValueBool()
}
```

And other similar repeated patterns.

## Code Issue

```go
if !value.IsNull() && !value.IsUnknown() {
	// Do something
}
```

## Fix

Refactor into helper functions to abstract the check, increasing readability and maintainability. For example:

```go
func isKnown(value basetypes.Value) bool {
	return !value.IsNull() && !value.IsUnknown()
}

// Usage:
if isKnown(environmentSource.AllowBingSearch) {
	environmentDto.Properties.BingChatEnabled = environmentSource.AllowBingSearch.ValueBool()
}
```

If such a utility function is already present in a helpers package, use it consistently everywhere.

---

**This markdown should be saved as:**  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/models_null_unknown_checks_low.md`
