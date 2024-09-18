// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package connectors

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
	_ datasource.DataSource              = &DataSource{}
	_ datasource.DataSourceWithConfigure = &DataSource{}
)

func NewConnectorsDataSource() datasource.DataSource {
	return &DataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_connectors",
	}
}

type DataSource struct {
	ConnectorsClient Client
	ProviderTypeName string
	TypeName         string
}

type ListDataSourceModel struct {
	Timeouts   timeouts.Value    `tfsdk:"timeouts"`
	Id         types.String      `tfsdk:"id"`
	Connectors []DataSourceModel `tfsdk:"connectors"`
}

type DataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
	DisplayName types.String `tfsdk:"display_name"`
	Tier        types.String `tfsdk:"tier"`
	Publisher   types.String `tfsdk:"publisher"`
	Unblockable types.Bool   `tfsdk:"unblockable"`
}

func ConvertFromConnectorDto(connectorDto Dto) DataSourceModel {
	return DataSourceModel{
		Id:          types.StringValue(connectorDto.Id),
		Name:        types.StringValue(connectorDto.Name),
		Type:        types.StringValue(connectorDto.Type),
		Description: types.StringValue(connectorDto.Properties.Description),
		DisplayName: types.StringValue(connectorDto.Properties.DisplayName),
		Tier:        types.StringValue(connectorDto.Properties.Tier),
		Publisher:   types.StringValue(connectorDto.Properties.Publisher),
		Unblockable: types.BoolValue(connectorDto.Properties.Unblockable),
	}
}

func (d *DataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *DataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of available connectors within a specific Power Platform tenant. Each connector represents a service that can be used to enhance the capabilities of Power Apps, Power Automate, and Power Virtual Agents. The returned list includes both standard and custom connectors, providing a comprehensive view of the services that can be integrated into your Power Platform solutions. The list can be used to understand what services are readily available for use within your tenant, and can assist in planning and developing new applications or flows. It's important to note that the availability of connectors may vary based on the specific licenses and permissions assigned within your tenant.",
		MarkdownDescription: "Fetches the list of available connectors within a specific Power Platform tenant. Each connector represents a service that can be used to enhance the capabilities of Power Apps, Power Automate, and Power Virtual Agents. The returned list includes both standard and custom connectors, providing a comprehensive view of the services that can be integrated into your Power Platform solutions. The list can be used to understand what services are readily available for use within your tenant, and can assist in planning and developing new applications or flows. It's important to note that the availability of connectors may vary based on the specific licenses and permissions assigned within your tenant.\n\nAdditional Resources:\n\n* [Connectors Overview](https://learn.microsoft.com/connectors/connectors)\n",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"id": schema.StringAttribute{
				Computed: true,
			},
			"connectors": schema.ListNestedAttribute{
				Description:         "List of Connectors",
				MarkdownDescription: "List of Connectors",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Id",
							Description:         "id",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name",
							Description:         "Name",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Display name",
							Description:         "Display name",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Type",
							Description:         "Type",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description",
							Description:         "Description",
							Computed:            true,
						},
						"tier": schema.StringAttribute{
							MarkdownDescription: "Tier",
							Description:         "Tier",
							Computed:            true,
						},
						"publisher": schema.StringAttribute{
							MarkdownDescription: "Publisher",
							Description:         "Publisher",
							Computed:            true,
						},
						"unblockable": schema.BoolAttribute{
							MarkdownDescription: "Indicates if the connector can be blocked in a Data Loss Prevention policy. If true, the connector has to be in 'Non-Business' connectors group.",
							Description:         "Indicates if the connector can be blocked in a Data Loss Prevention policy. If true, the connector has to be in 'Non-Business' connectors group.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *DataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.ConnectorsClient = NewConnectorsClient(client)
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ListDataSourceModel
	resp.State.Get(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE CONNECTORS START: %s", d.ProviderTypeName))

	timeout, diags := state.Timeouts.Read(ctx, constants.DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	connectors, err := d.ConnectorsClient.GetConnectors(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, connector := range connectors {
		connectorModel := ConvertFromConnectorDto(connector)
		state.Connectors = append(state.Connectors, connectorModel)
	}
	state.Id = types.StringValue(strconv.Itoa(len(connectors)))

	diags = resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE CONNECTORS END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
