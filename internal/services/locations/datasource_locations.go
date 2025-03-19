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
	_ datasource.DataSource              = &DataSource{}
	_ datasource.DataSourceWithConfigure = &DataSource{}
)

func NewLocationsDataSource() datasource.DataSource {
	return &DataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "locations",
		},
	}
}

type DataSource struct {
	helpers.TypeInfo
	LocationsClient client
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches the list of available Dynamics 365 locations. For more information see [Power Platform Geos](https://learn.microsoft.com/power-platform/admin/regions-overview)",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"locations": schema.ListNestedAttribute{
				MarkdownDescription: "List of available locations",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Unique identifier of the location",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the location",
							Computed:            true,
						},
						"display_name": schema.StringAttribute{
							MarkdownDescription: "Display name of the location",
							Computed:            true,
						},
						"code": schema.StringAttribute{
							MarkdownDescription: "Code of the location",
							Computed:            true,
						},
						"is_default": schema.BoolAttribute{
							MarkdownDescription: "Is the location default",
							Computed:            true,
						},
						"is_disabled": schema.BoolAttribute{
							MarkdownDescription: "Is the location disabled",
							Computed:            true,
						},
						"can_provision_database": schema.BoolAttribute{
							MarkdownDescription: "Can the location provision a database",
							Computed:            true,
						},
						"can_provision_customer_engagement_database": schema.BoolAttribute{
							MarkdownDescription: "Can the location provision a customer engagement database",
							Computed:            true,
						},
						"azure_regions": schema.ListAttribute{
							MarkdownDescription: "List of Azure regions",
							Computed:            true,
							ElementType:         types.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state DataSourceModel
	resp.State.Get(ctx, &state)
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	locations, err := d.LocationsClient.GetLocations(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, location := range locations.Value {
		state.Value = append(state.Value, DataModel{
			ID:                                     location.ID,
			Name:                                   location.Name,
			DisplayName:                            location.Properties.DisplayName,
			Code:                                   location.Properties.Code,
			IsDefault:                              location.Properties.IsDefault,
			IsDisabled:                             location.Properties.IsDisabled,
			CanProvisionDatabase:                   location.Properties.CanProvisionDatabase,
			CanProvisionCustomerEngagementDatabase: location.Properties.CanProvisionCustomerEngagementDatabase,
			AzureRegions:                           location.Properties.AzureRegions,
		})
	}

	diags := resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
