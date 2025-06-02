### Issue 3

# Title

Uncontrolled nil-pointer dereference in `WrapIntoProviderError()`

##

`/workspaces/terraform-provider-power-platform/internal/customerrors/provider_error.go`

## Problem

The `WrapIntoProviderError()` function creates a `ProviderError` with a wrapped error. However, it does not sufficiently protect against all possible forms of invalid input. For instance, `err` could be `nil`, and such cases are not being handled adequately, leading to risky dereferencing or improper error chaining.

## Impact

- May cause potential runtime panics (e.g., nil pointer dereference).
- Impacts error clarity during debugging and propagation.
- High Severity: Causes runtime vulnerabilities in error handling logic.

## Location

`WrapIntoProviderError` function definition.

## Code Issue

```go
func WrapIntoProviderError(err error, errorCode ErrorCode, msg string) error {
	if err == nil {
		return ProviderError{
			Err:       fmt.Errorf("%s", msg),
			ErrorCode: errorCode,
		}
	}
	return ProviderError{
		Err:       fmt.Errorf("%s: [%w]", msg, err),
		ErrorCode: errorCode,
	}
}
```

## Fix

Add better checks and handling within the function:

```go
func WrapIntoProviderError(err error, errorCode ErrorCode, msg string) error {
	if errorCode == "" {
		return fmt.Errorf("errorCode cannot be empty in WrapIntoProviderError")
	}

	if err == nil {
		return ProviderError{
			Err:       fmt.Errorf("%s", msg),
			ErrorCode: errorCode,
		}
	}

	return ProviderError{
		Err:       fmt.Errorf("%s: [%w]", msg, err),
		ErrorCode: errorCode,
	}
}
```