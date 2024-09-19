// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package authorization

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	UserClient UserClient
}

type SecurityRolesListDataSourceModel struct {
	Timeouts       timeouts.Value                `tfsdk:"timeouts"`
	Id             types.String                  `tfsdk:"id"`
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

func NewSecurityRolesDataSource() datasource.DataSource {
	return &SecurityRolesDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "security_roles",
		},
	}
}

func (_ *SecurityRolesDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of Dataverse security roles for a given environment and business unit",
		MarkdownDescription: "Fetches the list of Dataverse security roles for a given environment and business unit.  For more information see [About security roles and privileges](https://learn.microsoft.com/power-platform/admin/security-roles-privileges)",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"id": schema.StringAttribute{
				Description:         "Id of the read operation",
				MarkdownDescription: "Id of the read operation",
				Optional:            true,
			},
			"environment_id": schema.StringAttribute{
				Description:         "Id of the Dynamics 365 environment",
				MarkdownDescription: "Id of the Dynamics 365 environment",
				Required:            true,
			},
			"business_unit_id": schema.StringAttribute{
				Description: "Id of the business unit to filter the security roles",
				Optional:    true,
			},
			"security_roles": schema.ListNestedAttribute{
				Description:         "List of security roles",
				MarkdownDescription: "List of security roles",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"role_id": schema.StringAttribute{
							MarkdownDescription: "Security role id",
							Description:         "Security role id",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Security role name",
							Description:         "Security role name",
							Computed:            true,
						},
						"is_managed": schema.BoolAttribute{
							MarkdownDescription: "Is the security role managed",
							Description:         "Is the security role managed",
							Computed:            true,
						},
						"business_unit_id": schema.StringAttribute{
							MarkdownDescription: "Id of the business unit",
							Description:         "Id of the business unit",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *SecurityRolesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	clientApi := req.ProviderData.(*api.ProviderClient).Api
	if clientApi == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	d.UserClient = NewUserClient(clientApi)
}

func (d *SecurityRolesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *SecurityRolesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SecurityRolesListDataSourceModel
	resp.State.Get(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE SECURITY ROLES START: %s", d.ProviderTypeName))

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

	roles, err := d.UserClient.GetSecurityRoles(ctx, state.EnvironmentId.ValueString(), state.BusinessUnitId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", d.ProviderTypeName, d.TypeName), err.Error())
		return
	}

	state.Id = state.EnvironmentId

	for _, role := range roles {
		state.SecurityRoles = append(state.SecurityRoles, SecurityRoleDataSourceModel{
			RoleId:         types.StringValue(role.RoleId),
			Name:           types.StringValue(role.Name),
			IsManaged:      types.BoolValue(role.IsManaged),
			BusinessUnitId: types.StringValue(role.BusinessUnitId),
		})
	}

	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE SECURITY ROLES END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
