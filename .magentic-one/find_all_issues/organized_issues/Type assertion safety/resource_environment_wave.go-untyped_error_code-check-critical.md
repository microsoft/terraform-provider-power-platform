# Untyped Error Code Check in Error Handling

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/resource_environment_wave.go

## Problem

The code checks if an error corresponds to `customerrors.ERROR_OBJECT_NOT_FOUND` by using the `customerrors.Code()` function, but does not use Go 1.13+ idiomatic error wrapping and type assertions for robust error handling. Using error codes as strings or constants can introduce bugs and makes handling error types less safe and readable.

## Impact

This can lead to fragile error handling logic if the error string or code changes. It is less type-safe and can hide error context, making troubleshooting more difficult. Severity: **critical** because it can cause logic to silently misbehave if the underlying error value changes or is wrapped.

## Location

```go
	if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
		resp.State.RemoveResource(ctx)
		return
	}
```

## Code Issue

```go
if customerrors.Code(err) == customerrors.ERROR_OBJECT_NOT_FOUND {
	resp.State.RemoveResource(ctx)
	return
}
```

## Fix

Use Go's standard library error wrapping and unwrapping with errors.Is or custom sentinel error variables for robust error comparison:

```go
import "errors"

if errors.Is(err, customerrors.ErrObjectNotFound) {
	resp.State.RemoveResource(ctx)
	return
}
```

**Explanation:**
- This approach uses type-safe error handling and supports Go's error-wrapping best practices (introduced in Go 1.13).
- Adjust `customerrors` to export a proper error variable if not already present (e.g., `var ErrObjectNotFound = errors.New("object not found")`).
- Benefits include maintainability, improved debugging, and correctness.