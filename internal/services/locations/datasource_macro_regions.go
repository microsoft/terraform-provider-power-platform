// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package locations

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
	_ datasource.DataSource              = &macroRegionsDataSource{}
	_ datasource.DataSourceWithConfigure = &macroRegionsDataSource{}
)

func NewMacroRegionsDataSource() datasource.DataSource {
	return &macroRegionsDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "macro_regions",
		},
	}
}

type macroRegionsDataSource struct {
	helpers.TypeInfo
	LocationsClient client
}

func (d *macroRegionsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *macroRegionsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches the macro region geographies available to this tenant. A macro region geography is a data residency boundary; tenants without tenant-wide Advanced Data Residency must create environments with a `macro_region` instead of a `location`. For more information see [macro region geography](https://learn.microsoft.com/power-platform/admin/macro-regions).",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"macro_regions": schema.ListNestedAttribute{
				MarkdownDescription: "Macro region geographies available to this tenant. Empty when the tenant provisions by datacenter location. Use `macro_region_id` as the `macro_region` of a `powerplatform_environment`.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"macro_region_id": schema.StringAttribute{
							MarkdownDescription: "Identifier of the macro region geography (for example `eu-efta`, `north-america`).",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Display name of the macro region geography.",
							Computed:            true,
						},
						"data_residency_note": schema.StringAttribute{
							MarkdownDescription: "Data residency statement shown for this macro region geography. Null when the service does not provide one.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *macroRegionsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
		return
	}

	client, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.LocationsClient = newLocationsClient(client.Api)
}

func (d *macroRegionsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state macroRegionsDataSourceModel
	diags := resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	locations, err := d.LocationsClient.GetLocations(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
		return
	}

	// Tenants that provision by datacenter location return an empty array; that is not an error.
	state.MacroRegions = make([]macroRegionDataModel, 0, len(locations.MacroRegions))
	for _, macroRegion := range locations.MacroRegions {
		dataResidencyNote := types.StringNull()
		if macroRegion.DataResidencyNote != nil {
			dataResidencyNote = types.StringValue(*macroRegion.DataResidencyNote)
		}

		state.MacroRegions = append(state.MacroRegions, macroRegionDataModel{
			MacroRegionId:     types.StringValue(macroRegion.MacroRegionId),
			DisplayName:       types.StringValue(macroRegion.DisplayName),
			DataResidencyNote: dataResidencyNote,
		})
	}

	diags = resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
