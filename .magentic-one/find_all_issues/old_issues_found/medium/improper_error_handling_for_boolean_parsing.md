# Title

Improper Error Handling for Boolean Parsing

##

/workspaces/terraform-provider-power-platform/internal/helpers/config.go

## Problem

The `GetConfigBool` function fails to implement proper error logging or recovery mechanisms when parsing the environment variable as a boolean (`strconv.ParseBool`). Though a warning is logged, it doesn't include potential details about remediation or call for corrective measures.

## Impact

This affects debugging and production readiness, as it might be unclear why and how boolean parsing failed, especially in critical components relying on environment variables. Severity: **Medium**

## Location

This issue arises in the `GetConfigBool` function.

## Code Issue

```go
	if err != nil {
		tflog.Warn(ctx, "Failed to parse environment variable value as a boolean. Using default value instead.", map[string]any{environmentVariableName: value})
		return defaultValue
	}
```

## Fix

Include more detailed logging (e.g., considering log levels or structured details) and an option to raise metrics or signals for downstream monitoring tools to identify such parsing failures immediately in production.

```go
	if err != nil {
		tflog.Error(ctx, "Failed to parse environment variable value as a boolean. Ensure the environment variable is a valid boolean (true/false). Defaulting to fallback value.", map[string]any{
			"environment_variable": environmentVariableName,
			"config_value":         value,
			"default_value_used":   defaultValue,
		})
		return defaultValue
	}
```