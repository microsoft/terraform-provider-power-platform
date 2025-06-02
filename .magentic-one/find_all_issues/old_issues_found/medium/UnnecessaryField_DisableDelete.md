# Title
Unnecessary Inclusion of `disable_delete` Field in `UserResourceModel`

##

`/workspaces/terraform-provider-power-platform/internal/services/authorization/models.go`

## Problem

The `UserResourceModel` struct includes the field `DisableDelete`. This field's purpose is unclear within the broader context and may result in complexity or confusion during maintenance.

## Impact

Including unnecessary fields increases code clutter and the cognitive load required to understand and interact with the codebase. Severity: Medium.

## Location

The issue is located within the `UserResourceModel` declaration.

## Code Issue

```go
type UserResourceModel struct {
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
	Id                types.String   `tfsdk:"id"`
	EnvironmentId     types.String   `tfsdk:"environment_id"`
	AadId             types.String   `tfsdk:"aad_id"`
	BusinessUnitId    types.String   `tfsdk:"business_unit_id"`
	SecurityRoles     []string       `tfsdk:"security_roles"`
	UserPrincipalName types.String   `tfsdk:"user_principal_name"`
	FirstName         types.String   `tfsdk:"first_name"`
	LastName          types.String   `tfsdk:"last_name"`
	DisableDelete     types.Bool     `tfsdk:"disable_delete"`
}
```

## Fix

Perform a system-wide review to determine if `DisableDelete` is actively used. If it is deemed unnecessary, remove it to simplify the code.

```go
type UserResourceModel struct {
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
	Id                types.String   `tfsdk:"id"`
	EnvironmentId     types.String   `tfsdk:"environment_id"`
	AadId             types.String   `tfsdk:"aad_id"`
	BusinessUnitId    types.String   `tfsdk:"business_unit_id"`
	SecurityRoles     []string       `tfsdk:"security_roles"`
	UserPrincipalName types.String   `tfsdk:"user_principal_name"`
	FirstName         types.String   `tfsdk:"first_name"`
	LastName          types.String   `tfsdk:"last_name"`
}
```