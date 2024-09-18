// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connection

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/constants"
)

var (
	_ datasource.DataSource              = &ConnectionsDataSource{}
	_ datasource.DataSourceWithConfigure = &ConnectionsDataSource{}
)

func NewConnectionSharesDataSource() datasource.DataSource {
	return &SharesDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_connection_shares",
	}
}

type SharesDataSource struct {
	ConnectionsClient ConnectionsClient
	ProviderTypeName  string
	TypeName          string
}

type SharesListDataSourceModel struct {
	Timeouts      timeouts.Value          `tfsdk:"timeouts"`
	Id            types.String            `tfsdk:"id"`
	EnvironmentId types.String            `tfsdk:"environment_id"`
	ConnectorName types.String            `tfsdk:"connector_name"`
	ConnectionId  types.String            `tfsdk:"connection_id"`
	Shares        []SharesDataSourceModel `tfsdk:"shares"`
}

type SharesDataSourceModel struct {
	Id        types.String                   `tfsdk:"id"`
	RoleName  types.String                   `tfsdk:"role_name"`
	Principal SharesPrincipalDataSourceModel `tfsdk:"principal"`
}

type SharesPrincipalDataSourceModel struct {
	EntraId     types.String `tfsdk:"entra_object_id"`
	DisplayName types.String `tfsdk:"display_name"`
}

func (d *SharesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *SharesDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
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

func (d *SharesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *SharesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state SharesListDataSourceModel
	resp.State.Get(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE START: %s", d.ProviderTypeName))

	timeout, diags := state.Timeouts.Read(ctx, constants.DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

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

	diags = resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func ConvertFromConnectionSharesDto(connection ShareConnectionResponseDto) SharesDataSourceModel {
	share := SharesDataSourceModel{
		Id:       types.StringValue(connection.Name),
		RoleName: types.StringValue(connection.Properties.RoleName),
		Principal: SharesPrincipalDataSourceModel{
			EntraId:     types.StringValue(connection.Properties.Principal["id"].(string)),
			DisplayName: types.StringValue(connection.Properties.Principal["displayName"].(string)),
		},
	}
	return share
}
