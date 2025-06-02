# Issue 3: Inefficient handling of generic types in `Find`

##

/workspaces/terraform-provider-power-platform/internal/helpers/array/array.go

## Problem

The `Find` function neglects to handle scenarios where no element in the array satisfies the predicate without exposing unnecessary ambiguity, as the generic type's zero value is returned. If the zero value is a legitimate instance of the type being searched for, it could lead to misinterpretation or subtle bugs.

## Impact

While generic zero values are convenient, they can lead to potential bugs or misinterpretation in critical sections of code, especially when developers assume the returned value corresponds to valid data. The accidental use of such zero values can result in logic failures.

**Severity:** High

## Location

File location: `Find` function in `/workspaces/terraform-provider-power-platform/internal/helpers/array/array.go`.

## Code Issue

Current code:

```go
	for _, v := range arr {
		if predicate(v) {
			return v
		}
	}
	var result T
	return result
```

## Fix

Update the function to return an additional boolean value indicating success/failure in finding a matching element. This enhances clarity and provides better handling for cases where no element satisfies the predicate.

```go
func Find[T comparable](arr []T, predicate func(T) bool) (T, bool) {
	for _, v := range arr {
		if predicate(v) {
			return v, true
		}
	}
	var result T
	return result, false
}
```

This modified implementation ensures the caller can accurately determine whether a match was found, minimizing the chance of logic errors.
