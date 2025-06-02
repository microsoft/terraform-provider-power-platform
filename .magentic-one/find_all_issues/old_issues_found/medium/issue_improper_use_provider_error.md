### Issue 1

# Title

Improper use of `error` interface for `ProviderError` type

##

`/workspaces/terraform-provider-power-platform/internal/customerrors/provider_error.go`

## Problem

The `ProviderError` struct is being asserted as an `error` interface using `var _ error = ProviderError{}`. While syntactically correct, the `ProviderError` struct does not contain pointer semantics (e.g., it is a value receiver), which is generally not the best practice for `error` implementations. Doing so may lead to unintentional behavior if someone creates a pointer reference to `ProviderError` during use.

## Impact

- May create unexpected behavior while performing runtime type assertions.
- Loss of error context or misbehavior in error handling.
- Medium Severity: While not breaking immediately, it might introduce subtle bugs in complex applications.

## Location

Line defining `var _ error = ProviderError{}`.

## Code Issue

```go
var _ error = ProviderError{}
```

## Fix

Change the receiver functions of the `ProviderError` type to pointer semantics by using `*ProviderError` in the associated methods (`Error()`).

```go
func (e *ProviderError) Error() string {
	if e.Err == nil {
		return string(e.ErrorCode)
	}

	return fmt.Sprintf("%s: %s", e.ErrorCode, e.Err.Error())
}

func Unwrap(err error) error {
	if e, ok := err.(*ProviderError); ok {
		return errors.Unwrap(e.Err)
	}

	return errors.Unwrap(err)
}

func Code(err error) ErrorCode {
	if e, ok := err.(*ProviderError); ok {
		return e.ErrorCode
	}

	return ""
}
```