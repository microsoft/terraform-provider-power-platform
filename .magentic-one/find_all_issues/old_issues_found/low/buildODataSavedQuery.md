# Title

Direct pointer comparison in `buildODataSavedQuery`, `buildODataUserQuery`, and similar functions

##

`/workspaces/terraform-provider-power-platform/internal/services/data_record/odata.go`

## Problem

Direct comparison of pointer values (`nil` or `non-nil`) is used to determine if a value is set. In Go, although this works, it is prone to bugs when dealing with complex or shared pointer data structures.

## Impact

Incorrect handling or mutation of pointers can lead to runtime errors or unexpected behavior. While this is less severe (low severity) than the first issue, it represents a design flaw that can impact maintainability.

## Location

Functions such as:
1. `buildODataSavedQuery`
2. `buildODataUserQuery`
3. `buildODataFilterPart`, etc.

## Code Issue

Example from `buildODataFilterPart`:
```go
if filter != nil {
    resultQuery = "$filter=" + url.QueryEscape(*filter)
}
```

## Fix

Include explicit pointer dereferencing checks to avoid unintended nil dereferences and improve readability.

```go
if filter != nil && *filter != "" {
    resultQuery = "$filter=" + url.QueryEscape(*filter)
}
```
