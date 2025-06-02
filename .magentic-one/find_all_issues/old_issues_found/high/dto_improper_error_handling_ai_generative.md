# Title

Improper error handling for type casting in `convertAiGenerativeSettings`

##

/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/dto.go

## Problem

In the `convertAiGenerativeSettings` function, type assertions like `aiGenerativeSettingsObj.(basetypes.ObjectValue)` could cause runtime panics if the assertion fails. This is because the type assertion has not been accompanied by proper error handling.

## Impact

- Possible runtime panic in case of an unexpected type.
- Reduces robustness of the code.
- Harder to debug unexpected crashes.

Severity: High

## Location

Function `convertAiGenerativeSettings`.

## Code Issue

```go
if err := aiGenerativeSettingsObj.(basetypes.ObjectValue).As(ctx, &aiGenerativeSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); err != nil {
    return fmt.Errorf("failed to convert ai generative settings: %v", err)
}
```

## Fix

Introduce a two-step type assertion check to verify the type before accessing it, and provide meaningful error logs if the assertion fails.

```go
objectValue, ok := aiGenerativeSettingsObj.(basetypes.ObjectValue)
if !ok {
    return fmt.Errorf("expected basetypes.ObjectValue but got type %T", aiGenerativeSettingsObj)
}

if err := objectValue.As(ctx, &aiGenerativeSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true}); err != nil {
    return fmt.Errorf("failed to convert ai generative settings: %v", err)
}
```

This improves reliability and ensures that unexpected type mismatches are handled safely.
