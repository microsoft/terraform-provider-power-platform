// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform_analytics_data_export

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var (
	_ datasource.DataSource              = &AnalyticsExportDataSource{}
	_ datasource.DataSourceWithConfigure = &AnalyticsExportDataSource{}
)

type AnalyticsExportDataSource struct {
	ProviderTypeName    string
	TypeInfo            helpers.TypeInfo
	AnalyticsExportData Client
}

func NewAnalyticsExportDataSource() datasource.DataSource {
	return &AnalyticsExportDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "Analytics_DataExport",
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

func (d *AnalyticsExportDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		MarkdownDescription: "Analytics Data Export configuration",
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
						"environments": schema.SetNestedAttribute{
							MarkdownDescription: "The environments configured for analytics data export",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"environment_id": schema.StringAttribute{
										MarkdownDescription: "The identifier of the environment",
										Computed:            true,
									},
									"organization_id": schema.StringAttribute{
										MarkdownDescription: "The identifier of the organization",
										Computed:            true,
									},
								},
							},
						},
						"status": schema.SetNestedAttribute{
							MarkdownDescription: "The status information for the analytics data export",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
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
								},
							},
						},
						"sink": schema.SingleNestedAttribute{
							MarkdownDescription: "The sink configuration for analytics data",
							Required:            true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "The resource ID of the sink",
									Computed:            true,
								},
								"type": schema.StringAttribute{
									MarkdownDescription: "The type of the sink",
									Required:            true,
								},
								"subscription_id": schema.StringAttribute{
									MarkdownDescription: "The Azure subscription ID",
									Optional:            true,
								},
								"resource_group_name": schema.StringAttribute{
									MarkdownDescription: "The Azure resource group name",
									Optional:            true,
								},
								"resource_name": schema.StringAttribute{
									MarkdownDescription: "The name of the sink resource",
									Computed:            true,
								},
								"key": schema.StringAttribute{
									MarkdownDescription: "The key for accessing the sink",
									Computed:            true,
								},
							},
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
							Required:            true,
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

	// Here we should set d.AnalyticsExportData from the provider client (missing)
	// e.g., client, ok := req.ProviderData.(*api.ProviderClient); if ok { d.AnalyticsExportData = client.AnalyticsExportData }
	tflog.Debug(ctx, "CONFIGURE: completed")
}

func (d *AnalyticsExportDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AnalyticsExportDatasourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Here state.ID is referenced but not defined in the schema. In common implementations,
	// either the schema defines an input ID or the data source uses a list/filter.
	id := state.ID.ValueString()

	apiUrl := &url.URL{
		Scheme: "https",
		Host:   getAnalyticsUrlMap()[""], // <-- empty key is unusual
		Path:   fmt.Sprintf("/api/v2/analyticsdataexport/%s", id),
	}

	var analyticsDataExport AnalyticsExportDatasourceModel
	_, err := d.AnalyticsExportData.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &analyticsDataExport)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Analytics Data Export", fmt.Sprintf("Unable to read analytics data export with ID %s: %s", id, err.Error()))
		return
	}

	// Set state
	diags = resp.State.Set(ctx, &analyticsDataExport)
	resp.Diagnostics.Append(diags...)
}
