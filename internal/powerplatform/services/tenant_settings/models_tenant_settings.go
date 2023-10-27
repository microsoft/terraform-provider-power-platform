package powerplatform

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type PowerPlatformSettingsDto struct {
	Search                 *SearchSettingsDto               `json:"search,omitempty"`
	TeamsIntegration       *TeamsIntegrationSettingsDto     `json:"teamsIntegration,omitempty"`
	PowerApps              *PowerAppsSettingsDto            `json:"powerApps,omitempty"`
	PowerAutomate          *PowerAutomateSettingsDto        `json:"powerAutomate,omitempty"`
	Environments           *EnvironmentSettingsDto          `json:"environments,omitempty"`
	Governance             *GovernanceSettingsDto           `json:"governance,omitempty"`
	Licensing              *LicenseSettingsDto              `json:"licensing,omitempty"`
	PowerPages             *PowerPagesSettingsDto           `json:"powerPages,omitempty"`
	Champions              *ChampionSettingsDto             `json:"champions,omitempty"`
	Intelligence           *IntelligenceSettingsDto         `json:"intelligence,omitempty"`
	ModelExperimentation   *ModelExperimentationSettingsDto `json:"modelExperimentation,omitempty"`
	CatalogSettings        *CatalogSettingsDto              `json:"catalogSettings,omitempty"`
	UserManagementSettings *UserManagementSettingsDto       `json:"userManagementSettings,omitempty"`
}

type TeamsIntegrationSettingsDto struct {
	ShareWithColleaguesUserLimit *int64 `json:"shareWithColleaguesUserLimit,omitempty"`
}

type PowerAutomateSettingsDto struct {
	DisableCopilot *bool `json:"disableCopilot,omitempty"`
}

type PowerAppsSettingsDto struct {
	DisableShareWithEveryone             *bool `json:"disableShareWithEveryone,omitempty"`
	EnableGuestsToMake                   *bool `json:"enableGuestsToMake,omitempty"`
	DisableMembersIndicator              *bool `json:"disableMembersIndicator,omitempty"`
	DisableMakerMatch                    *bool `json:"disableMakerMatch,omitempty"`
	DisableUnusedLicenseAssignment       *bool `json:"disableUnusedLicenseAssignment,omitempty"`
	DisableCreateFromImage               *bool `json:"disableCreateFromImage,omitempty"`
	DisableCreateFromFigma               *bool `json:"disableCreateFromFigma,omitempty"`
	DisableConnectionSharingWithEveryone *bool `json:"disableConnectionSharingWithEveryone,omitempty"`
}

type EnvironmentSettingsDto struct {
	DisablePreferredDataLocationForTeamsEnvironment *bool `json:"disablePreferredDataLocationForTeamsEnvironment,omitempty"`
}

type SearchSettingsDto struct {
	DisableDocsSearch      *bool `json:"disableDocsSearch,omitempty"`
	DisableCommunitySearch *bool `json:"disableCommunitySearch,omitempty"`
	DisableBingVideoSearch *bool `json:"disableBingVideoSearch,omitempty"`
}

type GovernanceSettingsDto struct {
	DisableAdminDigest                                 *bool              `json:"disableAdminDigest,omitempty"`
	DisableDeveloperEnvironmentCreationByNonAdminUsers *bool              `json:"disableDeveloperEnvironmentCreationByNonAdminUsers,omitempty"`
	EnableDefaultEnvironmentRouting                    *bool              `json:"enableDefaultEnvironmentRouting,omitempty"`
	Policy                                             *PolicySettingsDto `json:"policy,omitempty"`
}

type PolicySettingsDto struct {
	EnableDesktopFlowDataPolicyManagement *bool `json:"enableDesktopFlowDataPolicyManagement,omitempty"`
}

type LicenseSettingsDto struct {
	DisableBillingPolicyCreationByNonAdminUsers     *bool  `json:"disableBillingPolicyCreationByNonAdminUsers,omitempty"`
	EnableTenantCapacityReportForEnvironmentAdmins  *bool  `json:"enableTenantCapacityReportForEnvironmentAdmins,omitempty"`
	StorageCapacityConsumptionWarningThreshold      *int64 `json:"storageCapacityConsumptionWarningThreshold,omitempty"`
	EnableTenantLicensingReportForEnvironmentAdmins *bool  `json:"enableTenantLicensingReportForEnvironmentAdmins,omitempty"`
	DisableUseOfUnassignedAIBuilderCredits          *bool  `json:"disableUseOfUnassignedAIBuilderCredits,omitempty"`
}

type PowerPagesSettingsDto struct {
}

type ChampionSettingsDto struct {
	DisableChampionsInvitationReachout   *bool `json:"disableChampionsInvitationReachout,omitempty"`
	DisableSkillsMatchInvitationReachout *bool `json:"disableSkillsMatchInvitationReachout,omitempty"`
}

type IntelligenceSettingsDto struct {
	DisableCopilot            *bool `json:"disableCopilot,omitempty"`
	EnableOpenAiBotPublishing *bool `json:"enableOpenAiBotPublishing,omitempty"`
}

type ModelExperimentationSettingsDto struct {
	EnableModelDataSharing *bool `json:"enableModelDataSharing,omitempty"`
	DisableDataLogging     *bool `json:"disableDataLogging,omitempty"`
}

type CatalogSettingsDto struct {
	PowerCatalogAudienceSetting *string `json:"powerCatalogAudienceSetting,omitempty"`
}

type UserManagementSettingsDto struct {
	EnableDeleteDisabledUserinAllEnvironments *bool `json:"enableDeleteDisabledUserinAllEnvironments,omitempty"`
}

type TenantSettingsDto struct {
	WalkMeOptOut                                   *bool                     `json:"walkMeOptOut,omitempty"`
	DisableNPSCommentsReachout                     *bool                     `json:"disableNPSCommentsReachout,omitempty"`
	DisableNewsletterSendout                       *bool                     `json:"disableNewsletterSendout,omitempty"`
	DisableEnvironmentCreationByNonAdminUsers      *bool                     `json:"disableEnvironmentCreationByNonAdminUsers,omitempty"`
	DisablePortalsCreationByNonAdminUsers          *bool                     `json:"disablePortalsCreationByNonAdminUsers,omitempty"`
	DisableSurveyFeedback                          *bool                     `json:"disableSurveyFeedback,omitempty"`
	DisableTrialEnvironmentCreationByNonAdminUsers *bool                     `json:"disableTrialEnvironmentCreationByNonAdminUsers,omitempty"`
	DisableCapacityAllocationByEnvironmentAdmins   *bool                     `json:"disableCapacityAllocationByEnvironmentAdmins,omitempty"`
	DisableSupportTicketsVisibleByAllUsers         *bool                     `json:"disableSupportTicketsVisibleByAllUsers,omitempty"`
	PowerPlatform                                  *PowerPlatformSettingsDto `json:"powerPlatform,omitempty"`
}

func ConvertFromTenantSettingsModel(ctx context.Context, tenantSettings TenantSettingsSourceModel) TenantSettingsDto {
	tenantSettingsDto := TenantSettingsDto{}

	if !tenantSettings.WalkMeOptOut.IsNull() && !tenantSettings.WalkMeOptOut.IsUnknown() {
		tenantSettingsDto.WalkMeOptOut = tenantSettings.WalkMeOptOut.ValueBoolPointer()
	}
	if !tenantSettings.DisableNPSCommentsReachout.IsNull() && !tenantSettings.DisableNPSCommentsReachout.IsUnknown() {
		tenantSettingsDto.DisableNPSCommentsReachout = tenantSettings.DisableNPSCommentsReachout.ValueBoolPointer()
	}
	if !tenantSettings.DisableNewsletterSendout.IsNull() && !tenantSettings.DisableNewsletterSendout.IsUnknown() {
		tenantSettingsDto.DisableNewsletterSendout = tenantSettings.DisableNewsletterSendout.ValueBoolPointer()
	}
	if !tenantSettings.DisableEnvironmentCreationByNonAdminUsers.IsNull() && !tenantSettings.DisableEnvironmentCreationByNonAdminUsers.IsUnknown() {
		tenantSettingsDto.DisableEnvironmentCreationByNonAdminUsers = tenantSettings.DisableEnvironmentCreationByNonAdminUsers.ValueBoolPointer()
	}
	if !tenantSettings.DisablePortalsCreationByNonAdminUsers.IsNull() && !tenantSettings.DisablePortalsCreationByNonAdminUsers.IsUnknown() {
		tenantSettingsDto.DisablePortalsCreationByNonAdminUsers = tenantSettings.DisablePortalsCreationByNonAdminUsers.ValueBoolPointer()
	}
	if !tenantSettings.DisableSurveyFeedback.IsNull() && !tenantSettings.DisableSurveyFeedback.IsUnknown() {
		tenantSettingsDto.DisableSurveyFeedback = tenantSettings.DisableSurveyFeedback.ValueBoolPointer()
	}
	if !tenantSettings.DisableTrialEnvironmentCreationByNonAdminUsers.IsNull() && !tenantSettings.DisableTrialEnvironmentCreationByNonAdminUsers.IsUnknown() {
		tenantSettingsDto.DisableTrialEnvironmentCreationByNonAdminUsers = tenantSettings.DisableTrialEnvironmentCreationByNonAdminUsers.ValueBoolPointer()
	}
	if !tenantSettings.DisableCapacityAllocationByEnvironmentAdmins.IsNull() && !tenantSettings.DisableCapacityAllocationByEnvironmentAdmins.IsUnknown() {
		tenantSettingsDto.DisableCapacityAllocationByEnvironmentAdmins = tenantSettings.DisableCapacityAllocationByEnvironmentAdmins.ValueBoolPointer()
	}
	if !tenantSettings.DisableSupportTicketsVisibleByAllUsers.IsNull() && !tenantSettings.DisableSupportTicketsVisibleByAllUsers.IsUnknown() {
		tenantSettingsDto.DisableSupportTicketsVisibleByAllUsers = tenantSettings.DisableSupportTicketsVisibleByAllUsers.ValueBoolPointer()
	}

	powerPlatformAttributes := tenantSettings.PowerPlatform.Attributes()
	searchObject := powerPlatformAttributes["search"]
	if !searchObject.IsNull() && !searchObject.IsUnknown() {
		var searchSettings SearchSettingsModel
		searchObject.(basetypes.ObjectValue).As(ctx, &searchSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &PowerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.Search = &SearchSettingsDto{}
		if !searchSettings.DisableDocsSearch.IsNull() && !searchSettings.DisableDocsSearch.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Search.DisableDocsSearch = searchSettings.DisableDocsSearch.ValueBoolPointer()
		}
		if !searchSettings.DisableCommunitySearch.IsNull() && !searchSettings.DisableCommunitySearch.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Search.DisableCommunitySearch = searchSettings.DisableCommunitySearch.ValueBoolPointer()
		}
		if !searchSettings.DisableBingVideoSearch.IsNull() && !searchSettings.DisableBingVideoSearch.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Search.DisableBingVideoSearch = searchSettings.DisableBingVideoSearch.ValueBoolPointer()
		}
	}
	teamIntegrationObject := powerPlatformAttributes["teams_integration"]
	if !teamIntegrationObject.IsNull() && !teamIntegrationObject.IsUnknown() {
		var teamsIntegrationSettings TeamsIntegrationSettings
		teamIntegrationObject.(basetypes.ObjectValue).As(ctx, &teamsIntegrationSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &PowerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.TeamsIntegration = &TeamsIntegrationSettingsDto{}
		if !teamsIntegrationSettings.ShareWithColleaguesUserLimit.IsNull() && !teamsIntegrationSettings.ShareWithColleaguesUserLimit.IsUnknown() {
			tenantSettingsDto.PowerPlatform.TeamsIntegration.ShareWithColleaguesUserLimit = teamsIntegrationSettings.ShareWithColleaguesUserLimit.ValueInt64Pointer()
		}
	}
	powerAppsObject := powerPlatformAttributes["power_apps"]
	if !powerAppsObject.IsNull() && !powerAppsObject.IsUnknown() {
		var powerAppsSettings PowerAppsSettings
		powerAppsObject.(basetypes.ObjectValue).As(ctx, &powerAppsSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &PowerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.PowerApps = &PowerAppsSettingsDto{}
		if !powerAppsSettings.DisableShareWithEveryone.IsNull() && !powerAppsSettings.DisableShareWithEveryone.IsUnknown() {
			tenantSettingsDto.PowerPlatform.PowerApps.DisableShareWithEveryone = powerAppsSettings.DisableShareWithEveryone.ValueBoolPointer()
		}
		if !powerAppsSettings.EnableGuestsToMake.IsNull() && !powerAppsSettings.EnableGuestsToMake.IsUnknown() {
			tenantSettingsDto.PowerPlatform.PowerApps.EnableGuestsToMake = powerAppsSettings.EnableGuestsToMake.ValueBoolPointer()
		}
		if !powerAppsSettings.DisableMembersIndicator.IsNull() && !powerAppsSettings.DisableMembersIndicator.IsUnknown() {
			tenantSettingsDto.PowerPlatform.PowerApps.DisableMembersIndicator = powerAppsSettings.DisableMembersIndicator.ValueBoolPointer()
		}
		if !powerAppsSettings.DisableMakerMatch.IsNull() && !powerAppsSettings.DisableMakerMatch.IsUnknown() {
			tenantSettingsDto.PowerPlatform.PowerApps.DisableMakerMatch = powerAppsSettings.DisableMakerMatch.ValueBoolPointer()
		}
		if !powerAppsSettings.DisableUnusedLicenseAssignment.IsNull() && !powerAppsSettings.DisableUnusedLicenseAssignment.IsUnknown() {
			tenantSettingsDto.PowerPlatform.PowerApps.DisableUnusedLicenseAssignment = powerAppsSettings.DisableUnusedLicenseAssignment.ValueBoolPointer()
		}
		if !powerAppsSettings.DisableCreateFromImage.IsNull() && !powerAppsSettings.DisableCreateFromImage.IsUnknown() {
			tenantSettingsDto.PowerPlatform.PowerApps.DisableCreateFromImage = powerAppsSettings.DisableCreateFromImage.ValueBoolPointer()
		}
		if !powerAppsSettings.DisableCreateFromFigma.IsNull() && !powerAppsSettings.DisableCreateFromFigma.IsUnknown() {
			tenantSettingsDto.PowerPlatform.PowerApps.DisableCreateFromFigma = powerAppsSettings.DisableCreateFromFigma.ValueBoolPointer()
		}
		if !powerAppsSettings.DisableConnectionSharingWithEveryone.IsNull() && !powerAppsSettings.DisableConnectionSharingWithEveryone.IsUnknown() {
			tenantSettingsDto.PowerPlatform.PowerApps.DisableConnectionSharingWithEveryone = powerAppsSettings.DisableConnectionSharingWithEveryone.ValueBoolPointer()
		}
	}

	powerAutomateObject := powerPlatformAttributes["power_automate"]
	if !powerAutomateObject.IsNull() && !powerAutomateObject.IsUnknown() {
		var powerAutomateSettings PowerAutomateSettings
		powerAutomateObject.(basetypes.ObjectValue).As(ctx, &powerAutomateSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &PowerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.PowerAutomate = &PowerAutomateSettingsDto{}
		if !powerAutomateSettings.DisableCopilot.IsNull() && !powerAutomateSettings.DisableCopilot.IsUnknown() {
			tenantSettingsDto.PowerPlatform.PowerAutomate.DisableCopilot = powerAutomateSettings.DisableCopilot.ValueBoolPointer()
		}
	}

	environmentsObject := powerPlatformAttributes["environments"]
	if !environmentsObject.IsNull() && !environmentsObject.IsUnknown() {
		var environmentsSettings EnvironmentsSettings
		environmentsObject.(basetypes.ObjectValue).As(ctx, &environmentsSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &PowerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.Environments = &EnvironmentSettingsDto{}
		if !environmentsSettings.DisablePreferredDataLocationForTeamsEnvironment.IsNull() && !environmentsSettings.DisablePreferredDataLocationForTeamsEnvironment.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Environments.DisablePreferredDataLocationForTeamsEnvironment = environmentsSettings.DisablePreferredDataLocationForTeamsEnvironment.ValueBoolPointer()
		}
	}

	governanceObject := powerPlatformAttributes["governance"]
	if !governanceObject.IsNull() && !governanceObject.IsUnknown() {
		var governanceSettings GovernanceSettings
		governanceObject.(basetypes.ObjectValue).As(ctx, &governanceSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &PowerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.Governance = &GovernanceSettingsDto{}

		if !governanceSettings.DisableAdminDigest.IsNull() && !governanceSettings.DisableAdminDigest.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Governance.DisableAdminDigest = governanceSettings.DisableAdminDigest.ValueBoolPointer()
		}
		if !governanceSettings.DisableDeveloperEnvironmentCreationByNonAdminUsers.IsNull() && !governanceSettings.DisableDeveloperEnvironmentCreationByNonAdminUsers.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Governance.DisableDeveloperEnvironmentCreationByNonAdminUsers = governanceSettings.DisableDeveloperEnvironmentCreationByNonAdminUsers.ValueBoolPointer()
		}
		if !governanceSettings.EnableDefaultEnvironmentRouting.IsNull() && !governanceSettings.EnableDefaultEnvironmentRouting.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Governance.EnableDefaultEnvironmentRouting = governanceSettings.EnableDefaultEnvironmentRouting.ValueBoolPointer()
		}
		policyObject := governanceSettings.Policy
		if !policyObject.IsNull() && !policyObject.IsUnknown() {
			var policySettings PolicySettings
			policyObject.As(ctx, &policySettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

			tenantSettingsDto.PowerPlatform.Governance.Policy = &PolicySettingsDto{}
			if !policySettings.EnableDesktopFlowDataPolicyManagement.IsNull() && !policySettings.EnableDesktopFlowDataPolicyManagement.IsUnknown() {
				tenantSettingsDto.PowerPlatform.Governance.Policy.EnableDesktopFlowDataPolicyManagement = policySettings.EnableDesktopFlowDataPolicyManagement.ValueBoolPointer()
			}
		}
	}

	licensingObject := powerPlatformAttributes["licensing"]
	if !licensingObject.IsNull() && !licensingObject.IsUnknown() {
		var licensingSettings LicensingSettings
		licensingObject.(basetypes.ObjectValue).As(ctx, &licensingSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &PowerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.Licensing = &LicenseSettingsDto{}
		if !licensingSettings.DisableBillingPolicyCreationByNonAdminUsers.IsNull() && !licensingSettings.DisableBillingPolicyCreationByNonAdminUsers.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Licensing.DisableBillingPolicyCreationByNonAdminUsers = licensingSettings.DisableBillingPolicyCreationByNonAdminUsers.ValueBoolPointer()
		}
		if !licensingSettings.EnableTenantCapacityReportForEnvironmentAdmins.IsNull() && !licensingSettings.EnableTenantCapacityReportForEnvironmentAdmins.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Licensing.EnableTenantCapacityReportForEnvironmentAdmins = licensingSettings.EnableTenantCapacityReportForEnvironmentAdmins.ValueBoolPointer()
		}
		if !licensingSettings.StorageCapacityConsumptionWarningThreshold.IsNull() && !licensingSettings.StorageCapacityConsumptionWarningThreshold.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Licensing.StorageCapacityConsumptionWarningThreshold = licensingSettings.StorageCapacityConsumptionWarningThreshold.ValueInt64Pointer()
		}
		if !licensingSettings.EnableTenantLicensingReportForEnvironmentAdmins.IsNull() && !licensingSettings.EnableTenantLicensingReportForEnvironmentAdmins.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Licensing.EnableTenantLicensingReportForEnvironmentAdmins = licensingSettings.EnableTenantLicensingReportForEnvironmentAdmins.ValueBoolPointer()
		}
		if !licensingSettings.DisableUseOfUnassignedAIBuilderCredits.IsNull() && !licensingSettings.DisableUseOfUnassignedAIBuilderCredits.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Licensing.DisableUseOfUnassignedAIBuilderCredits = licensingSettings.DisableUseOfUnassignedAIBuilderCredits.ValueBoolPointer()
		}
	}

	powerPagesObject := powerPlatformAttributes["power_pages"]
	if !powerPagesObject.IsNull() && !powerPagesObject.IsUnknown() {
		var powerPagesSettings PowerPagesSettings
		powerPagesObject.(basetypes.ObjectValue).As(ctx, &powerPagesSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &PowerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.PowerPages = &PowerPagesSettingsDto{}
	}

	championsObject := powerPlatformAttributes["champions"]
	if !championsObject.IsNull() && !championsObject.IsUnknown() {
		var championsSettings ChampionsSettings
		championsObject.(basetypes.ObjectValue).As(ctx, &championsSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &PowerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.Champions = &ChampionSettingsDto{}
		if !championsSettings.DisableChampionsInvitationReachout.IsNull() && !championsSettings.DisableChampionsInvitationReachout.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Champions.DisableChampionsInvitationReachout = championsSettings.DisableChampionsInvitationReachout.ValueBoolPointer()
		}
		if !championsSettings.DisableSkillsMatchInvitationReachout.IsNull() && !championsSettings.DisableSkillsMatchInvitationReachout.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Champions.DisableSkillsMatchInvitationReachout = championsSettings.DisableSkillsMatchInvitationReachout.ValueBoolPointer()
		}
	}

	intelligenceObject := powerPlatformAttributes["intelligence"]
	if !intelligenceObject.IsNull() && !intelligenceObject.IsUnknown() {
		var intelligenceSettings IntelligenceSettings
		intelligenceObject.(basetypes.ObjectValue).As(ctx, &intelligenceSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &PowerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.Intelligence = &IntelligenceSettingsDto{}
		if !intelligenceSettings.DisableCopilot.IsNull() && !intelligenceSettings.DisableCopilot.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Intelligence.DisableCopilot = intelligenceSettings.DisableCopilot.ValueBoolPointer()
		}
		if !intelligenceSettings.EnableOpenAiBotPublishing.IsNull() && !intelligenceSettings.EnableOpenAiBotPublishing.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Intelligence.EnableOpenAiBotPublishing = intelligenceSettings.EnableOpenAiBotPublishing.ValueBoolPointer()
		}
	}

	modelExperimentationObject := powerPlatformAttributes["model_experimentation"]
	if !modelExperimentationObject.IsNull() && !modelExperimentationObject.IsUnknown() {
		var modelExperimentationSettings ModelExperimentationSettings
		modelExperimentationObject.(basetypes.ObjectValue).As(ctx, &modelExperimentationSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		tenantSettingsDto.PowerPlatform.ModelExperimentation = &ModelExperimentationSettingsDto{}
		if !modelExperimentationSettings.EnableModelDataSharing.IsNull() && !modelExperimentationSettings.EnableModelDataSharing.IsUnknown() {
			tenantSettingsDto.PowerPlatform.ModelExperimentation.EnableModelDataSharing = modelExperimentationSettings.EnableModelDataSharing.ValueBoolPointer()
		}
		if !modelExperimentationSettings.DisableDataLogging.IsNull() && !modelExperimentationSettings.DisableDataLogging.IsUnknown() {
			tenantSettingsDto.PowerPlatform.ModelExperimentation.DisableDataLogging = modelExperimentationSettings.DisableDataLogging.ValueBoolPointer()
		}
	}

	catalogSettingsObject := powerPlatformAttributes["catalog_settings"]
	if !catalogSettingsObject.IsNull() && !catalogSettingsObject.IsUnknown() {
		var catalogSettings CatalogSettingsSettings
		catalogSettingsObject.(basetypes.ObjectValue).As(ctx, &catalogSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		tenantSettingsDto.PowerPlatform.CatalogSettings = &CatalogSettingsDto{}
		if !catalogSettings.PowerCatalogAudienceSetting.IsNull() && !catalogSettings.PowerCatalogAudienceSetting.IsUnknown() {
			tenantSettingsDto.PowerPlatform.CatalogSettings.PowerCatalogAudienceSetting = catalogSettings.PowerCatalogAudienceSetting.ValueStringPointer()
		}
	}

	userManagementSettingsObject := powerPlatformAttributes["user_management_settings"]
	if !userManagementSettingsObject.IsNull() && !userManagementSettingsObject.IsUnknown() {
		var userManagementSettings UserManagementSettings
		userManagementSettingsObject.(basetypes.ObjectValue).As(ctx, &userManagementSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		tenantSettingsDto.PowerPlatform.UserManagementSettings = &UserManagementSettingsDto{}
		if !userManagementSettings.EnableDeleteDisabledUserinAllEnvironments.IsNull() && !userManagementSettings.EnableDeleteDisabledUserinAllEnvironments.IsUnknown() {
			tenantSettingsDto.PowerPlatform.UserManagementSettings.EnableDeleteDisabledUserinAllEnvironments = userManagementSettings.EnableDeleteDisabledUserinAllEnvironments.ValueBoolPointer()
		}
	}

	return tenantSettingsDto
}

func ConvertFromTenantSettingsDto(tenantSettingsDto TenantSettingsDto) TenantSettingsSourceModel {
	tenantSettings := TenantSettingsSourceModel{
		Id:                         types.StringValue(""),
		WalkMeOptOut:               types.BoolValue(*tenantSettingsDto.WalkMeOptOut),
		DisableNPSCommentsReachout: types.BoolValue(*tenantSettingsDto.DisableNPSCommentsReachout),
		DisableNewsletterSendout:   types.BoolValue(*tenantSettingsDto.DisableNewsletterSendout),
		DisableEnvironmentCreationByNonAdminUsers:      types.BoolValue(*tenantSettingsDto.DisableEnvironmentCreationByNonAdminUsers),
		DisablePortalsCreationByNonAdminUsers:          types.BoolValue(*tenantSettingsDto.DisablePortalsCreationByNonAdminUsers),
		DisableSurveyFeedback:                          types.BoolValue(*tenantSettingsDto.DisableSurveyFeedback),
		DisableTrialEnvironmentCreationByNonAdminUsers: types.BoolValue(*tenantSettingsDto.DisableTrialEnvironmentCreationByNonAdminUsers),
		DisableCapacityAllocationByEnvironmentAdmins:   types.BoolValue(*tenantSettingsDto.DisableCapacityAllocationByEnvironmentAdmins),
		DisableSupportTicketsVisibleByAllUsers:         types.BoolValue(*tenantSettingsDto.DisableSupportTicketsVisibleByAllUsers),
	}

	attrTypesSearchProperties := map[string]attr.Type{
		"disable_docs_search":       types.BoolType,
		"disable_community_search":  types.BoolType,
		"disable_bing_video_search": types.BoolType,
	}

	attrValuesSearchProperties := map[string]attr.Value{
		"disable_docs_search":       types.BoolValue(*tenantSettingsDto.PowerPlatform.Search.DisableDocsSearch),
		"disable_community_search":  types.BoolValue(*tenantSettingsDto.PowerPlatform.Search.DisableCommunitySearch),
		"disable_bing_video_search": types.BoolValue(*tenantSettingsDto.PowerPlatform.Search.DisableBingVideoSearch),
	}

	attrTypesTeamsIntegrationProperties := map[string]attr.Type{
		"share_with_colleagues_user_limit": types.Int64Type,
	}

	attrValuesTeamsIntegrationProperties := map[string]attr.Value{
		"share_with_colleagues_user_limit": types.Int64Value(*tenantSettingsDto.PowerPlatform.TeamsIntegration.ShareWithColleaguesUserLimit),
	}

	attrTypesPowerAppsProperties := map[string]attr.Type{
		"disable_share_with_everyone":              types.BoolType,
		"enable_guests_to_make":                    types.BoolType,
		"disable_members_indicator":                types.BoolType,
		"disable_maker_match":                      types.BoolType,
		"disable_unused_license_assignment":        types.BoolType,
		"disable_create_from_image":                types.BoolType,
		"disable_create_from_figma":                types.BoolType,
		"disable_connection_sharing_with_everyone": types.BoolType,
	}

	attrValuesPowerAppsProperties := map[string]attr.Value{
		"disable_share_with_everyone":              types.BoolValue(*tenantSettingsDto.PowerPlatform.PowerApps.DisableShareWithEveryone),
		"enable_guests_to_make":                    types.BoolValue(*tenantSettingsDto.PowerPlatform.PowerApps.EnableGuestsToMake),
		"disable_members_indicator":                types.BoolValue(*tenantSettingsDto.PowerPlatform.PowerApps.DisableMembersIndicator),
		"disable_maker_match":                      types.BoolValue(*tenantSettingsDto.PowerPlatform.PowerApps.DisableMakerMatch),
		"disable_unused_license_assignment":        types.BoolValue(*tenantSettingsDto.PowerPlatform.PowerApps.DisableUnusedLicenseAssignment),
		"disable_create_from_image":                types.BoolValue(*tenantSettingsDto.PowerPlatform.PowerApps.DisableCreateFromImage),
		"disable_create_from_figma":                types.BoolValue(*tenantSettingsDto.PowerPlatform.PowerApps.DisableCreateFromFigma),
		"disable_connection_sharing_with_everyone": types.BoolValue(*tenantSettingsDto.PowerPlatform.PowerApps.DisableConnectionSharingWithEveryone),
	}

	attrTypesPowerAutomateProperties := map[string]attr.Type{
		"disable_copilot": types.BoolType,
	}

	attrValuesPowerAutomateProperties := map[string]attr.Value{
		"disable_copilot": types.BoolValue(*tenantSettingsDto.PowerPlatform.PowerAutomate.DisableCopilot),
	}

	attrTypesEnvironmentsProperties := map[string]attr.Type{
		"disable_preferred_data_location_for_teams_environment": types.BoolType,
	}

	attrValuesEnvironmentsProperties := map[string]attr.Value{
		"disable_preferred_data_location_for_teams_environment": types.BoolValue(*tenantSettingsDto.PowerPlatform.Environments.DisablePreferredDataLocationForTeamsEnvironment),
	}

	attrTypesGovernanceProperties := map[string]attr.Type{
		"disable_admin_digest": types.BoolType,
		"disable_developer_environment_creation_by_non_admin_users": types.BoolType,
		"enable_default_environment_routing":                        types.BoolType,
		"policy": types.ObjectType{AttrTypes: map[string]attr.Type{
			"enable_desktop_flow_data_policy_management": types.BoolType,
		}},
	}

	attrValuesGovernanceProperties := map[string]attr.Value{
		"disable_admin_digest": types.BoolValue(*tenantSettingsDto.PowerPlatform.Governance.DisableAdminDigest),
		"disable_developer_environment_creation_by_non_admin_users": types.BoolValue(*tenantSettingsDto.PowerPlatform.Governance.DisableDeveloperEnvironmentCreationByNonAdminUsers),
		"enable_default_environment_routing":                        types.BoolValue(*tenantSettingsDto.PowerPlatform.Governance.EnableDefaultEnvironmentRouting),
		"policy": types.ObjectValueMust(map[string]attr.Type{
			"enable_desktop_flow_data_policy_management": types.BoolType,
		}, map[string]attr.Value{
			"enable_desktop_flow_data_policy_management": types.BoolValue(*tenantSettingsDto.PowerPlatform.Governance.Policy.EnableDesktopFlowDataPolicyManagement),
		}),
	}

	attrTypesLicencingProperties := map[string]attr.Type{
		"disable_billing_policy_creation_by_non_admin_users":    types.BoolType,
		"enable_tenant_capacity_report_for_environment_admins":  types.BoolType,
		"storage_capacity_consumption_warning_threshold":        types.Int64Type,
		"enable_tenant_licensing_report_for_environment_admins": types.BoolType,
		"disable_use_of_unassigned_ai_builder_credits":          types.BoolType,
	}

	attrValuesLicencingProperties := map[string]attr.Value{
		"disable_billing_policy_creation_by_non_admin_users":    types.BoolValue(*tenantSettingsDto.PowerPlatform.Licensing.DisableBillingPolicyCreationByNonAdminUsers),
		"enable_tenant_capacity_report_for_environment_admins":  types.BoolValue(*tenantSettingsDto.PowerPlatform.Licensing.EnableTenantCapacityReportForEnvironmentAdmins),
		"storage_capacity_consumption_warning_threshold":        types.Int64Value(*tenantSettingsDto.PowerPlatform.Licensing.StorageCapacityConsumptionWarningThreshold),
		"enable_tenant_licensing_report_for_environment_admins": types.BoolValue(*tenantSettingsDto.PowerPlatform.Licensing.EnableTenantLicensingReportForEnvironmentAdmins),
		"disable_use_of_unassigned_ai_builder_credits":          types.BoolValue(*tenantSettingsDto.PowerPlatform.Licensing.DisableUseOfUnassignedAIBuilderCredits),
	}

	attrTypesPowerPagesProperties := map[string]attr.Type{}

	attrValuesPowerPagesProperties := map[string]attr.Value{}

	attrTypesChampionsProperties := map[string]attr.Type{
		"disable_champions_invitation_reachout":    types.BoolType,
		"disable_skills_match_invitation_reachout": types.BoolType,
	}

	attrValuesChampionsProperties := map[string]attr.Value{
		"disable_champions_invitation_reachout":    types.BoolValue(*tenantSettingsDto.PowerPlatform.Champions.DisableChampionsInvitationReachout),
		"disable_skills_match_invitation_reachout": types.BoolValue(*tenantSettingsDto.PowerPlatform.Champions.DisableSkillsMatchInvitationReachout),
	}

	attrTypesIntelligenceProperties := map[string]attr.Type{
		"disable_copilot":               types.BoolType,
		"enable_open_ai_bot_publishing": types.BoolType,
	}

	attrValuesIntelligenceProperties := map[string]attr.Value{
		"disable_copilot":               types.BoolValue(*tenantSettingsDto.PowerPlatform.Intelligence.DisableCopilot),
		"enable_open_ai_bot_publishing": types.BoolValue(*tenantSettingsDto.PowerPlatform.Intelligence.EnableOpenAiBotPublishing),
	}

	attrTypesModelExperimentationProperties := map[string]attr.Type{
		"enable_model_data_sharing": types.BoolType,
		"disable_data_logging":      types.BoolType,
	}

	attrValuesModelExperimentationProperties := map[string]attr.Value{
		"enable_model_data_sharing": types.BoolValue(*tenantSettingsDto.PowerPlatform.ModelExperimentation.EnableModelDataSharing),
		"disable_data_logging":      types.BoolValue(*tenantSettingsDto.PowerPlatform.ModelExperimentation.DisableDataLogging),
	}

	attrTypesCatalogSettingsProperties := map[string]attr.Type{
		"power_catalog_audience_setting": types.StringType,
	}

	attrValuesCatalogSettingsProperties := map[string]attr.Value{
		"power_catalog_audience_setting": types.StringValue(*tenantSettingsDto.PowerPlatform.CatalogSettings.PowerCatalogAudienceSetting),
	}

	attrTypesUserManagementSettings := map[string]attr.Type{
		"enable_delete_disabled_user_in_all_environments": types.BoolType,
	}

	attrValuesUserManagementSettings := map[string]attr.Value{
		"enable_delete_disabled_user_in_all_environments": types.BoolValue(*tenantSettingsDto.PowerPlatform.UserManagementSettings.EnableDeleteDisabledUserinAllEnvironments),
	}

	attrTypesPowerPlatformObject := map[string]attr.Type{
		"search":                   types.ObjectType{AttrTypes: attrTypesSearchProperties},
		"teams_integration":        types.ObjectType{AttrTypes: attrTypesTeamsIntegrationProperties},
		"power_apps":               types.ObjectType{AttrTypes: attrTypesPowerAppsProperties},
		"power_automate":           types.ObjectType{AttrTypes: attrTypesPowerAutomateProperties},
		"environments":             types.ObjectType{AttrTypes: attrTypesEnvironmentsProperties},
		"governance":               types.ObjectType{AttrTypes: attrTypesGovernanceProperties},
		"licensing":                types.ObjectType{AttrTypes: attrTypesLicencingProperties},
		"power_pages":              types.ObjectType{AttrTypes: attrTypesPowerPagesProperties},
		"champions":                types.ObjectType{AttrTypes: attrTypesChampionsProperties},
		"intelligence":             types.ObjectType{AttrTypes: attrTypesIntelligenceProperties},
		"model_experimentation":    types.ObjectType{AttrTypes: attrTypesModelExperimentationProperties},
		"catalog_settings":         types.ObjectType{AttrTypes: attrTypesCatalogSettingsProperties},
		"user_management_settings": types.ObjectType{AttrTypes: attrTypesUserManagementSettings},
	}

	attrValuesPowerPlatformObject := map[string]attr.Value{
		"search":                   types.ObjectValueMust(attrTypesSearchProperties, attrValuesSearchProperties),
		"teams_integration":        types.ObjectValueMust(attrTypesTeamsIntegrationProperties, attrValuesTeamsIntegrationProperties),
		"power_apps":               types.ObjectValueMust(attrTypesPowerAppsProperties, attrValuesPowerAppsProperties),
		"power_automate":           types.ObjectValueMust(attrTypesPowerAutomateProperties, attrValuesPowerAutomateProperties),
		"environments":             types.ObjectValueMust(attrTypesEnvironmentsProperties, attrValuesEnvironmentsProperties),
		"governance":               types.ObjectValueMust(attrTypesGovernanceProperties, attrValuesGovernanceProperties),
		"licensing":                types.ObjectValueMust(attrTypesLicencingProperties, attrValuesLicencingProperties),
		"power_pages":              types.ObjectValueMust(attrTypesPowerPagesProperties, attrValuesPowerPagesProperties),
		"champions":                types.ObjectValueMust(attrTypesChampionsProperties, attrValuesChampionsProperties),
		"intelligence":             types.ObjectValueMust(attrTypesIntelligenceProperties, attrValuesIntelligenceProperties),
		"model_experimentation":    types.ObjectValueMust(attrTypesModelExperimentationProperties, attrValuesModelExperimentationProperties),
		"catalog_settings":         types.ObjectValueMust(attrTypesCatalogSettingsProperties, attrValuesCatalogSettingsProperties),
		"user_management_settings": types.ObjectValueMust(attrTypesUserManagementSettings, attrValuesUserManagementSettings),
	}

	tenantSettings.PowerPlatform = types.ObjectValueMust(attrTypesPowerPlatformObject, attrValuesPowerPlatformObject)
	return tenantSettings
}
