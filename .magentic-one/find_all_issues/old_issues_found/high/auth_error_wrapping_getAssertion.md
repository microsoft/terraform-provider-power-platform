# Title

Error handling can lead to potential information loss when wrapping errors in `getAssertion`.

##

/workspaces/terraform-provider-power-platform/internal/api/auth.go

## Problem

In the `getAssertion` method, errors are wrapped using `fmt.Errorf` without using `errors.As` or `errors.Unwrap`. This approach limits error transparency and makes it harder to inspect and troubleshoot root causes.

## Impact

This can lead to difficulties in debugging and troubleshooting, as the wrapped error loses context and becomes harder to unwrap. Severity: **high**.

## Location

The issue exists in multiple parts of the `getAssertion` method, particularly in the error handling.

## Code Issue

Example:

```go
	return "", fmt.Errorf("reading token file: %v", err)
```

## Fix

Adopt modern Go error wrapping using the `fmt.Errorf` `%w` verb for compatibility with `errors.Unwrap` or `errors.As`.

```go
	return "", fmt.Errorf("reading token file: %w", err)
```

Explanation:
- `%w` allows wrapped errors to maintain their original type, making them ser-compatible with error inspection methods like `errors.Is` or `errors.As`.
- This improves the ability for applications or developers to handle specific error cases without losing low-level error details.
