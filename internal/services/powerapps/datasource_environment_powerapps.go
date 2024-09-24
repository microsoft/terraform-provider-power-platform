// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerapps

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
	_ datasource.DataSource              = &EnvironmentPowerAppsDataSource{}
	_ datasource.DataSourceWithConfigure = &EnvironmentPowerAppsDataSource{}
)

func NewEnvironmentPowerAppsDataSource() datasource.DataSource {
	return &EnvironmentPowerAppsDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_powerapps",
		},
	}
}

type EnvironmentPowerAppsDataSource struct {
	helpers.TypeInfo
	PowerAppssClient client
}

type EnvironmentPowerAppsListDataSourceModel struct {
	Timeouts  timeouts.Value                        `tfsdk:"timeouts"`
	PowerApps []EnvironmentPowerAppsDataSourceModel `tfsdk:"powerapps"`
}

type EnvironmentPowerAppsDataSourceModel struct {
	EnvironmentId types.String `tfsdk:"id"`
	DisplayName   types.String `tfsdk:"display_name"`
	Name          types.String `tfsdk:"name"`
	CreatedTime   types.String `tfsdk:"created_time"`
}

func ConvertFromPowerAppDto(powerAppDto powerAppBapiDto) EnvironmentPowerAppsDataSourceModel {
	return EnvironmentPowerAppsDataSourceModel{
		EnvironmentId: types.StringValue(powerAppDto.Properties.Environment.Name),
		DisplayName:   types.StringValue(powerAppDto.Properties.DisplayName),
		Name:          types.StringValue(powerAppDto.Name),
		CreatedTime:   types.StringValue(powerAppDto.Properties.CreatedTime),
	}
}

func (d *EnvironmentPowerAppsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *EnvironmentPowerAppsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of Power Apps in an environment",
		MarkdownDescription: "Fetches the list of Power Apps in an environment.  See [Manage Power Apps](https://learn.microsoft.com/power-platform/admin/admin-manage-apps) for more details about how this data is surfaced in Power Platform Admin Center.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
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

func (d *EnvironmentPowerAppsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
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

	d.PowerAppssClient = newPowerAppssClient(clientApi)
}

func (d *EnvironmentPowerAppsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state EnvironmentPowerAppsListDataSourceModel

	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apps, err := d.PowerAppssClient.GetPowerApps(ctx, "")
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, app := range apps {
		appModel := ConvertFromPowerAppDto(app)
		state.PowerApps = append(state.PowerApps, appModel)
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
