// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var (
	_ datasource.DataSource              = &EnvironmentApplicationPackagesDataSource{}
	_ datasource.DataSourceWithConfigure = &EnvironmentApplicationPackagesDataSource{}
)

func NewEnvironmentApplicationPackagesDataSource() datasource.DataSource {
	return &EnvironmentApplicationPackagesDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_environment_application_packages",
	}
}

type EnvironmentApplicationPackagesDataSource struct {
	ApplicationClient ApplicationClient
	ProviderTypeName  string
	TypeName          string
}

type EnvironmentApplicationPackagesListDataSourceModel struct {
	EnvironmentId types.String                                   `tfsdk:"environment_id"`
	Name          types.String                                   `tfsdk:"name"`
	PublisherName types.String                                   `tfsdk:"publisher_name"`
	Id            types.String                                   `tfsdk:"id"`
	Applications  []EnvironmentApplicationPackageDataSourceModel `tfsdk:"applications"`
}

type EnvironmentApplicationPackageDataSourceModel struct {
	ApplicationId         types.String `tfsdk:"application_id"`
	Name                  types.String `tfsdk:"application_name"`
	UniqueName            types.String `tfsdk:"unique_name"`
	Version               types.String `tfsdk:"version"`
	Description           types.String `tfsdk:"description"`
	PublisherId           types.String `tfsdk:"publisher_id"`
	PublisherName         types.String `tfsdk:"publisher_name"`
	LearnMoreUrl          types.String `tfsdk:"learn_more_url"`
	State                 types.String `tfsdk:"state"`
	ApplicationVisibility types.String `tfsdk:"application_visibility"`
}

func (d *EnvironmentApplicationPackagesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *EnvironmentApplicationPackagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of Dynamics 365 applications in a tenant",
		MarkdownDescription: "Fetches the list of Dynamics 365 applications in a tenant.  The data source can be filtered by name and publisher name.\n\nThis is functionally equivalent to the [Environment-level view of apps](https://learn.microsoft.com/en-us/power-platform/admin/manage-apps#environment-level-view-of-apps) in the Power Platform Admin Center or the [`pac application list` command from Power Platform CLI](https://learn.microsoft.com/en-us/power-platform/developer/cli/reference/application#pac-application-list).  This data source uses the [Get Environment Application Package](https://learn.microsoft.com/en-us/rest/api/power-platform/appmanagement/applications/get-environment-application-package) endpoint in the Power Platform API.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Id of the read operation",
				Optional:    true,
			},
			"environment_id": schema.StringAttribute{
				Description: "Id of the Dynamics 365 environment",
				Required:    true,
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
						"application_id": schema.StringAttribute{
							MarkdownDescription: "ApplicaitonId",
							Description:         "ApplicaitonId",
							Computed:            true,
						},
						"application_name": schema.StringAttribute{
							MarkdownDescription: "Name",
							Description:         "Name",
							Computed:            true,
						},
						"unique_name": schema.StringAttribute{
							MarkdownDescription: "Unique Name",
							Description:         "Unique Name",
							Computed:            true,
						},
						"version": schema.StringAttribute{
							MarkdownDescription: "Version",
							Description:         "Version",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Localized Description",
							Description:         "Localized Description",
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
						"learn_more_url": schema.StringAttribute{
							MarkdownDescription: "Learn More Url",
							Description:         "Learn More Url",
							Computed:            true,
						},
						"state": schema.StringAttribute{
							MarkdownDescription: "State",
							Description:         "State",
							Computed:            true,
						},
						"application_visibility": schema.StringAttribute{
							MarkdownDescription: "Application Visibility",
							Description:         "Application Visibility",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *EnvironmentApplicationPackagesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *EnvironmentApplicationPackagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var plan EnvironmentApplicationPackagesListDataSourceModel
	resp.State.Get(ctx, &plan)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE ENVIRONMENT APPLICATION PACKAGES START: %s", d.ProviderTypeName))

	plan.Id = types.StringValue(strconv.FormatInt(time.Now().Unix(), 10))
	plan.EnvironmentId = types.StringValue(plan.EnvironmentId.ValueString())
	plan.Name = types.StringValue(plan.Name.ValueString())
	plan.PublisherName = types.StringValue(plan.PublisherName.ValueString())

	dvExits, err := d.ApplicationClient.DataverseExists(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), err.Error())
	}

	if !dvExits {
		resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", plan.EnvironmentId.ValueString()), "")
		return
	}

	applications, err := d.ApplicationClient.GetApplicationsByEnvironmentId(ctx, plan.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, application := range applications {
		if (plan.Name.ValueString() != "" && plan.Name.ValueString() != application.Name) ||
			(plan.PublisherName.ValueString() != "" && plan.PublisherName.ValueString() != application.PublisherName) {
			continue
		}
		plan.Applications = append(plan.Applications, EnvironmentApplicationPackageDataSourceModel{
			ApplicationId:         types.StringValue(application.ApplicationId),
			Name:                  types.StringValue(application.Name),
			UniqueName:            types.StringValue(application.UniqueName),
			Version:               types.StringValue(application.Version),
			Description:           types.StringValue(application.Description),
			PublisherId:           types.StringValue(application.PublisherId),
			PublisherName:         types.StringValue(application.PublisherName),
			LearnMoreUrl:          types.StringValue(application.LearnMoreUrl),
			State:                 types.StringValue(application.State),
			ApplicationVisibility: types.StringValue(application.ApplicationVisibility),
		})
	}

	diags := resp.State.Set(ctx, &plan)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE ENVIRONMENT APPLICATION PACKAGES END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
