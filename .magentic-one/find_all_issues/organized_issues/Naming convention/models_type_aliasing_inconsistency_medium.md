# Type Aliasing Inconsistency

##

/workspaces/terraform-provider-power-platform/internal/services/connection/models.go

## Problem

The model structs representing similar concepts use inconsistent naming strategies. For example, `SharesDataSourceModel` and `ShareResourceModel` differ in singular/plural naming, while their internal types use slightly different names (`SharesPrincipalDataSourceModel` vs `SharePrincipalResourceModel`). This inconsistent naming can make the codebase difficult to maintain and understand.

## Impact

Medium. Inconsistent naming affects readability and maintainability. Developers may misinterpret the purpose or role of a struct, leading to subtle bugs or confusion when making changes or integrating with other components.

## Location

- Types:
  - `SharesDataSourceModel`
  - `ShareResourceModel`
  - `SharesPrincipalDataSourceModel`
  - `SharePrincipalResourceModel`
- File: `/workspaces/terraform-provider-power-platform/internal/services/connection/models.go`

## Code Issue

```go
type SharesDataSourceModel struct {
	// ...
	Principal SharesPrincipalDataSourceModel `tfsdk:"principal"`
}

type ShareResourceModel struct {
	// ...
	Principal     SharePrincipalResourceModel `tfsdk:"principal"`
}

type SharesPrincipalDataSourceModel struct {
	EntraId     types.String `tfsdk:"entra_object_id"`
	DisplayName types.String `tfsdk:"display_name"`
}

type SharePrincipalResourceModel struct {
	EntraObjectId types.String `tfsdk:"entra_object_id"`
	DisplayName   types.String `tfsdk:"display_name"`
}
```

## Fix

Unify naming for consistency. Use either singular or plural consistently and choose one naming pattern (e.g., always `SharePrincipalModel`, with `DataSource` or `Resource` suffix as needed):

```go
type SharePrincipalModel struct {
	EntraObjectId types.String `tfsdk:"entra_object_id"`
	DisplayName   types.String `tfsdk:"display_name"`
}

type SharesDataSourceModel struct {
	Id        types.String        `tfsdk:"id"`
	RoleName  types.String        `tfsdk:"role_name"`
	Principal SharePrincipalModel `tfsdk:"principal"`
}

type ShareResourceModel struct {
	Timeouts      timeouts.Value     `tfsdk:"timeouts"`
	Id            types.String       `tfsdk:"id"`
	EnvironmentId types.String       `tfsdk:"environment_id"`
	ConnectorName types.String       `tfsdk:"connector_name"`
	ConnectionId  types.String       `tfsdk:"connection_id"`
	RoleName      types.String       `tfsdk:"role_name"`
	Principal     SharePrincipalModel `tfsdk:"principal"`
}
```

This approach enhances consistency and makes refactoring or understanding the models easier.
