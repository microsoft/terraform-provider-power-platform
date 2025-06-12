# Structural Redundancy: Repeated `BusinessUnitId` Field

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/models.go

## Problem

The `BusinessUnitId` field appears in multiple models (`SecurityRolesListDataSourceModel`, `SecurityRoleDataSourceModel`, and `UserResourceModel`), which may suggest possible redundancy or missed opportunity for struct composition and reducing code duplication.

## Impact

Maintenance is more difficult if changes are needed in all locations, and the code can be less DRY (Don't Repeat Yourself). Severity: **low**, unless this leads to inconsistent usage or validation.

## Location

```go
type SecurityRolesListDataSourceModel struct {
	BusinessUnitId types.String
}
type SecurityRoleDataSourceModel struct {
	BusinessUnitId types.String
}
type UserResourceModel struct {
	BusinessUnitId types.String
}
```

## Code Issue

```go
BusinessUnitId types.String
```

## Fix

Consider defining a base struct for models that share common fields, then compose it:

```go
type BusinessUnitModel struct {
	BusinessUnitId types.String
}

type SecurityRolesListDataSourceModel struct {
	BusinessUnitModel
	...
}
type SecurityRoleDataSourceModel struct {
	BusinessUnitModel
	...
}
type UserResourceModel struct {
	BusinessUnitModel
	...
}
```
