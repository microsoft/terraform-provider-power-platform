# Repeated Code Patterns (DRY Violation) in Conversion Functions

## 
/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/dto.go

## Problem

The file contains several functions with highly similar logic patterns for DTO conversion and value extraction (e.g., `convertAiGenerativeSettings`, `convertAiGeneratedDesc`, `convertBackupRetention`, etc.). Each function repeats almost the same structure: fetch an attribute, check for null/unknown, convert type, prepare rule, append values. This is a classic violation of the DRY (Don't Repeat Yourself) principle.

## Impact

- **Severity:** Medium
- Increases maintenance cost and code size.
- More places for subtle bugs to arise.
- Difficult to ensure consistent logic/fixes when changing behavior.
- Hinders refactoring and readability.

## Location

Functions such as:

```go
func convertAiGenerativeSettings(ctx context.Context, attrs map[string]attr.Value, dto *EnvironmentGroupRuleSetValueSetDto) error { ... }
func convertAiGeneratedDesc(ctx context.Context, attrs map[string]attr.Value, dto *EnvironmentGroupRuleSetValueSetDto) error { ... }
func convertBackupRetention(ctx context.Context, attrs map[string]attr.Value, dto *EnvironmentGroupRuleSetValueSetDto) error { ... }
// ... many similar patterns for other attributes
```

## Code Issue

Each function starts with:

```go
obj := attrs["some_key"]
if !obj.IsNull() && !obj.IsUnknown() {
    var model someType
    if diags := obj.(basetypes.ObjectValue).As(ctx, &model, basetypes.ObjectAsOptions{...}); diags != nil {
        return fmt.Errorf("failed to convert %s: %v", "foo", diags)
    }
    // create rule, append to dto.Parameters
    // set rule.Value = append ...
}
return nil
```

## Fix

Consider genericizing or factoriziing the common extraction/checking pattern, for instance:

- Use a helper function for attribute extraction and type conversion.
- Use higher-order functions or closures if variations are small.
- If repetitive due to Go's typing, use code generation or templates to reduce the manual duplication.

Example pattern:

```go
func extractAndConvertAttr[T any](ctx context.Context, attrs map[string]attr.Value, key string, model *T) error {
    obj := attrs[key]
    if !obj.IsNull() && !obj.IsUnknown() {
        if diags := obj.(basetypes.ObjectValue).As(ctx, model, basetypes.ObjectAsOptions{...}); diags != nil {
            return fmt.Errorf("failed to convert %s: %v", key, diags)
        }
        return nil
    }
    return errors.New("attribute not found or is null/unknown")
}
```

You can then refactor each conversion function to call this utility, reducing boilerplate and improving maintainability.

---
