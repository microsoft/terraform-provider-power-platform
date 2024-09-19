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

type DataSourceModel struct {
	Timeouts timeouts.Value `tfsdk:"timeouts"`
	Id       types.Int64    `tfsdk:"id"`
	Value    []DataModel    `tfsdk:"locations"`
}

type DataModel struct {
	ID                                     string   `tfsdk:"id"`
	Name                                   string   `tfsdk:"name"`
	DisplayName                            string   `tfsdk:"display_name"`
	Code                                   string   `tfsdk:"code"`
	IsDefault                              bool     `tfsdk:"is_default"`
	IsDisabled                             bool     `tfsdk:"is_disabled"`
	CanProvisionDatabase                   bool     `tfsdk:"can_provision_database"`
	CanProvisionCustomerEngagementDatabase bool     `tfsdk:"can_provision_customer_engagement_database"`
	AzureRegions                           []string `tfsdk:"azure_regions"`
}

func NewLocationsDataSource() datasource.DataSource {
	return &DataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "locations",
		},
	}
}

type DataSource struct {
	helpers.TypeInfo
	LocationsClient Client
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

func (d *DataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of available Dynamics 365 locations",
		MarkdownDescription: "Fetches the list of available Dynamics 365 locations. For more information see [Power Platform Geos](https://learn.microsoft.com/power-platform/admin/regions-overview)",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"id": schema.Int64Attribute{
				Description:         "Id of the read operation",
				MarkdownDescription: "Id of the read operation",
				Optional:            true,
			},
			"locations": schema.ListNestedAttribute{
				Description:         "List of available locations",
				MarkdownDescription: "List of available locations",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Description: "Unique identifier of the location",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the location",
							Computed:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "Display name of the location",
							Computed:    true,
						},
						"code": schema.StringAttribute{
							Description: "Code of the location",
							Computed:    true,
						},
						"is_default": schema.BoolAttribute{
							Description: "Is the location default",
							Computed:    true,
						},
						"is_disabled": schema.BoolAttribute{
							Description: "Is the location disabled",
							Computed:    true,
						},
						"can_provision_database": schema.BoolAttribute{
							Description: "Can the location provision a database",
							Computed:    true,
						},
						"can_provision_customer_engagement_database": schema.BoolAttribute{
							Description: "Can the location provision a customer engagement database",
							Computed:    true,
						},
						"azure_regions": schema.ListAttribute{
							Description:         "List of Azure regions",
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
	d.LocationsClient = NewLocationsClient(clientApi)
}

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	var state DataSourceModel
	resp.State.Get(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE LOCATIONS START: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	locations, err := d.LocationsClient.GetLocations(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	state.Id = types.Int64Value(int64(len(locations.Value)))

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
