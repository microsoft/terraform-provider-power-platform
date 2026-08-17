// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package role_based_access

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var (
	_ datasource.DataSource              = &roleDefinitionsDataSource{}
	_ datasource.DataSourceWithConfigure = &roleDefinitionsDataSource{}
)

func NewRoleDefinitionsDataSource() datasource.DataSource {
	return &roleDefinitionsDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "role_definitions",
		},
	}
}

func (d *roleDefinitionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *roleDefinitionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches available [role definitions](https://learn.microsoft.com/en-us/rest/api/power-platform/authorization/role-based-access-control/list-role-definitions) from Power Platform. " +
			"Use this data source to discover available roles that can be assigned to service principals or users at the tenant, environment group or environment scope.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"role_definitions": schema.ListNestedAttribute{
				MarkdownDescription: "List of available role definitions",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"role_definition_id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier of the role definition",
							Computed:            true,
						},
						"role_definition_name": schema.StringAttribute{
							MarkdownDescription: "The display name of the role",
							Computed:            true,
						},
						"permissions": schema.ListAttribute{
							MarkdownDescription: "List of permissions granted by this role",
							Computed:            true,
							ElementType:         types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *roleDefinitionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *roleDefinitionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state roleDefinitionsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Listing role definitions")

	definitions, err := d.Client.ListRoleDefinitions(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list role definitions", err.Error())
		return
	}

	roleDefAttrTypes := map[string]attr.Type{
		"role_definition_id":   types.StringType,
		"role_definition_name": types.StringType,
		"permissions":          types.ListType{ElemType: types.StringType},
	}

	roleDefValues := make([]attr.Value, len(definitions))
	for i, def := range definitions {
		permValues := make([]attr.Value, len(def.Permissions))
		for j, perm := range def.Permissions {
			permValues[j] = types.StringValue(perm)
		}
		permList, diags := types.ListValue(types.StringType, permValues)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		roleDefObj, diags := types.ObjectValue(roleDefAttrTypes, map[string]attr.Value{
			"role_definition_id":   types.StringValue(def.RoleDefinitionId),
			"role_definition_name": types.StringValue(def.RoleDefinitionName),
			"permissions":          permList,
		})
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		roleDefValues[i] = roleDefObj
	}

	roleDefList, diags := types.ListValue(types.ObjectType{AttrTypes: roleDefAttrTypes}, roleDefValues)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state.RoleDefinitions = roleDefList
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
