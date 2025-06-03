# Title

Potential nil pointer dereference in `NewUUIDPointerValueMust`

##

/workspaces/terraform-provider-power-platform/internal/customtypes/uuid.go

## Problem

The function `NewUUIDPointerValueMust` dereferences the `value` pointer without checking if it is `nil`. If `value` is `nil`, this will cause a runtime panic due to dereferencing a nil pointer.

## Impact

This is a high-severity error handling and control flow issue. Dereferencing a nil pointer can cause the application to crash at runtime, leading to instability and potentially bringing down important processes or resources.

## Location

```go
func NewUUIDPointerValueMust(value *string) (UUID, diag.Diagnostics) {
	return NewUUIDValue(*value).ValueUUID()
}
```

## Code Issue

```go
func NewUUIDPointerValueMust(value *string) (UUID, diag.Diagnostics) {
	return NewUUIDValue(*value).ValueUUID()
}
```

## Fix

Always check for `nil` before dereferencing a pointer. You can mimic the pattern in `NewUUIDPointerValue` to handle `nil` values safely.

```go
func NewUUIDPointerValueMust(value *string) (UUID, diag.Diagnostics) {
	if value == nil {
		return NewUUIDNull().ValueUUID()
	}
	return NewUUIDValue(*value).ValueUUID()
}
```
