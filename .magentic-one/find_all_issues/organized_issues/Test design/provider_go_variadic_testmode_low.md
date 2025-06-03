# Title

Deprecated or Unidiomatic Usage: Variadic testModeEnabled Parameter

##

/workspaces/terraform-provider-power-platform/internal/provider/provider.go

## Problem

In the `NewPowerPlatformProvider` constructor, the `testModeEnabled ...bool` parameter uses a variadic form but is only ever checked for `[0]`. This creates a misleading function signature, as variadic suggests that multiple values are relevant, but only the first one is used.

## Impact

This usage can cause confusion for developers and users of the function, is non-idiomatic, and does not accurately reflect the function's intent. Severity: **low**.

## Location

```go
func NewPowerPlatformProvider(ctx context.Context, testModeEnabled ...bool) func() provider.Provider {
   ...
   if len(testModeEnabled) > 0 && testModeEnabled[0] {
      ...
      providerConfig.TestMode = true
   }
   ...
}
```

## Fix

Change to a standard boolean parameter:

```go
func NewPowerPlatformProvider(ctx context.Context, testModeEnabled bool) func() provider.Provider {
    ...
    if testModeEnabled {
        ...
        providerConfig.TestMode = true
    }
    ...
}
```

Update all call sites to pass an explicit boolean argument, which better communicates intent and allows for stronger type checking.
