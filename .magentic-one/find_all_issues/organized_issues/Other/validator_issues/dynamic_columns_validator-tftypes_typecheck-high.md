# Issue 3: Incorrect Type Handling for tftypes.Value Consistency

## 

/workspaces/terraform-provider-power-platform/internal/services/data_record/dynamic_columns_validator.go

## Problem

Currently, the code uses:

```go
if value.Type().Is(tftypes.Tuple{}) || value.Type().Is(tftypes.List{}) {
```

This is not the idiomatic way to check if a type is a List or Tuple within the `tftypes` package. The `.Is()` method expects an instance of `tftypes.Type`, so the code should use `tftypes.List{}` and `tftypes.Tuple{}` types from package root, not composite literals for the comparison.

## Impact

The code could behave unexpectedly, always returning false from the check, thus not warning users of incorrect relationships. Severity: **high**.

## Location

Within the `Validate` function, loop for attribute validation.

## Code Issue

```go
if value.Type().Is(tftypes.Tuple{}) || value.Type().Is(tftypes.List{}) {
```

## Fix

Declare target type variables and use them for comparison:

```go
if value.Type().Is(&tftypes.Tuple{}) || value.Type().Is(&tftypes.List{}) {
    // ...
}
```

Or, check using type switches:

```go
switch value.Type().(type) {
case *tftypes.List, *tftypes.Tuple:
    msg := fmt.Sprintf("Dynamic columns should use set collection with `toset([...])` for many-to-one relationships. Record attribute: '%s'", key)
    diags.AddWarning(msg, msg)
}
```
