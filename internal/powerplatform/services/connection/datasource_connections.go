// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"
	"strconv"
	"strings"

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

func NewConnectionsDataSource() datasource.DataSource {
	return &ConnectionsDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_connections",
	}
}

type ConnectionsDataSource struct {
	ConnectionsClient ConnectionsClient
	ProviderTypeName  string
	TypeName          string
}

type ConnectionsListDataSourceModel struct {
	Id            types.String                 `tfsdk:"id"`
	EnvironmentId types.String                 `tfsdk:"environment_id"`
	Connections   []ConnectionsDataSourceModel `tfsdk:"connections"`
}

type ConnectionsDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	DisplayName types.String `tfsdk:"display_name"`
	Status      []string     `tfsdk:"status"`
	//Parameters  types.String `tfsdk:"parameters"`
}

func (d *ConnectionsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *ConnectionsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches a list of \"Connections\" for a given environment. Each connection represents an connection instance to an external data source or service.",
		MarkdownDescription: "Fetches a list of [Connection](https://learn.microsoft.com/en-us/power-apps/maker/canvas-apps/add-manage-connections) for a given environment. Each connection represents an connection instance to an external data source or service.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"environment_id": schema.StringAttribute{
				Description:         "Environment Id. The unique identifier of the environment that the connection are associated with.",
				MarkdownDescription: "Environment Id. The unique identifier of the environment that the connection are associated with.",
				Required:            true,
			},
			"connections": schema.ListNestedAttribute{
				Description:         "List of Connections",
				MarkdownDescription: "List of Connections",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique connection id",
							Description:         "Unique connection id",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the connection.",
							Description:         "Name of the connection.",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Display name of the connection.",
							Description:         "Display name of the connection.",
							Computed:            true,
						},
						"status": schema.SetAttribute{
							Description:         "List of connection statuses",
							MarkdownDescription: "List of connection statuses",
							ElementType:         types.StringType,
							Computed:            true,
						},
						// "parameters": schema.StringAttribute{
						// 	Description:         "Connection parameters. Json string containing the authentication connection parameters. Depending on required authentication parameters of a given connector, the connection parameters can vary.",
						// 	MarkdownDescription: "Connection parameters. Json string containing the authentication connection parameters, (for example)[https://learn.microsoft.com/en-us/power-automate/desktop-flows/alm/alm-connection#create-a-connection-using-your-service-principal]. Depending on required authentication parameters of a given connector, the connection parameters can vary.",
						// 	Computed:            true,
						// },
					},
				},
			},
		},
	}
}

func (d *ConnectionsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *ConnectionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ConnectionsListDataSourceModel
	resp.State.Get(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE START: %s", d.ProviderTypeName))

	connections, err := d.ConnectionsClient.GetConnections(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, connection := range connections {
		connectionModel := ConvertFromConnectionDto(connection)
		state.Connections = append(state.Connections, connectionModel)
	}
	state.Id = types.StringValue(strconv.Itoa(len(connections)))

	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func ConvertFromConnectionDto(connection ConnectionDto) ConnectionsDataSourceModel {
	nameConnectorSplit := strings.Split(connection.Properties.ApiId, "/")
	nameConnector := nameConnectorSplit[len(nameConnectorSplit)-1]

	conn := ConnectionsDataSourceModel{
		Id:          types.StringValue(connection.Name),
		Name:        types.StringValue(nameConnector),
		DisplayName: types.StringValue(connection.Properties.DisplayName),
	}

	statuses := []string{}
	for _, status := range connection.Properties.Statuses {
		statuses = append(statuses, status.Status)
	}
	conn.Status = statuses

	// if connection.Properties.ConnectionParametersSet != nil {
	// 	p, _ := json.Marshal(connection.Properties.ConnectionParametersSet)
	// 	conn.Parameters = types.StringValue(string(p))
	// }

	return conn
}
