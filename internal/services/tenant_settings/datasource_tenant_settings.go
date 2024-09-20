// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant_settings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

type TenantSettingsDataSource struct {
	helpers.TypeInfo
	TenantSettingsClient TenantSettingsClient
}

type TenantSettingsSourceModel struct {
	Timeouts                                       timeouts.Value `tfsdk:"timeouts"`
	Id                                             types.String   `tfsdk:"id"`
	WalkMeOptOut                                   types.Bool     `tfsdk:"walk_me_opt_out"`
	DisableNPSCommentsReachout                     types.Bool     `tfsdk:"disable_nps_comments_reachout"`
	DisableNewsletterSendout                       types.Bool     `tfsdk:"disable_newsletter_sendout"`
	DisableEnvironmentCreationByNonAdminUsers      types.Bool     `tfsdk:"disable_environment_creation_by_non_admin_users"`
	DisablePortalsCreationByNonAdminUsers          types.Bool     `tfsdk:"disable_portals_creation_by_non_admin_users"`
	DisableSurveyFeedback                          types.Bool     `tfsdk:"disable_survey_feedback"`
	DisableTrialEnvironmentCreationByNonAdminUsers types.Bool     `tfsdk:"disable_trial_environment_creation_by_non_admin_users"`
	DisableCapacityAllocationByEnvironmentAdmins   types.Bool     `tfsdk:"disable_capacity_allocation_by_environment_admins"`
	DisableSupportTicketsVisibleByAllUsers         types.Bool     `tfsdk:"disable_support_tickets_visible_by_all_users"`
	PowerPlatform                                  types.Object   `tfsdk:"power_platform"`
}

type PowerPlatformSettingsModel struct {
	Search                 types.Map `tfsdk:"search"`
	TeamsIntegration       types.Map `tfsdk:"teams_integration"`
	PowerApps              types.Map `tfsdk:"power_apps"`
	PowerAutomate          types.Map `tfsdk:"power_automate"`
	Environments           types.Map `tfsdk:"environments"`
	Governance             types.Map `tfsdk:"governance"`
	Licensing              types.Map `tfsdk:"licensing"`
	PowerPages             types.Map `tfsdk:"power_pages"`
	Champions              types.Map `tfsdk:"champions"`
	Intelligence           types.Map `tfsdk:"intelligence"`
	ModelExperimentation   types.Map `tfsdk:"model_experimentation"`
	CatalogSettings        types.Map `tfsdk:"catalog_settings"`
	UserManagementSettings types.Map `tfsdk:"user_management_settings"`
}

type SearchSettingsModel struct {
	DisableDocsSearch      types.Bool `tfsdk:"disable_docs_search"`
	DisableCommunitySearch types.Bool `tfsdk:"disable_community_search"`
	DisableBingVideoSearch types.Bool `tfsdk:"disable_bing_video_search"`
}

type TeamsIntegrationSettings struct {
	ShareWithColleaguesUserLimit types.Int64 `tfsdk:"share_with_colleagues_user_limit"`
}

type PowerAppsSettings struct {
	DisableShareWithEveryone             types.Bool `tfsdk:"disable_share_with_everyone"`
	EnableGuestsToMake                   types.Bool `tfsdk:"enable_guests_to_make"`
	DisableMakerMatch                    types.Bool `tfsdk:"disable_maker_match"`
	DisableUnusedLicenseAssignment       types.Bool `tfsdk:"disable_unused_license_assignment"`
	DisableCreateFromImage               types.Bool `tfsdk:"disable_create_from_image"`
	DisableCreateFromFigma               types.Bool `tfsdk:"disable_create_from_figma"`
	DisableConnectionSharingWithEveryone types.Bool `tfsdk:"disable_connection_sharing_with_everyone"`
}

type PowerAutomateSettings struct {
	DisableCopilot types.Bool `tfsdk:"disable_copilot"`
}

type EnvironmentsSettings struct {
	DisablePreferredDataLocationForTeamsEnvironment types.Bool `tfsdk:"disable_preferred_data_location_for_teams_environment"`
}

type GovernanceSettings struct {
	DisableAdminDigest                                 types.Bool       `tfsdk:"disable_admin_digest"`
	DisableDeveloperEnvironmentCreationByNonAdminUsers types.Bool       `tfsdk:"disable_developer_environment_creation_by_non_admin_users"`
	EnableDefaultEnvironmentRouting                    types.Bool       `tfsdk:"enable_default_environment_routing"`
	EnvironmentRoutingAllMakers                        types.Bool       `tfsdk:"environment_routing_all_makers"`
	EnvironmentRoutingTargetEnvironmentGroupId         customtypes.UUID `tfsdk:"environment_routing_target_environment_group_id"`
	EnvironmentRoutingTargetSecurityGroupId            customtypes.UUID `tfsdk:"environment_routing_target_security_group_id"`
	Policy                                             types.Object     `tfsdk:"policy"`
}

type PolicySettings struct {
	EnableDesktopFlowDataPolicyManagement types.Bool `tfsdk:"enable_desktop_flow_data_policy_management"`
}

type LicensingSettings struct {
	DisableBillingPolicyCreationByNonAdminUsers     types.Bool  `tfsdk:"disable_billing_policy_creation_by_non_admin_users"`
	EnableTenantCapacityReportForEnvironmentAdmins  types.Bool  `tfsdk:"enable_tenant_capacity_report_for_environment_admins"`
	StorageCapacityConsumptionWarningThreshold      types.Int64 `tfsdk:"storage_capacity_consumption_warning_threshold"`
	EnableTenantLicensingReportForEnvironmentAdmins types.Bool  `tfsdk:"enable_tenant_licensing_report_for_environment_admins"`
	DisableUseOfUnassignedAIBuilderCredits          types.Bool  `tfsdk:"disable_use_of_unassigned_ai_builder_credits"`
}

type PowerPagesSettings struct {
}

type ChampionsSettings struct {
	DisableChampionsInvitationReachout   types.Bool `tfsdk:"disable_champions_invitation_reachout"`
	DisableSkillsMatchInvitationReachout types.Bool `tfsdk:"disable_skills_match_invitation_reachout"`
}

type IntelligenceSettings struct {
	DisableCopilot            types.Bool `tfsdk:"disable_copilot"`
	EnableOpenAiBotPublishing types.Bool `tfsdk:"enable_open_ai_bot_publishing"`
}

type ModelExperimentationSettings struct {
	EnableModelDataSharing types.Bool `tfsdk:"enable_model_data_sharing"`
	DisableDataLogging     types.Bool `tfsdk:"disable_data_logging"`
}

type CatalogSettingsSettings struct {
	PowerCatalogAudienceSetting types.String `tfsdk:"power_catalog_audience_setting"`
}

type UserManagementSettings struct {
	EnableDeleteDisabledUserinAllEnvironments types.Bool `tfsdk:"enable_delete_disabled_user_in_all_environments"`
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
	d.TenantSettingsClient = NewTenantSettingsClient(client)
}

func (d *TenantSettingsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	ctx, exitContext := helpers.EnterRequestContext(ctx, d.TypeInfo, req)
	defer exitContext()

	var state TenantSettingsSourceModel
	resp.Diagnostics.Append(resp.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenantSettings, err := d.TenantSettingsClient.GetTenantSettings(ctx)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.ProviderTypeName), err.Error())
		return
	}

	var configuredSettings TenantSettingsSourceModel
	req.Config.Get(ctx, &configuredSettings)
	state, _ = ConvertFromTenantSettingsDto(*tenantSettings, state.Timeouts)
	hash, err := tenantSettings.CalcObjectHash()
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Error calculating hash for %s", d.ProviderTypeName), err.Error())
	}
	state.Id = types.StringValue(*hash)

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
			"id": schema.StringAttribute{
				Description:         "Id of the read operation",
				MarkdownDescription: "Id of the read operation",
				Computed:            true,
			},
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
