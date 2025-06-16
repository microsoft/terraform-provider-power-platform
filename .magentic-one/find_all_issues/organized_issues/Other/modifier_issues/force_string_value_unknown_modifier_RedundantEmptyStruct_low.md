# Redundant empty struct definition

##
/workspaces/terraform-provider-power-platform/internal/modifiers/force_string_value_unknown_modifier.go

## Problem

The `forceStringValueUnknownModifier` struct is defined as an empty struct and does not contain any fields. If the struct will not need fields in the future, consider using a singleton pattern or reducing boilerplate.

## Impact

This is a low severity issue. Unnecessary struct definitions can introduce boilerplate and minor performance overhead through allocations if many instances are created.

## Location

```go
type forceStringValueUnknownModifier struct {
}
```

## Fix

If no state or fields will ever be needed, use a package-level variable as a singleton:

```go
var forceStringValueUnknownModifierInstance = &forceStringValueUnknownModifier{}
```

And modify your factory function to return this instance.
