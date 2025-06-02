# Lack of Validation for Struct Inputs

##

/workspaces/terraform-provider-power-platform/internal/helpers/typeinfo.go

## Problem

The `TypeInfo` struct does not provide any input validation when creating instances. For example, it is possible to create an object with an empty `TypeName`, which would lead to an invalid type string like `powerplatform_` in `FullTypeName`. Having such invalid names could propagate hidden errors in larger contexts.

## Impact

Medium. Absence of validation could result in invalid resource or data source type names, possibly causing issues during downstream operations or harming user experience.

## Location

```go
type TypeInfo struct {
	ProviderTypeName string
	TypeName         string
}
```
...
```go
func (t *TypeInfo) FullTypeName() string {
	if t.ProviderTypeName == "" {
		return fmt.Sprintf("powerplatform_%s", t.TypeName)
	}

	return fmt.Sprintf("%s_%s", t.ProviderTypeName, t.TypeName)
}
```

## Code Issue

```go
func (t *TypeInfo) FullTypeName() string {
	if t.ProviderTypeName == "" {
		return fmt.Sprintf("powerplatform_%s", t.TypeName)
	}

	return fmt.Sprintf("%s_%s", t.ProviderTypeName, t.TypeName)
}
```

## Fix

Add validation for `TypeName` when constructing `TypeInfo` or running `FullTypeName`, and return an error if it’s missing or invalid.

```go
func (t *TypeInfo) FullTypeName() (string, error) {
	if t.TypeName == "" {
		return "", fmt.Errorf("TypeName must not be empty")
	}

	if t.ProviderTypeName == "" {
		return fmt.Sprintf("powerplatform_%s", t.TypeName), nil
	}

	return fmt.Sprintf("%s_%s", t.ProviderTypeName, t.TypeName), nil
}
```
Or enforce TypeInfo creation only via constructor.
