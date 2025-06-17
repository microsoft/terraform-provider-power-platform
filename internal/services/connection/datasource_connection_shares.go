// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connection

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
	_ datasource.DataSource              = &ConnectionsDataSource{}
	_ datasource.DataSourceWithConfigure = &ConnectionsDataSource{}
)

func NewConnectionSharesDataSource() datasource.DataSource {
	return &SharesDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "connection_shares",
		},
	}
}

func (d *SharesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *SharesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Lists connection shares found in an environment. For more information about the connection sharing, see [Share a connection reference](https://learn.microsoft.com/en-us/power-automate/connection-sharing)",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Environment Id. The unique identifier of the environment that the connection are associated with.",
				Required:            true,
			},
			"connector_name": schema.StringAttribute{
				MarkdownDescription: "Connector Name. The unique identifier of the connector that the connection are associated with.",
				Required:            true,
			},
			"connection_id": schema.StringAttribute{
				MarkdownDescription: "Connection Id. The unique identifier of the connection that the shares are associated with.",
				Required:            true,
			},
			"shares": schema.ListNestedAttribute{
				MarkdownDescription: "List of shares for a given connection.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique share id",
							Computed:            true,
						},
						"role_name": schema.StringAttribute{
							MarkdownDescription: "Role name of the share",
							Computed:            true,
						},
						"principal": schema.SingleNestedAttribute{
							MarkdownDescription: "",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"entra_object_id": schema.StringAttribute{
									MarkdownDescription: "Entra Object Id of the principal",
									Computed:            true,
								},
								"display_name": schema.StringAttribute{
									MarkdownDescription: "Principal Display Name",
									Computed:            true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *SharesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.ConnectionsClient = newConnectionsClient(client.Api)
}
func (d *SharesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state SharesListDataSourceModel
	resp.State.Get(ctx, &state)

	connectionsList, err := d.ConnectionsClient.GetConnectionShares(ctx, state.EnvironmentId.ValueString(), state.ConnectorName.ValueString(), state.ConnectionId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf(
				"Failed to get connection shares for environment_id '%s', connector_name '%s', connection_id '%s'",
				state.EnvironmentId.ValueString(), state.ConnectorName.ValueString(), state.ConnectionId.ValueString(),
			),
			err.Error(),
		)
		return
	}

	for _, connection := range connectionsList.Value {
		connectionModel := ConvertFromConnectionSharesDto(connection)
		state.Shares = append(state.Shares, connectionModel)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func ConvertFromConnectionSharesDto(connection shareConnectionResponseDto) SharesDataSourceModel {
	share := SharesDataSourceModel{
		Id:       types.StringValue(connection.Name),
		RoleName: types.StringValue(connection.Properties.RoleName),
	}

	if displayName, ok := connection.Properties.Principal["displayName"].(string); ok {
		share.Principal.DisplayName = types.StringValue(displayName)
	} else {
		share.Principal.DisplayName = types.StringValue("")
	}

	if entraId, ok := connection.Properties.Principal["id"].(string); ok {
		share.Principal.EntraId = types.StringValue(entraId)
	} else {
		share.Principal.EntraId = types.StringValue("")
	}

	return share
}
