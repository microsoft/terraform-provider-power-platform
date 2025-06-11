# Title 

Function Naming - StringPtr

##

/workspaces/terraform-provider-power-platform/internal/helpers/config.go

## Problem

The `StringPtr` naming does not follow Go idiomatic naming, which typically suggests `StringPointer`.

## Impact

Low severity. Minor impact, mostly readability and consistency.

## Location

Line 81-83

## Code Issue

```go
func StringPtr(s string) *string {
	return &s
}
```

## Fix

Rename the function to `StringPointer` for consistency with Go naming conventions for such helpers.

```go
func StringPointer(s string) *string {
	return &s
}
```
