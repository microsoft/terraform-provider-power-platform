# String Comparison on Error Values in Tests

##

/workspaces/terraform-provider-power-platform/internal/api/client_test.go

## Problem

Several tests compare error messages via string matching (`err.Error() == ...` or `strings.HasPrefix(err.Error(), ...)`) to validate that the correct error occurred. This approach is brittle: if the error messages change due to upstream library changes, refactoring, or localization, the tests may start failing even if the error type is unchanged and correct. It's more robust to check for custom error types (using `errors.Is` or `errors.As`), or at least well-defined sentinel error values.

## Impact

Medium severity. Can cause fragile tests, false negatives, and unnecessary test maintenance.

## Location

Multiple places, for example:

```go
if err.Error() != "context deadline exceeded" {
	t.Errorf("Expected error message %s but got %s", "context deadline exceeded", err.Error())
}
```

```go
if !strings.HasPrefix(err.Error(), expectedError) {
	t.Errorf("Expected error message '%s' but got '%s'", expectedError, err.Error())
}
```

## Code Issue

Example:
```go
if err.Error() != "context deadline exceeded" {
	t.Errorf("Expected error message %s but got %s", "context deadline exceeded", err.Error())
}
```

Example:
```go
if !strings.HasPrefix(err.Error(), expectedError) {
	t.Errorf("Expected error message '%s' but got '%s'", expectedError, err.Error())
}
```

## Fix

Where possible, compare error types or use `errors.Is` against sentinel error values. For context cancellation, use `errors.Is(err, context.DeadlineExceeded)`:

```go
if !errors.Is(err, context.DeadlineExceeded) {
	t.Errorf("Expected error type context.DeadlineExceeded but got %v", err)
}
```

For domain-specific errors, prefer returning sentinel errors or error types that can be checked via `errors.Is` or `errors.As`, rather than matching error strings.
