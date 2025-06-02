# Incomplete Mapping of Feature State Strings

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/resource_environment_wave.go

## Problem

The function `mapFeatureStateToSchemaState` provides an incomplete mapping of potential API states to the possible schema states. The function maps only two specific string values from the API (`Upgrading` to `upgrading`, `ON` to `enabled`), and any other value is mapped to `error`, including potentially valid but unknown or new values. This will result in all new/unknown states being treated as errors.

## Impact

If Microsoft adds states to the API, the provider will interpret all of them as `error`, making troubleshooting or feature evolution difficult. Severity: **medium** because it risks misrepresentation of provider state and can cause downstream automation or monitoring to misbehave.

## Location

```go
func mapFeatureStateToSchemaState(apiState string) string {
	switch apiState {
	case "Upgrading":
		return "upgrading"
	case "ON":
		return "enabled"
	default:
		return "error"
	}
}
```

## Code Issue

```go
func mapFeatureStateToSchemaState(apiState string) string {
	switch apiState {
	case "Upgrading":
		return "upgrading"
	case "ON":
		return "enabled"
	default:
		return "error"
	}
}
```

## Fix

Implement logging for unmapped/unknown values and/or a passthrough or explicit error state with diagnostics:

```go
func mapFeatureStateToSchemaState(apiState string) string {
	switch apiState {
	case "Upgrading":
		return "upgrading"
	case "ON":
		return "enabled"
	case "Error", "ERROR":
		return "error"
	default:
		tflog.Warn(context.TODO(), fmt.Sprintf("Unknown feature state from API: %s", apiState))
		return "error"
	}
}
```

- Or consider returning an explicit error or separate status for unmapped states
- Add documentation for maintainers that API state values might change