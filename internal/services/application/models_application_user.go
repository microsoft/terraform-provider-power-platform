// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application

import (
	"sort"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type EnvironmentApplicationUserResource struct {
	helpers.TypeInfo
	ApplicationClient client
}

type EnvironmentApplicationUserResourceModel struct {
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
	Id                    types.String   `tfsdk:"id"`
	EnvironmentId         types.String   `tfsdk:"environment_id"`
	ApplicationId         types.String   `tfsdk:"application_id"`
	SystemUserId          types.String   `tfsdk:"system_user_id"`
	BusinessUnitId        types.String   `tfsdk:"business_unit_id"`
	SecurityRoles         []string       `tfsdk:"security_roles"`
	ResolvedSecurityRoles types.List     `tfsdk:"resolved_security_roles"`
}

var resolvedSecurityRoleObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"name":             types.StringType,
		"role_id":          types.StringType,
		"business_unit_id": types.StringType,
	},
}

func convertApplicationUserToModel(user *applicationUserDto) (types.String, []string, types.List) {
	roleNames := make([]string, 0, len(user.SecurityRoles))
	resolvedRoles := make([]attr.Value, 0, len(user.SecurityRoles))

	sort.Slice(user.SecurityRoles, func(i, j int) bool {
		return user.SecurityRoles[i].RoleId < user.SecurityRoles[j].RoleId
	})

	for _, role := range user.SecurityRoles {
		roleNames = append(roleNames, role.Name)
		resolvedRoles = append(resolvedRoles, types.ObjectValueMust(resolvedSecurityRoleObjectType.AttrTypes, map[string]attr.Value{
			"name":             types.StringValue(role.Name),
			"role_id":          types.StringValue(role.RoleId),
			"business_unit_id": types.StringValue(role.BusinessUnitId),
		}))
	}

	sort.Strings(roleNames)

	return types.StringValue(user.BusinessUnitId), roleNames, types.ListValueMust(resolvedSecurityRoleObjectType, resolvedRoles)
}
