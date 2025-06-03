# Struct field `Id` naming does not match Go naming conventions (resource identifiers commonly named `ID`)

##

/workspaces/terraform-provider-power-platform/internal/services/rest/resource_rest.go

## Problem

The struct field is named `Id` in `DataverseWebApiResourceModel`, rather than using the Go convention of `ID`. This goes against standard Go conventions and might cause confusion, as abbreviations should be capitalized in Go struct fields.

## Impact

Low (naming only, stylistic), but impacts maintainability and idiomaticity.

## Location

```go
type DataverseWebApiResourceModel struct {
	Timeouts timeouts.Value            `tfsdk:"timeouts"`
	Id       types.String              `tfsdk:"id"`
	...
}
```

## Code Issue

```go
	Id       types.String              `tfsdk:"id"`
```

## Fix

Rename to:

```go
	ID       types.String              `tfsdk:"id"`
```

And use `ID` throughout the codebase where the field is referenced.
