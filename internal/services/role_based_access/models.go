// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/customtypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

// Resource structs

type rbacRoleAssignmentResource struct {
	helpers.TypeInfo
	Client client
}

// Resource models

type rbacRoleAssignmentResourceModel struct {
	Timeouts           timeouts.Value                    `tfsdk:"timeouts"`
	Id                 types.String                      `tfsdk:"id"`
	ScopeType          types.String                      `tfsdk:"scope_type"`
	EnvironmentId      customtypes.UUID                  `tfsdk:"environment_id"`
	EnvironmentGroupId customtypes.UUID                  `tfsdk:"environment_group_id"`
	PrincipalId        customtypes.UUID                  `tfsdk:"principal_id"`
	PrincipalType      types.String                      `tfsdk:"principal_type"`
	RoleDefinitionId   customtypes.UUID                  `tfsdk:"role_definition_id"`
	RoleDefinitionName customtypes.CaseInsensitiveString `tfsdk:"role_definition_name"`
	Scope              types.String                      `tfsdk:"scope"`
	CreatedOn          types.String                      `tfsdk:"created_on"`
}

// assignmentScope is where the assignment applies. The kind comes from the `scope_type` attribute rather
// than being inferred, so an empty or missing id is a validation error instead of a tenant grant.
func (m rbacRoleAssignmentResourceModel) assignmentScope() assignmentScope {
	switch m.ScopeType.ValueString() {
	case scopeEnvironment:
		return environmentAssignmentScope(m.EnvironmentId.ValueString())
	case scopeEnvironmentGroup:
		return environmentGroupAssignmentScope(m.EnvironmentGroupId.ValueString())
	default:
		return tenantAssignmentScope()
	}
}

// Data source structs

type roleDefinitionsDataSource struct {
	helpers.TypeInfo
	Client client
}

type roleDefinitionsDataSourceModel struct {
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
	RoleDefinitions types.List     `tfsdk:"role_definitions"`
}

type rbacRoleAssignmentsDataSource struct {
	helpers.TypeInfo
	Client client
}

// Data source models

type rbacRoleAssignmentDataSourceModel struct {
	Id                         types.String `tfsdk:"id"`
	Scope                      types.String `tfsdk:"scope"`
	PrincipalType              types.String `tfsdk:"principal_type"`
	PrincipalId                types.String `tfsdk:"principal_id"`
	RoleDefinitionId           types.String `tfsdk:"role_definition_id"`
	CreatedByPrincipalType     types.String `tfsdk:"created_by_principal_type"`
	CreatedByPrincipalObjectId types.String `tfsdk:"created_by_principal_object_id"`
	CreatedOn                  types.String `tfsdk:"created_on"`
	ExpiresOn                  types.String `tfsdk:"expires_on"`
}

type rbacRoleAssignmentsDataSourceModel struct {
	Timeouts           timeouts.Value                      `tfsdk:"timeouts"`
	ScopeType          types.String                        `tfsdk:"scope_type"`
	EnvironmentId      customtypes.UUID                    `tfsdk:"environment_id"`
	EnvironmentGroupId customtypes.UUID                    `tfsdk:"environment_group_id"`
	RoleAssignments    []rbacRoleAssignmentDataSourceModel `tfsdk:"role_assignments"`
}

// assignmentScope is which assignments to read, taken from the explicit `scope_type` attribute.
func (m rbacRoleAssignmentsDataSourceModel) assignmentScope() assignmentScope {
	switch m.ScopeType.ValueString() {
	case scopeEnvironment:
		return environmentAssignmentScope(m.EnvironmentId.ValueString())
	case scopeEnvironmentGroup:
		return environmentGroupAssignmentScope(m.EnvironmentGroupId.ValueString())
	default:
		return tenantAssignmentScope()
	}
}

func convertRbacRoleAssignmentDtoToDataSourceModel(assignment roleAssignmentDto) rbacRoleAssignmentDataSourceModel {
	model := rbacRoleAssignmentDataSourceModel{
		Id:                         types.StringValue(assignment.RoleAssignmentId),
		Scope:                      types.StringValue(assignment.Scope),
		PrincipalType:              types.StringValue(assignment.PrincipalType),
		PrincipalId:                types.StringValue(assignment.PrincipalObjectId),
		RoleDefinitionId:           types.StringValue(assignment.RoleDefinitionId),
		CreatedByPrincipalType:     types.StringValue(assignment.CreatedByPrincipalType),
		CreatedByPrincipalObjectId: types.StringValue(assignment.CreatedByPrincipalObjectId),
		CreatedOn:                  types.StringValue(assignment.CreatedOn),
		ExpiresOn:                  types.StringNull(),
	}
	if assignment.ExpiresOn != nil {
		model.ExpiresOn = types.StringValue(*assignment.ExpiresOn)
	}
	return model
}
