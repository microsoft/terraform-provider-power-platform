// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package authorization

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var (
	_ datasource.DataSource              = &SecurityRolesDataSource{}
	_ datasource.DataSourceWithConfigure = &SecurityRolesDataSource{}
)

type SecurityRolesDataSource struct {
	helpers.TypeInfo
	UserClient client
}

func NewSecurityRolesDataSource() datasource.DataSource {
	return &SecurityRolesDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "security_roles",
		},
	}
}

func (d *SecurityRolesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches the list of Dataverse security roles for a given environment and business unit.  For more information see [About security roles and privileges](https://learn.microsoft.com/power-platform/admin/security-roles-privileges)",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Id of the Dynamics 365 environment",
				Required:            true,
			},
			"business_unit_id": schema.StringAttribute{
				MarkdownDescription: "Id of the business unit to filter the security roles",
				Optional:            true,
			},
			"security_roles": schema.ListNestedAttribute{
				MarkdownDescription: "List of security roles",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"role_id": schema.StringAttribute{
							MarkdownDescription: "Security role id",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Security role name",
							Computed:            true,
						},
						"is_managed": schema.BoolAttribute{
							MarkdownDescription: "Is the security role managed",
							Computed:            true,
						},
						"business_unit_id": schema.StringAttribute{
							MarkdownDescription: "Id of the business unit",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *SecurityRolesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
		return
	}

	client, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.UserClient = newUserClient(client.Api)
}

func (d *SecurityRolesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *SecurityRolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state SecurityRolesListDataSourceModel
	resp.State.Get(ctx, &state)

	if state.EnvironmentId.ValueString() == "" {
		resp.Diagnostics.AddError("environment_id connot be an empty string", "environment_id connot be an empty string")
		return
	}

	dvExits, err := d.UserClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
	tflog.Debug(ctx, fmt.Sprintf("Environment Id: %s", state.EnvironmentId.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
	}

	if !dvExits {
		resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), "")
		return
	}

	roles, err := d.UserClient.GetDataverseSecurityRoles(ctx, state.EnvironmentId.ValueString(), state.BusinessUnitId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", d.ProviderTypeName, d.TypeName), err.Error())
		return
	}

	for _, role := range roles {
		state.SecurityRoles = append(state.SecurityRoles, SecurityRoleDataSourceModel{
			RoleId:         types.StringValue(role.RoleId),
			Name:           types.StringValue(role.Name),
			IsManaged:      types.BoolValue(role.IsManaged),
			BusinessUnitId: types.StringValue(role.BusinessUnitId),
		})
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
