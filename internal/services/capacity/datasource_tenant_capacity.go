// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package capacity

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
)

var (
	_ datasource.DataSource              = &TenantCapacityDataSource{}
	_ datasource.DataSourceWithConfigure = &TenantCapacityDataSource{}
)

func NewTenantCapcityDataSource() datasource.DataSource {
	return &TenantCapacityDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_tenant_capacity",
	}
}

type TenantCapacityDataSource struct {
	CapacityClient   Client
	ProviderTypeName string
	TypeName         string
}

type TenantCapacityDataSourceModel struct {
	TenantId         types.String              `tfsdk:"tenant_id"`
	LicenseModelType types.String              `tfsdk:"license_model_type"`
	TenantCapacities []CapacityDataSourceModel `tfsdk:"tenant_capacities"`
}

type CapacityDataSourceModel struct {
	CapacityType  types.String               `tfsdk:"capacity_type"`
	CapacityUnits types.String               `tfsdk:"capacity_units"`
	TotalCapacity types.Float32              `tfsdk:"total_capacity"`
	MaxCapacity   types.Float32              `tfsdk:"max_capacity"`
	Consumption   ConsumptionDataSourceModel `tfsdk:"consumption"`
	Status        types.String               `tfsdk:"status"`
}

type ConsumptionDataSourceModel struct {
	Actual          types.Float32 `tfsdk:"actual"`
	Rated           types.Float32 `tfsdk:"rated"`
	ActualUpdatedOn types.String  `tfsdk:"actual_updated_on"`
	RatedUpdatedOn  types.String  `tfsdk:"rated_updated_on"`
}

func (d *TenantCapacityDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *TenantCapacityDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the capacity information for a given tenant.",
		MarkdownDescription: "Fetches the capacity information for a given tenant.",
		Attributes: map[string]schema.Attribute{
			"tenant_id": schema.StringAttribute{
				Required:    true,
				Description: "The tenant ID for which the capacity information is to be fetched.",
			},
			"license_model_type": schema.StringAttribute{
				Computed:    true,
				Description: "The license model type for which the capacity information is to be fetched.",
			},
			"tenant_capacities": schema.ListNestedAttribute{
				Description: "The list of capacities for the given tenant.",
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"capacity_type": schema.StringAttribute{
							Computed:    true,
							Description: "The type of the capacity.",
						},
						"capacity_units": schema.StringAttribute{
							Computed:    true,
							Description: "The units of the capacity.",
						},
						"total_capacity": schema.Float32Attribute{
							Computed:    true,
							Description: "The total capacity.",
						},
						"max_capacity": schema.Float32Attribute{
							Computed:    true,
							Description: "The maximum capacity.",
						},
						"consumption": schema.SingleNestedAttribute{
							Description: "The consumption details.",
							Computed:    true,
							Attributes: map[string]schema.Attribute{
								"actual": schema.Float32Attribute{
									Computed:    true,
									Description: "The actual consumption.",
								},
								"rated": schema.Float32Attribute{
									Computed:    true,
									Description: "The rated consumption.",
								},
								"actual_updated_on": schema.StringAttribute{
									Computed:    true,
									Description: "The actual consumption updated on.",
								},
								"rated_updated_on": schema.StringAttribute{
									Computed:    true,
									Description: "The rated consumption updated on.",
								},
							},
						},
						"status": schema.StringAttribute{
							Computed:    true,
							Description: "The status of the capacity.",
						},
					},
				},
			},
		},
	}
}

func (d *TenantCapacityDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.CapacityClient = NewCapacityClient(client)
}

func (d *TenantCapacityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state TenantCapacityDataSourceModel
	var tenantId string
	req.Config.GetAttribute(ctx, path.Root("tenant_id"), &tenantId)

	tenantCapacityDto, err := d.CapacityClient.GetTenantCapacity(ctx, tenantId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching tenant capacity",
			fmt.Sprintf("Error fetching tenant capacity: %v", err),
		)
		return
	}

	state.TenantId = types.StringValue(tenantCapacityDto.TenantId)
	state.LicenseModelType = types.StringValue(tenantCapacityDto.LicenseModelType)

	for _, capacity := range tenantCapacityDto.TenantCapacities {
		state.TenantCapacities = append(state.TenantCapacities, CapacityDataSourceModel{
			CapacityType:  types.StringValue(capacity.CapacityType),
			CapacityUnits: types.StringValue(capacity.CapacityUnits),
			TotalCapacity: types.Float32Value(capacity.TotalCapacity),
			MaxCapacity:   types.Float32Value(capacity.MaxCapacity),
			Consumption: ConsumptionDataSourceModel{
				Actual:          types.Float32Value(capacity.Consumption.Actual),
				Rated:           types.Float32Value(capacity.Consumption.Rated),
				ActualUpdatedOn: types.StringValue(capacity.Consumption.ActualUpdatedOn),
				RatedUpdatedOn:  types.StringValue(capacity.Consumption.RatedUpdatedOn),
			},
			Status: types.StringValue(capacity.Status),
		})
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
