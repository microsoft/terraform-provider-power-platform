// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var (
	_ datasource.DataSource              = &environmentGroupRoleBasedAccessAssignmentsDataSource{}
	_ datasource.DataSourceWithConfigure = &environmentGroupRoleBasedAccessAssignmentsDataSource{}
)

func NewEnvironmentGroupRoleBasedAccessAssignmentsDataSource() datasource.DataSource {
	return &environmentGroupRoleBasedAccessAssignmentsDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_group_role_based_access_assignments",
		},
	}
}

func (d *environmentGroupRoleBasedAccessAssignmentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *environmentGroupRoleBasedAccessAssignmentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches the [role assignments](https://learn.microsoft.com/en-us/rest/api/power-platform/authorization/role-based-access-control/list-environment-group-role-assignments) scoped to a given environment group in Power Platform. " +
			"Use this data source to discover which principals are assigned roles on an environment group.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"environment_group_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the environment group",
				Required:            true,
			},
			"role_assignments": roleAssignmentsAttribute("List of role assignments scoped to the environment group"),
		},
	}
}

func (d *environmentGroupRoleBasedAccessAssignmentsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		return
	}

	providerClient, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T.", req.ProviderData),
		)
		return
	}
	d.Client = newRoleBasedAccessClient(providerClient.Api)
}

func (d *environmentGroupRoleBasedAccessAssignmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state environmentGroupRoleBasedAccessAssignmentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Listing role assignments for environment group %s", state.EnvironmentGroupId.ValueString()))

	assignments, err := d.Client.ListEnvironmentGroupRoleAssignments(ctx, state.EnvironmentGroupId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to list environment group role assignments", err.Error())
		return
	}

	state.RoleAssignments = make([]roleAssignmentDataSourceModel, 0, len(assignments))
	for _, assignment := range assignments {
		state.RoleAssignments = append(state.RoleAssignments, convertRoleAssignmentDtoToDataSourceModel(assignment))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
