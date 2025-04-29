### Issue 2

# Title

Inadequate validation of `errorCode` in `NewProviderError()`

##

`/workspaces/terraform-provider-power-platform/internal/customerrors/provider_error.go`

## Problem

The `NewProviderError()` function does not validate the `ErrorCode` parameter before creating and returning a `ProviderError` instance. This could lead to invalid or unrecognized error codes being used throughout the codebase.

## Impact

- Makes maintaining and debugging the application harder.
- Enforceability and reliability of error codes are reduced.
- Medium Severity: Could propagate improper error codes throughout modules.

## Location

`NewProviderError` function definition.

## Code Issue

```go
func NewProviderError(errorCode ErrorCode, format string, args ...any) error {
	return ProviderError{
		Err:       fmt.Errorf(format, args...),
		ErrorCode: errorCode,
	}
}
```

## Fix

Validate `errorCode` against known constants before assigning it:

```go
func NewProviderError(errorCode ErrorCode, format string, args ...any) error {
	switch errorCode {
	case ERROR_OBJECT_NOT_FOUND,
		ERROR_ENVIRONMENT_URL_NOT_FOUND,
		ERROR_ENVIRONMENTS_IN_ENV_GROUP,
		ERROR_POLICY_ASSIGNED_TO_ENV_GROUP,
		ERROR_ENVIRONMENT_SETTINGS_FAILED,
		ERROR_ENVIRONMENT_CREATION:
		// Valid errorCode cases.
	default:
		return fmt.Errorf("invalid errorCode provided: %s", errorCode)
	}

	return ProviderError{
		Err:       fmt.Errorf(format, args...),
		ErrorCode: errorCode,
	}
}
```