// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package application

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var (
	_ datasource.DataSource              = &TenantApplicationPackagesDataSource{}
	_ datasource.DataSourceWithConfigure = &TenantApplicationPackagesDataSource{}
)

func NewTenantApplicationPackagesDataSource() datasource.DataSource {
	return &TenantApplicationPackagesDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_tenant_application_packages",
	}
}

type TenantApplicationPackagesDataSource struct {
	ApplicationClient ApplicationClient
	ProviderTypeName  string
	TypeName          string
}

type TenantApplicationPackagesListDataSourceModel struct {
	Name          types.String                              `tfsdk:"name"`
	PublisherName types.String                              `tfsdk:"publisher_name"`
	Id            types.String                              `tfsdk:"id"`
	Applications  []TenantApplicationPackageDataSourceModel `tfsdk:"applications"`
}

type TenantApplicationPackageDataSourceModel struct {
	ApplicationId          types.String                                   `tfsdk:"application_id"`
	ApplicationDescprition types.String                                   `tfsdk:"application_descprition"`
	Name                   types.String                                   `tfsdk:"application_name"`
	LearnMoreUrl           types.String                                   `tfsdk:"learn_more_url"`
	LocalizedDescription   types.String                                   `tfsdk:"localized_description"`
	LocalizedName          types.String                                   `tfsdk:"localized_name"`
	PublisherId            types.String                                   `tfsdk:"publisher_id"`
	PublisherName          types.String                                   `tfsdk:"publisher_name"`
	UniqueName             types.String                                   `tfsdk:"unique_name"`
	ApplicationVisibility  types.String                                   `tfsdk:"application_visibility"`
	CatalogVisibility      types.String                                   `tfsdk:"catalog_visibility"`
	LastError              []TenantApplicationErrorDetailsDataSourceModel `tfsdk:"last_error"`
}

type TenantApplicationErrorDetailsDataSourceModel struct {
	ErrorCode  types.String `tfsdk:"error_code"`
	ErrorName  types.String `tfsdk:"error_name"`
	Message    types.String `tfsdk:"message"`
	Source     types.String `tfsdk:"source"`
	StatusCode types.Int64  `tfsdk:"status_code"`
	Type       types.String `tfsdk:"type"`
}

func (d *TenantApplicationPackagesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *TenantApplicationPackagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of Dynamics 365 tenant level applications",
		MarkdownDescription: "Fetches the list of Dynamics 365 tenant level applications",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Id of the read operation",
				Optional:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the Dynamics 365 application",
				Optional:    true,
			},
			"publisher_name": schema.StringAttribute{
				Description: "Publisher Name of the Dynamics 365 application",
				Optional:    true,
			},
			"applications": schema.ListNestedAttribute{
				Description:         "List of Applications",
				MarkdownDescription: "List of Applications",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"application_visibility": schema.StringAttribute{
							MarkdownDescription: "Application Visibility",
							Description:         "Application Visibility",
							Computed:            true,
						},
						"catalog_visibility": schema.StringAttribute{
							MarkdownDescription: "Catalog Visibility",
							Description:         "Catalog Visibility",
							Computed:            true,
						},
						"application_id": schema.StringAttribute{
							MarkdownDescription: "ApplicaitonId",
							Description:         "ApplicaitonId",
							Computed:            true,
						},
						"application_descprition": schema.StringAttribute{
							MarkdownDescription: "Applicaiton Description",
							Description:         "Applicaiton Description",
							Computed:            true,
						},
						"application_name": schema.StringAttribute{
							MarkdownDescription: "Name",
							Description:         "Name",
							Computed:            true,
						},
						"learn_more_url": schema.StringAttribute{
							MarkdownDescription: "Learn More Url",
							Description:         "Learn More Url",
							Computed:            true,
						},
						"localized_description": schema.StringAttribute{
							MarkdownDescription: "Localized Description",
							Description:         "Localized Description",
							Computed:            true,
						},
						"localized_name": schema.StringAttribute{
							MarkdownDescription: "Localized Name",
							Description:         "Localized Name",
							Computed:            true,
						},
						"publisher_id": schema.StringAttribute{
							MarkdownDescription: "Publisher Id",
							Description:         "Publisher Id",
							Computed:            true,
						},
						"publisher_name": schema.StringAttribute{
							MarkdownDescription: "Publisher Name",
							Description:         "Publisher Name",
							Computed:            true,
						},
						"unique_name": schema.StringAttribute{
							MarkdownDescription: "Unique Name",
							Description:         "Unique Name",
							Computed:            true,
						},
						"last_error": schema.ListNestedAttribute{
							Description:         "Last Error",
							MarkdownDescription: "Last Error",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"error_code": schema.StringAttribute{
										MarkdownDescription: "Error Code",
										Description:         "Error Code",
										Computed:            true,
									},
									"error_name": schema.StringAttribute{
										MarkdownDescription: "Error Name",
										Description:         "Error Name",
										Computed:            true,
									},
									"message": schema.StringAttribute{
										MarkdownDescription: "Message",
										Description:         "Message",
										Computed:            true,
									},
									"source": schema.StringAttribute{
										MarkdownDescription: "Source",
										Description:         "Source",
										Computed:            true,
									},
									"status_code": schema.Int64Attribute{
										MarkdownDescription: "Status Code",
										Description:         "Status Code",
										Computed:            true,
									},
									"type": schema.StringAttribute{
										MarkdownDescription: "Type",
										Description:         "Type",
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
	d.ApplicationClient = NewApplicationClient(clientApi)
}

func (d *TenantApplicationPackagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan TenantApplicationPackagesListDataSourceModel
	resp.State.Get(ctx, &plan)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE TENANT APPLICATION PACKAGES START: %s", d.ProviderTypeName))

	plan.Name = types.StringValue(plan.Name.ValueString())
	plan.PublisherName = types.StringValue(plan.PublisherName.ValueString())

	applications, err := d.ApplicationClient.GetTenantApplications(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	plan.Id = types.StringValue(strconv.Itoa(len(applications)))

	for _, application := range applications {
		if (plan.Name.ValueString() != "" && plan.Name.ValueString() != application.ApplicationName) ||
			(plan.PublisherName.ValueString() != "" && plan.PublisherName.ValueString() != application.PublisherName) {
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
		plan.Applications = append(plan.Applications, app)
	}

	diags := resp.State.Set(ctx, &plan)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE TENANT APPLICATION PACKAGES END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
