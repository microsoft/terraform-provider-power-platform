# Issue 5: Unused Context Parameter in Several Functions

##

/workspaces/terraform-provider-power-platform/internal/helpers/config.go

## Problem

Both `GetConfigString` and `GetConfigMultiString` receive a `ctx context.Context` parameter but do not use it. This can be misleading to function callers, suggesting context-related behavior when there is none.

## Impact

Low severity, but reduces clarity/maintainability and might confuse users as to context utility.

## Location

Lines 13-39

## Code Issue

```go
func GetConfigString(ctx context.Context, configValue basetypes.StringValue, environmentVariableName string, defaultValue string) string {
	// ctx is unused
	if !configValue.IsNull() {
		return configValue.ValueString()
	} else if value, ok := os.LookupEnv(environmentVariableName); ok && value != "" {
		return value
	}
	return defaultValue
}

// Similar in GetConfigMultiString
```

## Fix

Remove the unused parameter if not needed, or use it if logging/side-effects may be added in the future.

```go
func GetConfigString(configValue basetypes.StringValue, environmentVariableName string, defaultValue string) string {
	if !configValue.IsNull() {
		return configValue.ValueString()
	} else if value, ok := os.LookupEnv(environmentVariableName); ok && value != "" {
		return value
	}
	return defaultValue
}
```
