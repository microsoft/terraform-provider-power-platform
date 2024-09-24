// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant_settings

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/terraform-provider-power-platform/internal/customtypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

type TenantSettingsDataSource struct {
	helpers.TypeInfo
	TenantSettingsClient client
}

type TenantSettingsDataSourceModel struct {
	Timeouts                                       timeouts.Value `tfsdk:"timeouts"`
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

type TenantSettingsResourceModel struct {
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

type TenantSettingsResource struct {
	helpers.TypeInfo
	TenantSettingClient client
}
