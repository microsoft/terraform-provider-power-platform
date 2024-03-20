// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var (
	_ datasource.DataSource              = &SecurityRolesDataSource{}
	_ datasource.DataSourceWithConfigure = &SecurityRolesDataSource{}
)

type SecurityRolesDataSource struct {
	UserClient       UserClient
	ProviderTypeName string
	TypeName         string
}

type SecurityRolesListDataSourceModel struct {
	Id            types.String                  `tfsdk:"id"`
	EnvironmentId types.String                  `tfsdk:"environment_id"`
	SecurityRoles []SecurityRoleDataSourceModel `tfsdk:"security_roles"`
}

type SecurityRoleDataSourceModel struct {
	RoleId         types.String `tfsdk:"role_id"`
	Name           types.String `tfsdk:"name"`
	IsManaged      types.Bool   `tfsdk:"is_managed"`
	BusinessUnitId types.String `tfsdk:"business_unit_id"`
}

func NewSecurityRolesDataSource() datasource.DataSource {
	return &SecurityRolesDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_security_roles",
	}
}

func (d *SecurityRolesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of Dataverse security roles for a given environment and business unit",
		MarkdownDescription: "Fetches the list of Dataverse security roles for a given environment and business unit",
		Attributes: map[string]schema.Attribute{
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
			// "business_unit_id": schema.StringAttribute{
			// 	Description: "Id of the business unit",
			// 	Required:    true,
			// },
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

func (d *SecurityRolesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *SecurityRolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SecurityRolesListDataSourceModel
	resp.State.Get(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE SECURITY ROLES START: %s", d.ProviderTypeName))

	if state.EnvironmentId.ValueString() == "" {
		resp.Diagnostics.AddError("environment_id connot be an empty string", "environment_id connot be an empty string")
		return
	}

	roles, err := d.UserClient.GetSecurityRoles(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
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

	state.Id = types.StringValue(fmt.Sprint((time.Now().Unix())))

	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE SECURITY ROLES END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
