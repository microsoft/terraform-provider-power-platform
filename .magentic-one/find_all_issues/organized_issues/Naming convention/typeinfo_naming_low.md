# Naming Consistency for Struct Fields

##

/workspaces/terraform-provider-power-platform/internal/helpers/typeinfo.go

## Problem

The field names in the `TypeInfo` struct use mixed naming conventions. `ProviderTypeName` is quite descriptive, but `TypeName` is somewhat generic. Although this is not strictly incorrect, clearer naming such as `ResourceTypeName` or `SpecificTypeName` might improve code readability.

## Impact

Low. This is a minor readability issue. The current names are functional, but more descriptive names could reduce future confusion and improve maintainability.

## Location

```go
type TypeInfo struct {
	ProviderTypeName string
	TypeName         string
}
```

## Code Issue

```go
type TypeInfo struct {
	ProviderTypeName string
	TypeName         string
}
```

## Fix

Consider renaming `TypeName` to something more descriptive to clarify what type name is represented.

```go
type TypeInfo struct {
	ProviderTypeName string
	ResourceTypeName string // Was TypeName
}
```
