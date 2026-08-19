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
	_ datasource.DataSource              = &roleBasedAccessAssignmentsDataSource{}
	_ datasource.DataSourceWithConfigure = &roleBasedAccessAssignmentsDataSource{}
)

// roleAssignmentsAttribute returns the shared schema of the list of role assignments returned by the RBAC API.
func roleAssignmentsAttribute(markdownDescription string) schema.ListNestedAttribute {
	return schema.ListNestedAttribute{
		MarkdownDescription: markdownDescription,
		Computed:            true,
		NestedObject: schema.NestedAttributeObject{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					MarkdownDescription: "The unique identifier of the role assignment",
					Computed:            true,
				},
				"scope": schema.StringAttribute{
					MarkdownDescription: "The scope the role assignment applies to",
					Computed:            true,
				},
				"principal_type": schema.StringAttribute{
					MarkdownDescription: "The type of principal the role is assigned to (e.g., `ApplicationUser`, `User`)",
					Computed:            true,
				},
				"enterprise_application_object_id": schema.StringAttribute{
					MarkdownDescription: "The object ID of the enterprise application (service principal) or user the role is assigned to. For `ApplicationUser` principals this is the enterprise application object ID, not the application (client) ID",
					Computed:            true,
				},
				"role_definition_id": schema.StringAttribute{
					MarkdownDescription: "The ID of the assigned role definition",
					Computed:            true,
				},
				"created_by_principal_type": schema.StringAttribute{
					MarkdownDescription: "The type of principal that created the role assignment",
					Computed:            true,
				},
				"created_by_principal_object_id": schema.StringAttribute{
					MarkdownDescription: "The object ID of the principal that created the role assignment",
					Computed:            true,
				},
				"created_on": schema.StringAttribute{
					MarkdownDescription: "The timestamp when the role assignment was created",
					Computed:            true,
				},
				"expires_on": schema.StringAttribute{
					MarkdownDescription: "The timestamp when the role assignment expires, or `null` if it does not expire",
					Computed:            true,
				},
			},
		},
	}
}

func NewRoleBasedAccessAssignmentsDataSource() datasource.DataSource {
	return &roleBasedAccessAssignmentsDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "role_based_access_assignments",
		},
	}
}

func (d *roleBasedAccessAssignmentsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *roleBasedAccessAssignmentsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches the tenant scoped [role assignments](https://learn.microsoft.com/en-us/rest/api/power-platform/authorization/role-based-access-control/list-role-assignments) in Power Platform. " +
			"Use this data source to discover which principals are assigned tenant level roles.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"role_assignments": roleAssignmentsAttribute("List of tenant scoped role assignments"),
		},
	}
}

func (d *roleBasedAccessAssignmentsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *roleBasedAccessAssignmentsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state roleBasedAccessAssignmentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Listing tenant role assignments")

	assignments, err := d.Client.ListRoleAssignments(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list role assignments", err.Error())
		return
	}

	state.RoleAssignments = make([]roleAssignmentDataSourceModel, 0, len(assignments))
	for _, assignment := range assignments {
		state.RoleAssignments = append(state.RoleAssignments, convertRoleAssignmentDtoToDataSourceModel(assignment))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
