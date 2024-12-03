// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform_analytics_data_export

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var (
	_ datasource.DataSource              = &AnalyticsExportDataSource{}
	_ datasource.DataSourceWithConfigure = &AnalyticsExportDataSource{}
)

func NewAnalyticsExportDataSource() datasource.DataSource {
	return &AnalyticsExportDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "Analytics_DataExport",
		},
	}
}

func (d *AnalyticsExportDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *AnalyticsExportDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	defer exitContext()
	resp.Schema = schema.Schema{
		Description: "Data source to retrieve Analytics Data Export information.",
		MarkdownDescription: "Data source to retrieve Analytics Data Export information.\n\nThis functionality allows you to connect your Power Platform org to App Insights instance and export the telemetry.\n\n* [Export Data to Application Insights](https://learn.microsoft.com/en-us/power-platform/admin/set-up-export-application-insights)",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: true,
				Update: true,
				Delete: true,
				Read:   true,
			}),			
            "id": schema.StringAttribute{
                MarkdownDescription: "Unique ID of the Analytics Data Export.",
                Required:            true,
            },
            "source": schema.StringAttribute{
                MarkdownDescription: "Source of the Analytics Data Export.",
                Computed:            true,
            },
            "ai_type": schema.StringAttribute{
                MarkdownDescription: "Type of AI for the Analytics Data Export.",
                Computed:            true,
            },
            "environments": schema.ListNestedAttribute{
                MarkdownDescription: "List of environments associated with the Analytics Data Export.",
                Computed:            true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "environment_id": schema.StringAttribute{
                            MarkdownDescription: "ID of the environment.",
                            Computed:            true,
                        },
                    },
                }
            },
            "package_name": schema.StringAttribute{
                MarkdownDescription: "Package name of the Analytics Data Export.",
                Computed:            true,
            },
            "resource_provider": schema.StringAttribute{
                MarkdownDescription: "Resource provider of the Analytics Data Export.",
                Computed:            true,
            },
            "scenarios": schema.ListAttribute{
                MarkdownDescription: "List of scenarios for the Analytics Data Export.",
                Computed:            true,
                ElementType:         types.StringType,
            },
            "sink": schema.SingleNestedAttribute{
                MarkdownDescription: "Sink information for the Analytics Data Export.",
                Computed:            true,
                Attributes: map[string]schema.Attribute{
                    "key": schema.StringAttribute{
                        MarkdownDescription: "Key of the sink.",
                        Computed:            true,
					},
                    "type": schema.StringAttribute{
                        MarkdownDescription: "Type of the sink.",
                        Computed:            true,
                    },
                    // Add other sink attributes as necessary
                },
            },
            "status": schema.ListNestedAttribute{
                MarkdownDescription: "Status information for the Analytics Data Export.",
                Computed:            true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "name": schema.StringAttribute{
                            MarkdownDescription: "Name of the status.",
                            Computed:            true,
                        },
                        "state": schema.StringAttribute{
                            MarkdownDescription: "State of the status.",
                            Computed:            true,
                        },
                        "last_run_on": schema.StringAttribute{
                            MarkdownDescription: "Last run time of the status.",
                            Computed:            true,
                        },
                        // Add other status attributes as necessary			
					},
				},
		
			},
		},
	}
}

func (d *AnalyticsExportDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	tflog.Debug(ctx, "CONFIGURE: completed")
}

func (d *AnalyticsExportDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state AnalyticsExportDatasourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiUrl := &url.URL{
		Scheme: "https",
        Host:   d.GetConfig().Urls.PowerPlatformAnalyticsUrl,
		Path:   "/api/v2/analyticsdataexport/",
	}

	analyticsDataExport := &AnalyticsExportDataSource{}
	_, err := d.ApiClient.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, analyticsDataExport)
	if err != nil {
		resp.Diagnostics.AddError("Error reading Analytics Data Export", err.Error())
		return
	}

	diags = resp.State.Set(ctx, analyticsDataExport)
	resp.Diagnostics.Append(diags...)
}
