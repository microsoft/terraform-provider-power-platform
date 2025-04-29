# Title

Improper Type Assertion in Error Switch Block

##

`/workspaces/terraform-provider-power-platform/internal/api/client_test.go`

## Problem

In the function `TestUnitApiClient_GetConfig`, the `err` variable is being type-asserted to the `customerrors.UrlFormatError` type inside a `switch`. However, this assertion is unsafe, as `err` could be `nil`, leading to a potential unintended `default` case behavior.

For instance:
- If `err` is `nil`, `reflect.TypeOf(err)` would still execute, causing an unexpected output.
- Better defensiveness is needed when dealing with errors in test cases.

## Impact

- **Severity: Medium**
- Misleading error messages in the test results (`Expected error type ... but got ...`).
- Makes debugging more difficult, especially in cases of test failures.

## Location

First test function, in the body of `TestUnitApiClient_GetConfig()`.

```go
	switch err.(type) {
	case customerrors.UrlFormatError:
		return
	default:
		t.Errorf("Expected error type %s but got %s", reflect.TypeOf(customerrors.UrlFormatError{}), reflect.TypeOf(err))
	}
```

## Code Issue

```go
	switch err.(type) {
	case customerrors.UrlFormatError:
		return
	default:
		t.Errorf("Expected error type %s but got %s", reflect.TypeOf(customerrors.UrlFormatError{}), reflect.TypeOf(err))
	}
```

## Fix

Validate that `err` is not `nil` before attempting a type check. Alternatively, use the `errors.As()` method for safer type assertions in error chains.

```go
	if err == nil {
		t.Error("Expected an error but got nil")
		return
	}

	switch err.(type) {
	case customerrors.UrlFormatError:
		return
	default:
		t.Errorf("Expected error type %s but got %s", reflect.TypeOf(customerrors.UrlFormatError{}), reflect.TypeOf(err))
	}
```

Alternatively, using `errors.As()`:

```go
	var urlErr customerrors.UrlFormatError
	if !errors.As(err, &urlErr) {
		t.Errorf("Expected error type %T but got %T", urlErr, err)
	}
```
