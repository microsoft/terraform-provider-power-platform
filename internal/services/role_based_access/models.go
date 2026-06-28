// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

// Resource structs

type roleBasedAccessAssignmentResource struct {
	helpers.TypeInfo
	Client RoleBasedAccessClient
}

type environmentGroupRoleBasedAccessAssignmentResource struct {
	helpers.TypeInfo
	Client RoleBasedAccessClient
}

type environmentRoleBasedAccessAssignmentResource struct {
	helpers.TypeInfo
	Client RoleBasedAccessClient
}

// Resource models

type roleBasedAccessAssignmentResourceModel struct {
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
	Id                types.String   `tfsdk:"id"`
	PrincipalObjectId types.String   `tfsdk:"principal_object_id"`
	PrincipalType     types.String   `tfsdk:"principal_type"`
	RoleDefinitionId  types.String   `tfsdk:"role_definition_id"`
	Scope             types.String   `tfsdk:"scope"`
	CreatedOn         types.String   `tfsdk:"created_on"`
}

type environmentGroupRoleBasedAccessAssignmentResourceModel struct {
	Timeouts           timeouts.Value `tfsdk:"timeouts"`
	Id                 types.String   `tfsdk:"id"`
	EnvironmentGroupId types.String   `tfsdk:"environment_group_id"`
	PrincipalObjectId  types.String   `tfsdk:"principal_object_id"`
	PrincipalType      types.String   `tfsdk:"principal_type"`
	RoleDefinitionId   types.String   `tfsdk:"role_definition_id"`
	Scope              types.String   `tfsdk:"scope"`
	CreatedOn          types.String   `tfsdk:"created_on"`
}

type environmentRoleBasedAccessAssignmentResourceModel struct {
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
	Id                types.String   `tfsdk:"id"`
	EnvironmentId     types.String   `tfsdk:"environment_id"`
	PrincipalObjectId types.String   `tfsdk:"principal_object_id"`
	PrincipalType     types.String   `tfsdk:"principal_type"`
	RoleDefinitionId  types.String   `tfsdk:"role_definition_id"`
	Scope             types.String   `tfsdk:"scope"`
	CreatedOn         types.String   `tfsdk:"created_on"`
}

// Data source structs

type roleDefinitionsDataSource struct {
	helpers.TypeInfo
	Client RoleBasedAccessClient
}

type roleDefinitionsDataSourceModel struct {
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
	RoleDefinitions types.List     `tfsdk:"role_definitions"`
}

type roleDefinitionModel struct {
	RoleDefinitionId   types.String `tfsdk:"role_definition_id"`
	RoleDefinitionName types.String `tfsdk:"role_definition_name"`
	Permissions        types.List   `tfsdk:"permissions"`
}
