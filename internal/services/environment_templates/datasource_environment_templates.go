// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_templates

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
	_ datasource.DataSource              = &EnvironmentTemplatesDataSource{}
	_ datasource.DataSourceWithConfigure = &EnvironmentTemplatesDataSource{}
)

func NewEnvironmentTemplatesDataSource() datasource.DataSource {
	return &EnvironmentTemplatesDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "environment_templates",
		},
	}
}

type EnvironmentTemplatesDataSource struct {
	helpers.TypeInfo
	EnvironmentTemplatesClient EnvironmentTemplatesClient
}

type EnvironmentTemplatesDataSourceModel struct {
	Timeouts  timeouts.Value                  `tfsdk:"timeouts"`
	Id        types.Int64                     `tfsdk:"id"`
	Location  types.String                    `tfsdk:"location"`
	Templates []EnvironmentTemplatesDataModel `tfsdk:"environment_templates"`
}

type EnvironmentTemplatesDataModel struct {
	Category                     string `tfsdk:"category"`
	ID                           string `tfsdk:"id"`
	Name                         string `tfsdk:"name"`
	DisplayName                  string `tfsdk:"display_name"`
	Location                     string `tfsdk:"location"`
	IsDisabled                   bool   `tfsdk:"is_disabled"`
	DisabledReasonCode           string `tfsdk:"disabled_reason_code"`
	DisabledReasonMessage        string `tfsdk:"disabled_reason_message"`
	IsCustomerEngagement         bool   `tfsdk:"is_customer_engagement"`
	IsSupportedForResetOperation bool   `tfsdk:"is_supported_for_reset_operation"`
}

func (d *EnvironmentTemplatesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}

func (d *EnvironmentTemplatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of Dynamics 365 environment templates.",
		MarkdownDescription: "Fetches the list of Dynamics 365 environment templates.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Read: true,
			}),
			"id": schema.Int64Attribute{
				Description: "Id of the read operation",
				Optional:    true,
			},
			"location": schema.StringAttribute{
				Description: "Location of the environment templates",
				Required:    true,
			},
			"environment_templates": schema.ListNestedAttribute{
				Description:         "List of available environment templates",
				MarkdownDescription: "List of available environment templates",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"category": schema.StringAttribute{
							Description: "Category of the environment template",
							Computed:    true,
						},
						"id": schema.StringAttribute{
							Description: "Unique identifier of the environment template",
							Computed:    true,
						},
						"name": schema.StringAttribute{
							Description: "Name of the environment template",
							Computed:    true,
						},
						"display_name": schema.StringAttribute{
							Description: "Display name of the environment template",
							Computed:    true,
						},
						"location": schema.StringAttribute{
							Description: "Location of the environment template",
							Computed:    true,
						},
						"is_disabled": schema.BoolAttribute{
							Description: "Indicates if the environment template is disabled",
							Computed:    true,
						},
						"disabled_reason_code": schema.StringAttribute{
							Description: "Code of the reason why the environment template is disabled",
							Computed:    true,
						},
						"disabled_reason_message": schema.StringAttribute{
							Description: "Message of the reason why the environment template is disabled",
							Computed:    true,
						},
						"is_customer_engagement": schema.BoolAttribute{
							Description: "Indicates if the environment template is for customer engagement",
							Computed:    true,
						},
						"is_supported_for_reset_operation": schema.BoolAttribute{
							Description: "Indicates if the environment template is supported for reset operation",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (d *EnvironmentTemplatesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
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
	d.EnvironmentTemplatesClient = NewEnvironmentTemplatesClient(clientApi)
}

func (d *EnvironmentTemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	appendToList := func(items []ItemDto, category string, list *[]EnvironmentTemplatesDataModel) {
		for _, item := range items {
			*list = append(*list, EnvironmentTemplatesDataModel{
				Category:                     category,
				ID:                           item.ID,
				Name:                         item.Name,
				DisplayName:                  item.Properties.DisplayName,
				Location:                     item.Location,
				IsDisabled:                   item.Properties.IsDisabled,
				DisabledReasonCode:           item.Properties.DisabledReason.Code,
				DisabledReasonMessage:        item.Properties.DisabledReason.Message,
				IsCustomerEngagement:         item.Properties.IsCustomerEngagement,
				IsSupportedForResetOperation: item.Properties.IsSupportedForResetOperation,
			})
		}
	}

	var state EnvironmentTemplatesDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	environment_templates, err := d.EnvironmentTemplatesClient.GetEnvironmentTemplatesByLocation(ctx, state.Location.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	state.Templates = make([]EnvironmentTemplatesDataModel, 0)
	appendToList(environment_templates.Standard, "standard", &state.Templates)
	appendToList(environment_templates.Premium, "premium", &state.Templates)
	appendToList(environment_templates.Developer, "developer", &state.Templates)
	appendToList(environment_templates.Basic, "basic", &state.Templates)
	appendToList(environment_templates.Production, "production", &state.Templates)
	appendToList(environment_templates.Sandbox, "sandbox", &state.Templates)
	appendToList(environment_templates.Trial, "trial", &state.Templates)
	appendToList(environment_templates.Default, "default", &state.Templates)
	appendToList(environment_templates.Support, "support", &state.Templates)
	appendToList(environment_templates.SubscriptionBasedTrial, "subscriptionBasedTrial", &state.Templates)
	appendToList(environment_templates.Teams, "teams", &state.Templates)
	appendToList(environment_templates.Platform, "platform", &state.Templates)

	state.Id = types.Int64Value(int64(len(state.Templates)))
	state.Location = types.StringValue(state.Location.ValueString())

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
