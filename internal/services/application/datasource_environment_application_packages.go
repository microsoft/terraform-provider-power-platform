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
	_ datasource.DataSource              = &EnvironmentApplicationPackagesDataSource{}
	_ datasource.DataSourceWithConfigure = &EnvironmentApplicationPackagesDataSource{}
)

func NewEnvironmentApplicationPackagesDataSource() datasource.DataSource {
	return &EnvironmentApplicationPackagesDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_application_packages",
		},
	}
}

type EnvironmentApplicationPackagesDataSource struct {
	helpers.TypeInfo
	ApplicationClient client
}

type EnvironmentApplicationPackagesListDataSourceModel struct {
	Timeouts      timeouts.Value                                 `tfsdk:"timeouts"`
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

func (d *EnvironmentApplicationPackagesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *EnvironmentApplicationPackagesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of Dynamics 365 applications in a tenant",
		MarkdownDescription: "Fetches the list of Dynamics 365 applications in a tenant.  The data source can be filtered by name and publisher name.\n\nThis is functionally equivalent to the [Environment-level view of apps](https://learn.microsoft.com/power-platform/admin/manage-apps#environment-level-view-of-apps) in the Power Platform Admin Center or the [`pac application list` command from Power Platform CLI](https://learn.microsoft.com/power-platform/developer/cli/reference/application#pac-application-list).  This data source uses the [Get Environment Application Package](https://learn.microsoft.com/rest/api/power-platform/appmanagement/applications/get-environment-application-package) endpoint in the Power Platform API.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
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

func (d *EnvironmentApplicationPackagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	var state EnvironmentApplicationPackagesListDataSourceModel
	resp.State.Get(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE ENVIRONMENT APPLICATION PACKAGES START: %s", d.ProviderTypeName))

	state.EnvironmentId = types.StringValue(state.EnvironmentId.ValueString())
	state.Name = types.StringValue(state.Name.ValueString())
	state.PublisherName = types.StringValue(state.PublisherName.ValueString())

	dvExits, err := d.ApplicationClient.DataverseExists(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when checking if Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), err.Error())
	}

	if !dvExits {
		resp.Diagnostics.AddError(fmt.Sprintf("No Dataverse exists in environment '%s'", state.EnvironmentId.ValueString()), "")
		return
	}

	applications, err := d.ApplicationClient.GetApplicationsByEnvironmentId(ctx, state.EnvironmentId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	for _, application := range applications {
		if (state.Name.ValueString() != "" && state.Name.ValueString() != application.Name) ||
			(state.PublisherName.ValueString() != "" && state.PublisherName.ValueString() != application.PublisherName) {
			continue
		}
		state.Applications = append(state.Applications, EnvironmentApplicationPackageDataSourceModel{
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

	state.Id = types.StringValue(fmt.Sprintf("%s_%d", state.EnvironmentId.ValueString(), len(applications)))
	diags := resp.State.Set(ctx, &state)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE ENVIRONMENT APPLICATION PACKAGES END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
