# Title

Improper Use of Error Variable in Custom Error Wrapping for Not Found Solution (GetSolutionUniqueName)

##

internal/services/solution/api_solution.go

## Problem

In the `GetSolutionUniqueName` method, if `len(solutions.Value) == 0`, the returned error is created by wrapping `err` (which is `nil` at this point) with a new error via `customerrors.WrapIntoProviderError`. Passing `nil` as the error argument can be misleading and is not idiomatic Go practice. The same pattern occurs in `GetSolutionById`.

## Impact

This can result in Go errors whose underlying cause is `nil`, reducing code clarity and making debugging harder. Severity: **medium**, because it impairs error traceability and could cause confusion in diagnostics.

## Location

```go
if len(solutions.Value) == 0 {
	return nil, customerrors.WrapIntoProviderError(err, customerrors.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("solution with unique name '%s' not found", name))
}
```

Similar:
```go
if len(solutions.Value) == 0 {
	return nil, customerrors.WrapIntoProviderError(err, customerrors.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("solution with id '%s' not found", solutionId))
}
```

## Code Issue

```go
if len(solutions.Value) == 0 {
	return nil, customerrors.WrapIntoProviderError(err, customerrors.ERROR_OBJECT_NOT_FOUND, fmt.Sprintf("solution with unique name '%s' not found", name))
}
```

## Fix

Create a new standard error to wrap, so the underlying error is meaningful.

```go
if len(solutions.Value) == 0 {
	baseErr := fmt.Errorf("solution with unique name '%s' not found", name)
	return nil, customerrors.WrapIntoProviderError(baseErr, customerrors.ERROR_OBJECT_NOT_FOUND, baseErr.Error())
}
```

And for GetSolutionById:
```go
if len(solutions.Value) == 0 {
	baseErr := fmt.Errorf("solution with id '%s' not found", solutionId)
	return nil, customerrors.WrapIntoProviderError(baseErr, customerrors.ERROR_OBJECT_NOT_FOUND, baseErr.Error())
}
```
