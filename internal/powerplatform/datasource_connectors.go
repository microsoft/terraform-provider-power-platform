package powerplatform

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	powerplatform_bapi "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi"
	models "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/bapi/models"
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
	BapiApiClient    powerplatform_bapi.ApiClientInterface
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
}

func ConvertFromConnectorDto(environmentDto models.ConnectorDto) ConnectorsDataSourceModel {
	return ConnectorsDataSourceModel{
		Id:          types.StringValue(environmentDto.Id),
		Name:        types.StringValue(environmentDto.Name),
		Type:        types.StringValue(environmentDto.Type),
		Description: types.StringValue(environmentDto.Properties.Description),
		DisplayName: types.StringValue(environmentDto.Properties.DisplayName),
		Tier:        types.StringValue(environmentDto.Properties.Tier),
		Publisher:   types.StringValue(environmentDto.Properties.Publisher),
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

	client, ok := req.ProviderData.(*PowerPlatformProvider).bapiClient.(powerplatform_bapi.ApiClientInterface)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.BapiApiClient = client
}

func (d *ConnectorsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state ConnectorsListDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE CONNECTORS START: %s", d.ProviderTypeName))

	connectors, err := d.BapiApiClient.GetConnectors(ctx)
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
