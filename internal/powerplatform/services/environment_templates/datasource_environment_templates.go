// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var (
	_ datasource.DataSource              = &EnvironmentTemplatesDataSource{}
	_ datasource.DataSourceWithConfigure = &EnvironmentTemplatesDataSource{}
)

func NewEnvironmentTemplatesDataSource() datasource.DataSource {
	return &EnvironmentTemplatesDataSource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_environment_templates",
	}
}

type EnvironmentTemplatesDataSource struct {
	EnvironmentTemplatesClient EnvironmentTemplatesClient
	ProviderTypeName           string
	TypeName                   string
}

type EnvironmentTemplatesDataSourceModel struct {
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

func (d *EnvironmentTemplatesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + d.TypeName
}

func (d *EnvironmentTemplatesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Fetches the list of Dynamics 365 environment templates.",
		MarkdownDescription: "Fetches the list of Dynamics 365 environment templates.",
		Attributes: map[string]schema.Attribute{
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
	d.EnvironmentTemplatesClient = NewEnvironmentTemplatesClient(clientApi)
}

func (d *EnvironmentTemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	appendToList := func(items []EnvironmentTemplateItem, category string, list *[]EnvironmentTemplatesDataModel) {
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

	var plan EnvironmentTemplatesDataSourceModel
	resp.State.Get(ctx, &plan)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE ENVIRONMENT TEMPLATES START: %s", d.ProviderTypeName))

	environment_templates, err := d.EnvironmentTemplatesClient.GetEnvironmentTemplatesByLocation(ctx, plan.Location.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	plan.Templates = make([]EnvironmentTemplatesDataModel, 0)
	appendToList(environment_templates.Standard, "standard", &plan.Templates)
	appendToList(environment_templates.Premium, "premium", &plan.Templates)
	appendToList(environment_templates.Developer, "developer", &plan.Templates)
	appendToList(environment_templates.Basic, "basic", &plan.Templates)
	appendToList(environment_templates.Production, "production", &plan.Templates)
	appendToList(environment_templates.Sandbox, "sandbox", &plan.Templates)
	appendToList(environment_templates.Trial, "trial", &plan.Templates)
	appendToList(environment_templates.Default, "default", &plan.Templates)
	appendToList(environment_templates.Support, "support", &plan.Templates)
	appendToList(environment_templates.SubscriptionBasedTrial, "subscriptionBasedTrial", &plan.Templates)
	appendToList(environment_templates.Teams, "teams", &plan.Templates)
	appendToList(environment_templates.Platform, "platform", &plan.Templates)

	plan.Id = types.Int64Value(int64(len(plan.Templates)))
	plan.Location = types.StringValue(plan.Location.ValueString())

	diags := resp.State.Set(ctx, &plan)

	tflog.Debug(ctx, fmt.Sprintf("READ DATASOURCE ENVIRONMENT TEMPLATES END: %s", d.ProviderTypeName))

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
