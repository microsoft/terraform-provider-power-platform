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

	client, ok := req.ProviderData.(*api.ProviderClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected ProviderData Type",
			fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	d.TenantSettingsClient = newTenantSettingsClient(client.Api)
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
		MarkdownDescription: "Fetches Power Platform Tenant Settings.  See [Tenant Settings Overview](https://learn.microsoft.com/power-platform/admin/tenant-settings) for more information.",
		Attributes: map[string]schema.Attribute{
			"timeouts": timeouts.Attributes(ctx, timeouts.Opts{
				Create: false,
				Update: false,
				Delete: false,
				Read:   false,
			}),
			"walk_me_opt_out": schema.BoolAttribute{
				MarkdownDescription: "Walk Me Opt Out",
				Computed:            true,
			},
			"disable_nps_comments_reachout": schema.BoolAttribute{
				MarkdownDescription: "Disable NPS Comments Reachout",
				Computed:            true,
			},
			"disable_newsletter_sendout": schema.BoolAttribute{
				MarkdownDescription: "Disable Newsletter Sendout",
				Computed:            true,
			},
			"disable_environment_creation_by_non_admin_users": schema.BoolAttribute{
				MarkdownDescription: "Disable Environment Creation By Non Admin Users",
				Computed:            true,
			},
			"disable_portals_creation_by_non_admin_users": schema.BoolAttribute{
				MarkdownDescription: "Disable Portals Creation By Non Admin Users",
				Computed:            true,
			},
			"disable_survey_feedback": schema.BoolAttribute{
				MarkdownDescription: "Disable Survey Feedback",
				Computed:            true,
			},
			"disable_trial_environment_creation_by_non_admin_users": schema.BoolAttribute{
				MarkdownDescription: "Disable Trial Environment Creation By Non Admin Users",
				Computed:            true,
			},
			"disable_capacity_allocation_by_environment_admins": schema.BoolAttribute{
				MarkdownDescription: "Disable Capacity Allocation By Environment Admins",
				Computed:            true,
			},
			"disable_support_tickets_visible_by_all_users": schema.BoolAttribute{
				MarkdownDescription: "Disable Support Tickets Visible By All Users",
				Computed:            true,
			},
			"power_platform": schema.SingleNestedAttribute{
				MarkdownDescription: "Power Platform",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"search": schema.SingleNestedAttribute{
						MarkdownDescription: "Search",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"disable_docs_search": schema.BoolAttribute{
								MarkdownDescription: "Disable Docs Search",
								Computed:            true,
							},
							"disable_community_search": schema.BoolAttribute{
								MarkdownDescription: "Disable Community Search",
								Computed:            true,
							},
							"disable_bing_video_search": schema.BoolAttribute{
								MarkdownDescription: "Disable Bing Video Search",
								Computed:            true,
							},
						},
					},
					"teams_integration": schema.SingleNestedAttribute{
						MarkdownDescription: "Teams Integration",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"share_with_colleagues_user_limit": schema.Int64Attribute{
								MarkdownDescription: "Share With Colleagues User Limit",
								Computed:            true,
							},
						},
					},
					"power_apps": schema.SingleNestedAttribute{
						MarkdownDescription: "Power Apps",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"disable_share_with_everyone": schema.BoolAttribute{
								MarkdownDescription: "Disable Share With Everyone",
								Computed:            true,
							},
							"enable_guests_to_make": schema.BoolAttribute{
								MarkdownDescription: "Enable Guests To Make",
								Computed:            true,
							},
							"disable_maker_match": schema.BoolAttribute{
								MarkdownDescription: "Disable Maker Match",
								Computed:            true,
							},
							"disable_unused_license_assignment": schema.BoolAttribute{
								MarkdownDescription: "Disable Unused License Assignment",
								Computed:            true,
							},
							"disable_create_from_image": schema.BoolAttribute{
								MarkdownDescription: "Disable Create From Image",
								Computed:            true,
							},
							"disable_create_from_figma": schema.BoolAttribute{
								MarkdownDescription: "Disable Create From Figma",
								Computed:            true,
							},
							"disable_connection_sharing_with_everyone": schema.BoolAttribute{
								MarkdownDescription: "Disable Connection Sharing With Everyone",
								Computed:            true,
							},
						},
					},
					"power_automate": schema.SingleNestedAttribute{
						MarkdownDescription: "Power Automate",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"disable_copilot": schema.BoolAttribute{
								MarkdownDescription: "Disable Copilot",
								Computed:            true,
							},
						},
					},
					"environments": schema.SingleNestedAttribute{
						MarkdownDescription: "Environments",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"disable_preferred_data_location_for_teams_environment": schema.BoolAttribute{
								MarkdownDescription: "Disable Preferred Data Location For Teams Environment",
								Computed:            true,
							},
						},
					},
					"governance": schema.SingleNestedAttribute{
						MarkdownDescription: "Governance",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"disable_admin_digest": schema.BoolAttribute{
								MarkdownDescription: "Disable Admin Digest",
								Computed:            true,
							},
							"disable_developer_environment_creation_by_non_admin_users": schema.BoolAttribute{
								MarkdownDescription: "Disable Developer Environment Creation By Non Admin Users",
								Computed:            true,
							},
							"enable_default_environment_routing": schema.BoolAttribute{
								MarkdownDescription: "Enable Default Environment Routing",
								Computed:            true,
							},
							"environment_routing_all_makers": schema.BoolAttribute{
								MarkdownDescription: "Select who can be routed to a new personal developer environment. (All Makers = true, New Makers = false)",
								Computed:            true,
							},
							"environment_routing_target_environment_group_id": schema.StringAttribute{
								MarkdownDescription: "Assign newly created personal developer environments to a specific environment group",
								Computed:            true,
								CustomType:          customtypes.UUIDType{},
							},
							"environment_routing_target_security_group_id": schema.StringAttribute{
								MarkdownDescription: "Restrict routing to members of the following security group. (00000000-0000-0000-0000-000000000000 allows all users)",
								Computed:            true,
								CustomType:          customtypes.UUIDType{},
							},
							"policy": schema.SingleNestedAttribute{
								MarkdownDescription: "Policy",
								Computed:            true,
								Attributes: map[string]schema.Attribute{
									"enable_desktop_flow_data_policy_management": schema.BoolAttribute{
										MarkdownDescription: "Enable Desktop Flow Data Policy Management",
										Computed:            true,
									},
								},
							},
						},
					},
					"licensing": schema.SingleNestedAttribute{
						MarkdownDescription: "Licensing",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"disable_billing_policy_creation_by_non_admin_users": schema.BoolAttribute{
								MarkdownDescription: "Disable Billing Policy Creation By Non Admin Users",
								Computed:            true,
							},
							"enable_tenant_capacity_report_for_environment_admins": schema.BoolAttribute{
								MarkdownDescription: "Enable Tenant Capacity Report For Environment Admins",
								Computed:            true,
							},
							"storage_capacity_consumption_warning_threshold": schema.Int64Attribute{
								MarkdownDescription: "Storage Capacity Consumption Warning Threshold",
								Computed:            true,
							},
							"enable_tenant_licensing_report_for_environment_admins": schema.BoolAttribute{
								MarkdownDescription: "Enable Tenant Licensing Report For Environment Admins",
								Computed:            true,
							},
							"disable_use_of_unassigned_ai_builder_credits": schema.BoolAttribute{
								MarkdownDescription: "Disable Use Of Unassigned AI Builder Credits",
								Computed:            true,
							},
						},
					},
					"power_pages": schema.SingleNestedAttribute{
						MarkdownDescription: "Power Pages",
						Computed:            true,
						Attributes:          map[string]schema.Attribute{},
					},
					"champions": schema.SingleNestedAttribute{
						MarkdownDescription: "Champions",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"disable_champions_invitation_reachout": schema.BoolAttribute{
								MarkdownDescription: "Disable Champions Invitation Reachout",
								Computed:            true,
							},
							"disable_skills_match_invitation_reachout": schema.BoolAttribute{
								MarkdownDescription: "Disable Skills Match Invitation Reachout",
								Computed:            true,
							},
						},
					},
					"intelligence": schema.SingleNestedAttribute{
						MarkdownDescription: "Intelligence",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"disable_copilot": schema.BoolAttribute{
								MarkdownDescription: "Disable Copilot",
								Computed:            true,
							},
							"enable_open_ai_bot_publishing": schema.BoolAttribute{
								MarkdownDescription: "Enable Open AI Bot Publishing",
								Computed:            true,
							},
						},
					},
					"model_experimentation": schema.SingleNestedAttribute{
						MarkdownDescription: "Model Experimentation",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"enable_model_data_sharing": schema.BoolAttribute{
								MarkdownDescription: "Enable Model Data Sharing",
								Computed:            true,
							},
							"disable_data_logging": schema.BoolAttribute{
								MarkdownDescription: "Disable Data Logging",
								Computed:            true,
							},
						},
					},
					"catalog_settings": schema.SingleNestedAttribute{
						MarkdownDescription: "Catalog Settings",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"power_catalog_audience_setting": schema.StringAttribute{
								MarkdownDescription: "Power Catalog Audience Setting",
								Computed:            true,
							},
						},
					},
					"user_management_settings": schema.SingleNestedAttribute{
						MarkdownDescription: "User Management Settings",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"enable_delete_disabled_user_in_all_environments": schema.BoolAttribute{
								MarkdownDescription: "Enable Delete Disabled User In All Environments",
								Computed:            true,
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
