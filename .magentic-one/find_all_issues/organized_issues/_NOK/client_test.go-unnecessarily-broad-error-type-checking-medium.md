# Unnecessarily Broad Error Type Checking

##

/workspaces/terraform-provider-power-platform/internal/api/client_test.go

## Problem

In `TestUnitApiClient_GetConfig`, the test uses `switch err.(type)` to check for the custom error `customerrors.UrlFormatError`, but it does not use a type assertion or `errors.As` to confirm that the `err` object actually wraps the target type. In Go, error wrapping is common, and simply checking the type with a type switch may miss wrapped errors—if the actual error is wrapped, it will fail the type switch but would pass an `errors.As` check.

## Impact

Medium severity. Test may incorrectly fail if the error is wrapped (e.g., with `fmt.Errorf("...: %w", err)`). Best practices dictate using `errors.As` for custom error type checks to handle wrapped errors. False negatives in tests can hinder future refactoring and reliability of the codebase.

## Location

```go
switch err.(type) {
case customerrors.UrlFormatError:
	return
default:
	t.Errorf("Expected error type %s but got %s", reflect.TypeOf(customerrors.UrlFormatError{}), reflect.TypeOf(err))
}
```

## Fix

Use the standard library's `errors.As()` function to check for (possibly wrapped) error types.

```go
var urlErr customerrors.UrlFormatError
if errors.As(err, &urlErr) {
	return
} else {
	t.Errorf("Expected error type %s but got %s", reflect.TypeOf(customerrors.UrlFormatError{}), reflect.TypeOf(err))
}
```
