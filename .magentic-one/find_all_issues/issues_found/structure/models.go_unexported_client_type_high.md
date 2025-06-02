# Title

Unexported Client Type in Struct Field

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/models.go

## Problem

The field `client` in `AdminManagementApplicationResource` is referenced as a type but is not defined or imported in the file. It also uses a lowercase name, indicating it might be an unexported (package-private) struct or interface, which can cause issues if this file is intended to be imported or used outside of its package, or if this type is not actually defined.

## Impact

If the `client` type is not defined anywhere, this will cause a compile-time error. If it is defined but unexported and needs to be used outside this package, it will not be accessible, leading to potential design and maintainability problems. Severity: high

## Location

Line 13, struct: `AdminManagementApplicationResource`

## Code Issue

```go
type AdminManagementApplicationResource struct {
	helpers.TypeInfo
	AdminManagementApplicationClient client
}
```

## Fix

Ensure that the `client` type is properly defined in the package and/or should be accessible if intended for broader use. If `client` is a variable name, it should be preceded by its type; otherwise, capitalize it if necessary for exporting.

Example fix if you meant to reference a `Client` of type `AdminManagementApplicationClient`:

```go
type AdminManagementApplicationResource struct {
	helpers.TypeInfo
	Client AdminManagementApplicationClient
}
```

Or, if `client` type is not defined, import or define it appropriately. If it should remain package-private and is defined elsewhere, clarify the code structure for maintainability.
