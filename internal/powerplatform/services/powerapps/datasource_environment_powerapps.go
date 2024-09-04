// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerapps

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/constants"
	"github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var (
	_ datasource.DataSource              = &EnvironmentPowerAppsDataSource{}
	_ datasource.DataSourceWithConfigure = &EnvironmentPowerAppsDataSource{}
)

func NewEnvironmentPowerAppsDataSource() datasource.DataSource {
	return &EnvironmentPowerAppsDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_environment_powerapps",
	}
}

type EnvironmentPowerAppsDataSource struct {
	PowerAppssClient PowerAppssClient
	ProviderTypeName string
	TypeName         string
}

type EnvironmentPowerAppsListDataSourceModel struct {
	Timeouts  timeouts.Value                        `tfsdk:"timeouts"`
	Id        types.String                          `tfsdk:"id"`
	PowerApps []EnvironmentPowerAppsDataSourceModel `tfsdk:"powerapps"`
}

type EnvironmentPowerAppsDataSourceModel struct {
	EnvironmentId types.String `tfsdk:"id"`
	DisplayName   types.String `tfsdk:"display_name"`
	Name          types.String `tfsdk:"name"`
	CreatedTime   types.String `tfsdk:"created_time"`
}

func ConvertFromPowerAppDto(powerAppDto PowerAppBapi) EnvironmentPowerAppsDataSourceModel {
	return EnvironmentPowerAppsDataSourceModel{
		EnvironmentId: types.StringValue(powerAppDto.Properties.Environment.Name),
		DisplayName:   types.StringValue(powerAppDto.Properties.DisplayName),
		Name:          types.StringValue(powerAppDto.Name),
		CreatedTime:   types.StringValue(powerAppDto.Properties.CreatedTime),
	}
}

func (d *EnvironmentPowerAppsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *EnvironmentPowerAppsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of Power Apps in an environment",
		MarkdownDescription: "Fetches the list of Power Apps in an environment.  See [Manage Power Apps](https://learn.microsoft.com/power-platform/admin/admin-manage-apps) for more details about how this data is surfaced in Power Platform Admin Center.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"id": schema.StringAttribute{
				Description:         "Id of the read operation",
				MarkdownDescription: "Id of the read operation",
				Computed:            true,
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

func (d *EnvironmentPowerAppsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientApi := req.ProviderData.(*api.ProviderClient).Api

	if clientApi == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.PowerAppssClient = NewPowerAppssClient(clientApi)
}

func (d *EnvironmentPowerAppsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state EnvironmentPowerAppsListDataSourceModel

	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE ENVIRONMENT POWERAPPS START: %s", d.ProviderTypeName))

	timeout, diags := state.Timeouts.Read(ctx, constants.DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	apps, err := d.PowerAppssClient.GetPowerApps(ctx, "")
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, app := range apps {
		appModel := ConvertFromPowerAppDto(app)
		state.PowerApps = append(state.PowerApps, appModel)
	}

	state.Id = types.StringValue(strconv.Itoa(len(apps)))

	diags = resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE ENVIRONMENT POWERAPPS END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
