// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant_settings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/api"
	"github.com/microsoft/terraform-provider-power-platform/internal/customtypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

var (
	_ datasource.DataSource              = &TenantSettingsDataSource{}
	_ datasource.DataSourceWithConfigure = &TenantSettingsDataSource{}
)

func NewTenantSettingsDataSource() datasource.DataSource {
	return &TenantSettingsDataSource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "tenant_settings",
		},
	}
}

func (d *TenantSettingsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	if req.ProviderData == nil {
		// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
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
	d.TenantSettingsClient = newTenantSettingsClient(client)
}

func (d *TenantSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state TenantSettingsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantSettings, err := d.TenantSettingsClient.GetTenantSettings(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	var configuredSettings TenantSettingsDataSourceModel
	req.Config.Get(ctx, &configuredSettings)
	state, _ = convertFromTenantSettingsDto[TenantSettingsDataSourceModel](*tenantSettings, state.Timeouts)
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (d *TenantSettingsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()
	resp.Schema = schema.Schema{
		Description:         "Power Platform Tenant Settings Data Source",
		MarkdownDescription: "Fetches Power Platform Tenant Settings.  See [Tenant Settings Overview](https://learn.microsoft.com/power-platform/admin/tenant-settings) for more information.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: false,
				Update: false,
				Delete: false,
				Read:   false,
			}),
			"walk_me_opt_out": schema.BoolAttribute{
				Description: "Walk Me Opt Out",
				Computed:    true,
			},
			"disable_nps_comments_reachout": schema.BoolAttribute{
				Description: "Disable NPS Comments Reachout",
				Computed:    true,
			},
			"disable_newsletter_sendout": schema.BoolAttribute{
				Description: "Disable Newsletter Sendout",
				Computed:    true,
			},
			"disable_environment_creation_by_non_admin_users": schema.BoolAttribute{
				Description: "Disable Environment Creation By Non Admin Users",
				Computed:    true,
			},
			"disable_portals_creation_by_non_admin_users": schema.BoolAttribute{
				Description: "Disable Portals Creation By Non Admin Users",
				Computed:    true,
			},
			"disable_survey_feedback": schema.BoolAttribute{
				Description: "Disable Survey Feedback",
				Computed:    true,
			},
			"disable_trial_environment_creation_by_non_admin_users": schema.BoolAttribute{
				Description: "Disable Trial Environment Creation By Non Admin Users",
				Computed:    true,
			},
			"disable_capacity_allocation_by_environment_admins": schema.BoolAttribute{
				Description: "Disable Capacity Allocation By Environment Admins",
				Computed:    true,
			},
			"disable_support_tickets_visible_by_all_users": schema.BoolAttribute{
				Description: "Disable Support Tickets Visible By All Users",
				Computed:    true,
			},
			"power_platform": schema.SingleNestedAttribute{
				Description: "Power Platform",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"search": schema.SingleNestedAttribute{
						Description: "Search",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"disable_docs_search": schema.BoolAttribute{
								Description: "Disable Docs Search",
								Computed:    true,
							},
							"disable_community_search": schema.BoolAttribute{
								Description: "Disable Community Search",
								Computed:    true,
							},
							"disable_bing_video_search": schema.BoolAttribute{
								Description: "Disable Bing Video Search",
								Computed:    true,
							},
						},
					},
					"teams_integration": schema.SingleNestedAttribute{
						Description: "Teams Integration",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"share_with_colleagues_user_limit": schema.Int64Attribute{
								Description: "Share With Colleagues User Limit",
								Computed:    true,
							},
						},
					},
					"power_apps": schema.SingleNestedAttribute{
						Description: "Power Apps",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"disable_share_with_everyone": schema.BoolAttribute{
								Description: "Disable Share With Everyone",
								Computed:    true,
							},
							"enable_guests_to_make": schema.BoolAttribute{
								Description: "Enable Guests To Make",
								Computed:    true,
							},
							"disable_maker_match": schema.BoolAttribute{
								Description: "Disable Maker Match",
								Computed:    true,
							},
							"disable_unused_license_assignment": schema.BoolAttribute{
								Description: "Disable Unused License Assignment",
								Computed:    true,
							},
							"disable_create_from_image": schema.BoolAttribute{
								Description: "Disable Create From Image",
								Computed:    true,
							},
							"disable_create_from_figma": schema.BoolAttribute{
								Description: "Disable Create From Figma",
								Computed:    true,
							},
							"disable_connection_sharing_with_everyone": schema.BoolAttribute{
								Description: "Disable Connection Sharing With Everyone",
								Computed:    true,
							},
						},
					},
					"power_automate": schema.SingleNestedAttribute{
						Description: "Power Automate",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"disable_copilot": schema.BoolAttribute{
								Description: "Disable Copilot",
								Computed:    true,
							},
						},
					},
					"environments": schema.SingleNestedAttribute{
						Description: "Environments",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"disable_preferred_data_location_for_teams_environment": schema.BoolAttribute{
								Description: "Disable Preferred Data Location For Teams Environment",
								Computed:    true,
							},
						},
					},
					"governance": schema.SingleNestedAttribute{
						Description: "Governance",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"disable_admin_digest": schema.BoolAttribute{
								Description: "Disable Admin Digest",
								Computed:    true,
							},
							"disable_developer_environment_creation_by_non_admin_users": schema.BoolAttribute{
								Description: "Disable Developer Environment Creation By Non Admin Users",
								Computed:    true,
							},
							"enable_default_environment_routing": schema.BoolAttribute{
								Description: "Enable Default Environment Routing",
								Computed:    true,
							},
							"environment_routing_all_makers": schema.BoolAttribute{
								Description: "Select who can be routed to a new personal developer environment. (All Makers = true, New Makers = false)",
								Computed:    true,
							},
							"environment_routing_target_environment_group_id": schema.StringAttribute{
								Description: "Assign newly created personal developer environments to a specific environment group",
								Computed:    true,
								CustomType:  customtypes.UUIDType{},
							},
							"environment_routing_target_security_group_id": schema.StringAttribute{
								Description: "Restrict routing to members of the following security group. (00000000-0000-0000-0000-000000000000 allows all users)",
								Computed:    true,
								CustomType:  customtypes.UUIDType{},
							},
							"policy": schema.SingleNestedAttribute{
								Description: "Policy",
								Computed:    true,
								Attributes: map[string]schema.Attribute{
									"enable_desktop_flow_data_policy_management": schema.BoolAttribute{
										Description: "Enable Desktop Flow Data Policy Management",
										Computed:    true,
									},
								},
							},
						},
					},
					"licensing": schema.SingleNestedAttribute{
						Description: "Licensing",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"disable_billing_policy_creation_by_non_admin_users": schema.BoolAttribute{
								Description: "Disable Billing Policy Creation By Non Admin Users",
								Computed:    true,
							},
							"enable_tenant_capacity_report_for_environment_admins": schema.BoolAttribute{
								Description: "Enable Tenant Capacity Report For Environment Admins",
								Computed:    true,
							},
							"storage_capacity_consumption_warning_threshold": schema.Int64Attribute{
								Description: "Storage Capacity Consumption Warning Threshold",
								Computed:    true,
							},
							"enable_tenant_licensing_report_for_environment_admins": schema.BoolAttribute{
								Description: "Enable Tenant Licensing Report For Environment Admins",
								Computed:    true,
							},
							"disable_use_of_unassigned_ai_builder_credits": schema.BoolAttribute{
								Description: "Disable Use Of Unassigned AI Builder Credits",
								Computed:    true,
							},
						},
					},
					"power_pages": schema.SingleNestedAttribute{
						Description: "Power Pages",
						Computed:    true,
						Attributes:  map[string]schema.Attribute{},
					},
					"champions": schema.SingleNestedAttribute{
						Description: "Champions",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"disable_champions_invitation_reachout": schema.BoolAttribute{
								Description: "Disable Champions Invitation Reachout",
								Computed:    true,
							},
							"disable_skills_match_invitation_reachout": schema.BoolAttribute{
								Description: "Disable Skills Match Invitation Reachout",
								Computed:    true,
							},
						},
					},
					"intelligence": schema.SingleNestedAttribute{
						Description: "Intelligence",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"disable_copilot": schema.BoolAttribute{
								Description: "Disable Copilot",
								Computed:    true,
							},
							"enable_open_ai_bot_publishing": schema.BoolAttribute{
								Description: "Enable Open AI Bot Publishing",
								Computed:    true,
							},
						},
					},
					"model_experimentation": schema.SingleNestedAttribute{
						Description: "Model Experimentation",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"enable_model_data_sharing": schema.BoolAttribute{
								Description: "Enable Model Data Sharing",
								Computed:    true,
							},
							"disable_data_logging": schema.BoolAttribute{
								Description: "Disable Data Logging",
								Computed:    true,
							},
						},
					},
					"catalog_settings": schema.SingleNestedAttribute{
						Description: "Catalog Settings",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"power_catalog_audience_setting": schema.StringAttribute{
								Description: "Power Catalog Audience Setting",
								Computed:    true,
							},
						},
					},
					"user_management_settings": schema.SingleNestedAttribute{
						Description: "User Management Settings",
						Computed:    true,
						Attributes: map[string]schema.Attribute{
							"enable_delete_disabled_user_in_all_environments": schema.BoolAttribute{
								Description: "Enable Delete Disabled User In All Environments",
								Computed:    true,
							},
						},
					},
				},
			},
		},
	}
}

// Metadata returns the metadata for the resource, which includes the resource type name.
func (d *TenantSettingsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	// update our own internal storage of the provider type name.
	d.ProviderTypeName = req.ProviderTypeName

	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	// Set the type name for the resource to providername_resourcename.
	resp.TypeName = d.FullTypeName()
	tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}
