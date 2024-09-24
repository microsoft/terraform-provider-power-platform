// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package authorization

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type SecurityRolesListDataSourceModel struct {
	Timeouts       timeouts.Value                `tfsdk:"timeouts"`
	EnvironmentId  types.String                  `tfsdk:"environment_id"`
	BusinessUnitId types.String                  `tfsdk:"business_unit_id"`
	SecurityRoles  []SecurityRoleDataSourceModel `tfsdk:"security_roles"`
}

type SecurityRoleDataSourceModel struct {
	RoleId         types.String `tfsdk:"role_id"`
	Name           types.String `tfsdk:"name"`
	IsManaged      types.Bool   `tfsdk:"is_managed"`
	BusinessUnitId types.String `tfsdk:"business_unit_id"`
}

type UserResource struct {
	helpers.TypeInfo
	UserClient client
}

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
