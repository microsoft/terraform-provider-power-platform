// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package analytics_data_export

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
	"github.com/microsoft/terraform-provider-power-platform/internal/services/tenant"
)

var (
	_ datasource.DataSource              = &AnalyticsExportDataSource{}
	_ datasource.DataSourceWithConfigure = &AnalyticsExportDataSource{}
)

type AnalyticsExportDataSource struct {
	helpers.TypeInfo
	analyticsExportClient Client
}

func NewAnalyticsExportDataSource() datasource.DataSource {
	return &AnalyticsExportDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "analytics_data_exports",
		},
	}
}

func (d *AnalyticsExportDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Set the proper resource type (common pattern)
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func statusNestedAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the status entry",
			Computed:            true,
		},
		"state": schema.StringAttribute{
			MarkdownDescription: "The current state of the analytics component",
			Computed:            true,
		},
		"last_run_on": schema.StringAttribute{
			MarkdownDescription: "The timestamp of the last execution",
			Computed:            true,
		},
		"message": schema.StringAttribute{
			MarkdownDescription: "Any message associated with the status",
			Computed:            true,
		},
	}
}

func sinkNestedAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The resource ID of the sink",
			Computed:            true,
		},
		"type": schema.StringAttribute{
			MarkdownDescription: "The type of the sink",
			Computed:            true,
		},
		"subscription_id": schema.StringAttribute{
			MarkdownDescription: "The Azure subscription ID",
			Computed:            true,
		},
		"resource_group_name": schema.StringAttribute{
			MarkdownDescription: "The Azure resource group name",
			Computed:            true,
		},
		"resource_name": schema.StringAttribute{
			MarkdownDescription: "The name of the sink resource",
			Computed:            true,
		},
		"key": schema.StringAttribute{
			MarkdownDescription: "The key for accessing the sink",
			Computed:            true,
		},
	}
}

func (d *AnalyticsExportDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		MarkdownDescription: "Analytics Data Export configurations. See [documentation](https://learn.microsoft.com/en-us/power-platform/admin/set-up-export-application-insights) for more details.\n\n" +
			"**Note:** This resource is available as **preview**\n\n" +
			"**Known Limitations:** This resource is not supported for with service principal authentication.",
		Attributes: map[string]schema.Attribute{
			"exports": schema.SetNestedAttribute{
				MarkdownDescription: "The collection of analytics data exports",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "The unique identifier of the analytics data export",
							Computed:            true,
						},
						"source": schema.StringAttribute{
							MarkdownDescription: "The source of the analytics data",
							Computed:            true,
						},
						"environments": schema.SetAttribute{
							MarkdownDescription: "The environment IDs configured for analytics data export",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"status": schema.SetNestedAttribute{
							MarkdownDescription: "The status information for the analytics data export",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: statusNestedAttributes(),
							},
						},
						"sink": schema.SingleNestedAttribute{
							MarkdownDescription: "The sink configuration for analytics data",
							Computed:            true,
							Attributes:          sinkNestedAttributes(),
						},
						"package_name": schema.StringAttribute{
							MarkdownDescription: "The package name for the analytics data",
							Computed:            true,
						},
						"scenarios": schema.ListAttribute{
							MarkdownDescription: "The list of scenarios covered by this analytics export",
							Computed:            true,
							ElementType:         types.StringType,
						},
						"resource_provider": schema.StringAttribute{
							MarkdownDescription: "The resource provider for the analytics data",
							Computed:            true,
						},
						"ai_type": schema.StringAttribute{
							MarkdownDescription: "The AI type for the analytics data",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *AnalyticsExportDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig. It's ok.
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

	d.analyticsExportClient = NewAnalyticsExportClient(client.Api, tenant.NewTenantClient(client.Api))
}

func (d *AnalyticsExportDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var config AnalyticsDataExportModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// Fetch the analytics data export
	analyticsDataExport, err := d.analyticsExportClient.GetAnalyticsDataExport(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching analytics data export",
			fmt.Sprintf("Unable to fetch analytics data export: %s", err),
		)
		return
	}
	if analyticsDataExport == nil {
		resp.Diagnostics.AddError(
			"Analytics data export not found",
			"Unable to find analytics data export with the specified ID",
		)
		return
	}

	// Map the response to the model
	// Map each analytics data export item to a model
	exports := make([]AnalyticsDataModel, 0, len(*analyticsDataExport))
	for _, export := range *analyticsDataExport {
		if model := convertDtoToModel(&export); model != nil {
			exports = append(exports, *model)
		}
	}

	model := &AnalyticsDataExportModel{
		Exports: exports,
	}

	// Set state
	diags = resp.State.Set(ctx, model)
	resp.Diagnostics.Append(diags...)
}
