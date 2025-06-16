# Issue 1: Misleading Function Naming for GetConfigMultiString

##

/workspaces/terraform-provider-power-platform/internal/helpers/config.go

## Problem

The function `GetConfigMultiString` accepts an array of environment variable names, iterates through them, and returns the first non-empty value found. However, the function name does not clearly indicate that it returns a single string (the first set), not a concatenation or list of strings (which "MultiString" might imply).

## Impact

This can confuse future maintainers or users of this helper, making them think it returns multiple strings or a composite value. This is a low severity issue, but it affects code clarity and maintainability.

## Location

Line 25-39

## Code Issue

```go
func GetConfigMultiString(ctx context.Context, configValue basetypes.StringValue, environmentVariableNames []string, defaultValue string) string {
	if !configValue.IsNull() {
		return configValue.ValueString()
	}

	for _, k := range environmentVariableNames {
		if value, ok := os.LookupEnv(k); ok && value != "" {
			return value
		}
	}

	return defaultValue
}
```

## Fix

Consider renaming the function to reflect that it gets a single string from multiple environment sources, not multiple strings. For example: `GetConfigFirstAvailableString`.

```go
// GetConfigFirstAvailableString returns the value of configValue if set, otherwise the first set environment variable, else the defaultValue.
func GetConfigFirstAvailableString(ctx context.Context, configValue basetypes.StringValue, environmentVariableNames []string, defaultValue string) string {
	if !configValue.IsNull() {
		return configValue.ValueString()
	}

	for _, k := range environmentVariableNames {
		if value, ok := os.LookupEnv(k); ok && value != "" {
			return value
		}
	}

	return defaultValue
}
```
