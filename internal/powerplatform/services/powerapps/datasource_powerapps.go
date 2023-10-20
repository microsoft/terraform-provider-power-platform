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
	_ datasource.DataSource              = &PowerAppsDataSource{}
	_ datasource.DataSourceWithConfigure = &PowerAppsDataSource{}
)

func NewPowerAppsDataSource() datasource.DataSource {
	return &PowerAppsDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_powerapps",
	}
}

type PowerAppsDataSource struct {
	PowerAppssClient PowerAppssClient
	ProviderTypeName string
	TypeName         string
}

type PowerAppsListDataSourceModel struct {
	Id        types.String               `tfsdk:"id"`
	PowerApps []PowerAppsDataSourceModel `tfsdk:"powerapps"`
}

type PowerAppsDataSourceModel struct {
	EnvironmentId types.String `tfsdk:"id"`
	DisplayName   types.String `tfsdk:"display_name"`
	Name          types.String `tfsdk:"name"`
	CreatedTime   types.String `tfsdk:"created_time"`
}

func ConvertFromPowerAppDto(powerAppDto PowerAppBapi) PowerAppsDataSourceModel {
	return PowerAppsDataSourceModel{
		EnvironmentId: types.StringValue(powerAppDto.Properties.Environment.Name),
		DisplayName:   types.StringValue(powerAppDto.Properties.DisplayName),
		Name:          types.StringValue(powerAppDto.Name),
		CreatedTime:   types.StringValue(powerAppDto.Properties.CreatedTime),
	}
}

func (d *PowerAppsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *PowerAppsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of Power Apps in a tenant",
		MarkdownDescription: "Fetches the list of Power Apps in a tenant",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"powerapps": schema.ListNestedAttribute{
				Description:         "List of Power Apps",
				MarkdownDescription: "List of Power Apps",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique environment id (guid)",
							Description:         "Unique environment id (guid)",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Display name",
							Description:         "Display name",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name",
							Description:         "Name",
							Computed:            true,
						},
						"created_time": schema.StringAttribute{
							MarkdownDescription: "Created time",
							Description:         "Created time",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *PowerAppsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientBapi := req.ProviderData.(*clients.ProviderClient).BapiApi.Client
	clientDv := req.ProviderData.(*clients.ProviderClient).DataverseApi.Client

	if clientBapi == nil || clientDv == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.PowerAppssClient = NewPowerAppssClient(clientBapi, clientDv)
}

func (d *PowerAppsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state PowerAppsListDataSourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE POWERAPPS START: %s", d.ProviderTypeName))

	apps, err := d.PowerAppssClient.GetPowerApps(ctx, "")
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, app := range apps {
		appModel := ConvertFromPowerAppDto(app)
		state.PowerApps = append(state.PowerApps, appModel)
	}

	state.Id = types.StringValue(fmt.Sprint((time.Now().Unix())))

	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE POWERAPPS END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
