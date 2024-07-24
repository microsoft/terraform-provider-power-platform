// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var (
	_ datasource.DataSource              = &ConnectionsDataSource{}
	_ datasource.DataSourceWithConfigure = &ConnectionsDataSource{}
)

func NewConnectionSharesDataSource() datasource.DataSource {
	return &ConnectionSharesDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_connection_shares",
	}
}

type ConnectionSharesDataSource struct {
	ConnectionsClient ConnectionsClient
	ProviderTypeName  string
	TypeName          string
}

type ConnectionSharesListDataSourceModel struct {
	Id            types.String                      `tfsdk:"id"`
	EnvironmentId types.String                      `tfsdk:"environment_id"`
	ConnectorName types.String                      `tfsdk:"connector_name"`
	ConnectionId  types.String                      `tfsdk:"connection_id"`
	Shares        []ConnectionSharesDataSourceModel `tfsdk:"shares"`
}

type ConnectionSharesDataSourceModel struct {
	Id         types.String                              `tfsdk:"id"`
	Properties ConnectionSharesPropertiesDataSourceModel `tfsdk:"properties"`
}

type ConnectionSharesPropertiesDataSourceModel struct {
	RoleName  types.String                            `tfsdk:"role_name"`
	Principal ConnectionShresPrincipalDataSourceModel `tfsdk:"principal"`
}

type ConnectionShresPrincipalDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	DisplayName types.String `tfsdk:"display_name"`
}

func (d *ConnectionSharesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *ConnectionSharesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"environment_id": schema.StringAttribute{
				Description:         "Environment Id. The unique identifier of the environment that the connection are associated with.",
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
						"properties": schema.SingleNestedAttribute{
							MarkdownDescription: "",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"role_name": schema.StringAttribute{
									MarkdownDescription: "Role name of the share",
									Computed:            true,
								},
								"principal": schema.SingleNestedAttribute{
									MarkdownDescription: "",
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"id": schema.StringAttribute{
											MarkdownDescription: "Principal Id",
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
			},
		},
	}
}

func (d *ConnectionSharesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client := req.ProviderData.(*api.ProviderClient).Api

	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.ConnectionsClient = NewConnectionsClient(client)
}
func (d *ConnectionSharesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ConnectionSharesListDataSourceModel
	resp.State.Get(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE START: %s", d.ProviderTypeName))

	connectionsList, err := d.ConnectionsClient.GetConnectionShares(ctx, state.EnvironmentId.ValueString(), state.ConnectorName.ValueString(), state.ConnectionId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to get connection shares", err.Error())
		return
	}

	for _, connection := range connectionsList.Value {
		connectionModel := ConvertFromConnectionSharesDto(connection)
		state.Shares = append(state.Shares, connectionModel)
	}
	state.Id = types.StringValue(strconv.Itoa(len(connectionsList.Value)))

	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func ConvertFromConnectionSharesDto(connection ShareConnectionResponseDto) ConnectionSharesDataSourceModel {
	share := ConnectionSharesDataSourceModel{
		Id: types.StringValue(connection.Name),
		Properties: ConnectionSharesPropertiesDataSourceModel{
			RoleName: types.StringValue(connection.Properties.RoleName),
			Principal: ConnectionShresPrincipalDataSourceModel{
				Id:          types.StringValue(connection.Properties.Principal.Id),
				DisplayName: types.StringValue(connection.Properties.Principal.DisplayName),
			},
		},
	}

	return share
}
