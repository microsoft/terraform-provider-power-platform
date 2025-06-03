# Return Value Inconsistency and Nil Return

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/api_environment_wave.go

## Problem

The `GetFeature` method returns `nil, nil` if the feature is not found. This can cause bugs if the caller does not check for a nil pointer before accessing the returned feature. Consider returning a well-defined error for "not found" cases.

## Impact

May lead to nil pointer dereference bugs later in the code. Severity: **medium**

## Location

In the `GetFeature` method:

```go
	return nil, nil
```

## Code Issue

```go
	return nil, nil
```

## Fix

Return a sentinel error to indicate not found:

```go
	return nil, fmt.Errorf("feature %s not found in environment %s", featureName, environmentId)
```

Or, if you want to keep the current structure, update all callers to handle the `nil, nil` return safely.
