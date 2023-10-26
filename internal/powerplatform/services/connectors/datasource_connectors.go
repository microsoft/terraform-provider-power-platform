package powerplatform

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	clients "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/clients"
)

var (
	_ datasource.DataSource              = &ConnectorsDataSource{}
	_ datasource.DataSourceWithConfigure = &ConnectorsDataSource{}
)

func NewConnectorsDataSource() datasource.DataSource {
	return &ConnectorsDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_connectors",
	}
}

type ConnectorsDataSource struct {
	ConnectorsClient ConnectorsClient
	ProviderTypeName string
	TypeName         string
}

type ConnectorsListDataSourceModel struct {
	Id         types.String                `tfsdk:"id"`
	Connectors []ConnectorsDataSourceModel `tfsdk:"connectors"`
}

type ConnectorsDataSourceModel struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Type        types.String `tfsdk:"type"`
	Description types.String `tfsdk:"description"`
	DisplayName types.String `tfsdk:"display_name"`
	Tier        types.String `tfsdk:"tier"`
	Publisher   types.String `tfsdk:"publisher"`
	Unblockable types.Bool   `tfsdk:"unblockable"`
}

func ConvertFromConnectorDto(connectorDto ConnectorDto) ConnectorsDataSourceModel {
	return ConnectorsDataSourceModel{
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

func (d *ConnectorsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *ConnectorsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of available connectors in a Power Platform tenant",
		MarkdownDescription: "Fetches the list of available connectors in a Power Platform tenant",
		Attributes: map[string]schema.Attribute{
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

func (d *ConnectorsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client := req.ProviderData.(*clients.ProviderClient).BapiApi.Client

	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.ConnectorsClient = NewConnectorsClient(client)
}

func (d *ConnectorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ConnectorsListDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE CONNECTORS START: %s", d.ProviderTypeName))

	connectors, err := d.ConnectorsClient.GetConnectors(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, connector := range connectors {
		connectorModel := ConvertFromConnectorDto(connector)
		state.Connectors = append(state.Connectors, connectorModel)
	}
	state.Id = types.StringValue(fmt.Sprint((time.Now().Unix())))

	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE CONNECTORS END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
