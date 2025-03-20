// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application

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
	_ datasource.DataSource              = &TenantApplicationPackagesDataSource{}
	_ datasource.DataSourceWithConfigure = &TenantApplicationPackagesDataSource{}
)

func NewTenantApplicationPackagesDataSource() datasource.DataSource {
	return &TenantApplicationPackagesDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "tenant_application_packages",
		},
	}
}

func (d *TenantApplicationPackagesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *TenantApplicationPackagesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches the list of Dynamics 365 tenant level applications",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the Dynamics 365 application",
				Optional:            true,
			},
			"publisher_name": schema.StringAttribute{
				MarkdownDescription: "Publisher Name of the Dynamics 365 application",
				Optional:            true,
			},
			"applications": schema.ListNestedAttribute{
				MarkdownDescription: "List of Applications",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"application_visibility": schema.StringAttribute{
							MarkdownDescription: "Application Visibility",
							Computed:            true,
						},
						"catalog_visibility": schema.StringAttribute{
							MarkdownDescription: "Catalog Visibility",
							Computed:            true,
						},
						"application_id": schema.StringAttribute{
							MarkdownDescription: "ApplicaitonId",
							Computed:            true,
						},
						"application_descprition": schema.StringAttribute{
							MarkdownDescription: "Applicaiton Description",
							Computed:            true,
						},
						"application_name": schema.StringAttribute{
							MarkdownDescription: "Name",
							Computed:            true,
						},
						"learn_more_url": schema.StringAttribute{
							MarkdownDescription: "Learn More Url",
							Computed:            true,
						},
						"localized_description": schema.StringAttribute{
							MarkdownDescription: "Localized Description",
							Computed:            true,
						},
						"localized_name": schema.StringAttribute{
							MarkdownDescription: "Localized Name",
							Computed:            true,
						},
						"publisher_id": schema.StringAttribute{
							MarkdownDescription: "Publisher Id",
							Computed:            true,
						},
						"publisher_name": schema.StringAttribute{
							MarkdownDescription: "Publisher Name",
							Computed:            true,
						},
						"unique_name": schema.StringAttribute{
							MarkdownDescription: "Unique Name",
							Computed:            true,
						},
						"last_error": schema.ListNestedAttribute{
							MarkdownDescription: "Last Error",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"error_code": schema.StringAttribute{
										MarkdownDescription: "Error Code",
										Computed:            true,
									},
									"error_name": schema.StringAttribute{
										MarkdownDescription: "Error Name",
										Computed:            true,
									},
									"message": schema.StringAttribute{
										MarkdownDescription: "Message",
										Computed:            true,
									},
									"source": schema.StringAttribute{
										MarkdownDescription: "Source",
										Computed:            true,
									},
									"status_code": schema.Int64Attribute{
										MarkdownDescription: "Status Code",
										Computed:            true,
									},
									"type": schema.StringAttribute{
										MarkdownDescription: "Type",
										Computed:            true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *TenantApplicationPackagesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
	d.ApplicationClient = newApplicationClient(client.Api)
}

func (d *TenantApplicationPackagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state TenantApplicationPackagesListDataSourceModel
	resp.State.Get(ctx, &state)

	state.Name = types.StringValue(state.Name.ValueString())
	state.PublisherName = types.StringValue(state.PublisherName.ValueString())

	applications, err := d.ApplicationClient.GetTenantApplications(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, application := range applications {
		if (state.Name.ValueString() != "" && state.Name.ValueString() != application.ApplicationName) ||
			(state.PublisherName.ValueString() != "" && state.PublisherName.ValueString() != application.PublisherName) {
			continue
		}
		app := TenantApplicationPackageDataSourceModel{
			ApplicationId:          types.StringValue(application.ApplicationId),
			ApplicationDescprition: types.StringValue(application.ApplicationDescription),
			Name:                   types.StringValue(application.ApplicationName),
			LearnMoreUrl:           types.StringValue(application.LearnMoreUrl),
			LocalizedDescription:   types.StringValue(application.LocalizedDescription),
			LocalizedName:          types.StringValue(application.LocalizedName),
			PublisherId:            types.StringValue(application.PublisherId),
			PublisherName:          types.StringValue(application.PublisherName),
			UniqueName:             types.StringValue(application.UniqueName),
			ApplicationVisibility:  types.StringValue(application.ApplicationVisibility),
			CatalogVisibility:      types.StringValue(application.CatalogVisibility),
		}
		if application.LastError != nil {
			app.LastError = append(app.LastError, TenantApplicationErrorDetailsDataSourceModel{
				ErrorCode:  types.StringValue(application.LastError.ErrorCode),
				ErrorName:  types.StringValue(application.LastError.ErrorName),
				Message:    types.StringValue(application.LastError.Message),
				Source:     types.StringValue(application.LastError.Source),
				StatusCode: types.Int64Value(application.LastError.StatusCode),
				Type:       types.StringValue(application.LastError.Type),
			})
		}
		state.Applications = append(state.Applications, app)
	}

	diags := resp.State.Set(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
