# Unchecked Type Assertions can Panic

## 
/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/dto.go

## Problem

The code contains unchecked type assertions such as:

```go
aiGenerativeSettingsObj.(basetypes.ObjectValue).As(...)
```

If the value in `attrs["ai_generative_settings"]` is not of the expected type (`basetypes.ObjectValue`), this will cause a runtime panic, crashing the program. There are no checks (e.g., type assertion with second return value) to ensure the assertion succeeded before using the result.

## Impact

- **Severity:** High
- Unhandled panic paths can bring down the process (Terraform provider), possibly causing data loss or inconsistent state.
- Difficult to debug, as panics are abrupt and may not provide meaningful error messages to users or logs.
- Not idiomatic Go error handling.

## Location

Across all conversion functions. Example:

```go
if diags := aiGenerativeSettingsObj.(basetypes.ObjectValue).As(ctx, &aiGenerativeSettings, ...); diags != nil {
    ...
}
```

Similar patterns occur in other conversion helpers.

## Code Issue

```go
obj := attrs["ai_generative_settings"]
if !obj.IsNull() && !obj.IsUnknown() {
    var aiGenerativeSettings environmentGroupRuleSetAiGenerativeSettingsModel
    if diags := obj.(basetypes.ObjectValue).As(ctx, &aiGenerativeSettings, ...); diags != nil {
        ...
    }
}
```

## Fix

Use the "comma ok" type assertion form to ensure the type is correct before proceeding:

```go
obj := attrs["ai_generative_settings"]
if !obj.IsNull() && !obj.IsUnknown() {
    objVal, ok := obj.(basetypes.ObjectValue)
    if !ok {
        return fmt.Errorf("expected ai_generative_settings to be of type ObjectValue, got %T", obj)
    }
    var aiGenerativeSettings environmentGroupRuleSetAiGenerativeSettingsModel
    if diags := objVal.As(ctx, &aiGenerativeSettings, ...); diags != nil {
        ...
    }
}
```

Apply this change wherever type assertions are performed on `attr.Value` or similar interfaces.

---
