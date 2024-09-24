// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connection

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

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

func NewConnectionsDataSource() datasource.DataSource {
	return &ConnectionsDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "connections",
		},
	}
}

func (d *ConnectionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *ConnectionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		Description:         "Fetches a list of \"Connections\" for a given environment. Each connection represents an connection instance to an external data source or service.",
		MarkdownDescription: "Fetches a list of [Connection](https://learn.microsoft.com/en-us/power-apps/maker/canvas-apps/add-manage-connections) for a given environment. Each connection represents an connection instance to an external data source or service.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
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
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the connection.",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Display name of the connection.",
							Computed:            true,
						},
						"status": schema.SetAttribute{
							MarkdownDescription: "List of connection statuses",
							ElementType:         types.StringType,
							Computed:            true,
						},
						"connection_parameters": schema.StringAttribute{
							MarkdownDescription: "Connection parameters. Json string containing the authentication connection parameters (if connection is interactive, leave blank), (for example)[https://learn.microsoft.com/en-us/power-automate/desktop-flows/alm/alm-connection#create-a-connection-using-your-service-principal]. Depending on required authentication parameters of a given connector, the connection parameters can vary.",
							Computed:            true,
						},
						"connection_parameters_set": schema.StringAttribute{
							MarkdownDescription: "Set of connection parameters. Json string containing the authentication connection parameters (if connection is interactive, leave blank), (for example)[https://learn.microsoft.com/en-us/power-automate/desktop-flows/alm/alm-connection#create-a-connection-using-your-service-principal]. Depending on required authentication parameters of a given connector, the connection parameters can vary.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *ConnectionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
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

	d.ConnectionsClient = newConnectionsClient(client)
}

func (d *ConnectionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

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
	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func ConvertFromConnectionDto(connection connectionDto) ConnectionsDataSourceModel {
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

	if connection.Properties.ConnectionParametersSet != nil {
		p, _ := json.Marshal(connection.Properties.ConnectionParametersSet)
		conn.ConnectionParametersSet = types.StringValue(string(p))
	}

	if connection.Properties.ConnectionParameters != nil {
		p, _ := json.Marshal(connection.Properties.ConnectionParameters)
		conn.ConnectionParameters = types.StringValue(string(p))
	}

	return conn
}
