# Inconsistent return value in Find function

##

/workspaces/terraform-provider-power-platform/internal/helpers/array/array.go

## Problem

The `Find` function returns the zero value of type `T` when no item is found that satisfies the predicate. This makes it impossible for callers to distinguish between a genuine match on the zero value and a case where nothing was found, especially for types like `int`, `string`, or other types where the zero value can be meaningful data.

## Impact

This could result in subtle bugs if the caller assumes a result indicates a real element was found. Severity: **medium**

## Location

```go
// Find returns the first element in the array that satisfies the predicate.
func Find[T comparable](arr []T, predicate func(T) bool) T {
	for _, v := range arr {
		if predicate(v) {
			return v
		}
	}
	var result T
	return result
}
```

## Fix

Return a second boolean value to indicate whether an element was found, following the Go idiomatic pattern.

```go
// Find returns the first element in the array that satisfies the predicate.
// The second return value is true if an element was found, false otherwise.
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
