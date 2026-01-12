// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package tenant_settings

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/terraform-provider-power-platform/internal/customtypes"
	"github.com/microsoft/terraform-provider-power-platform/internal/helpers"
)

// Helper functions to reduce duplication in DTO conversion logic

// convertSimpleSettings creates an ObjectType and ObjectValue for simple attribute mappings.
// It handles the common pattern of nil checking and attribute conversion.
func convertSimpleSettings(attrTypes map[string]attr.Type, attrValues map[string]attr.Value) (basetypes.ObjectType, basetypes.ObjectValue) {
	objType := types.ObjectType{AttrTypes: attrTypes}
	if attrValues == nil {
		return objType, types.ObjectNull(attrTypes)
	}
	return objType, types.ObjectValueMust(attrTypes, attrValues)
}

// convertBoolSettingsMap creates attribute types map for boolean settings.
func convertBoolSettingsMap(keys []string) map[string]attr.Type {
	attrTypes := make(map[string]attr.Type, len(keys))
	for _, key := range keys {
		attrTypes[key] = types.BoolType
	}
	return attrTypes
}

// convertStringSettingsMap creates attribute types map for string settings.
func convertStringSettingsMap(keys []string) map[string]attr.Type {
	attrTypes := make(map[string]attr.Type, len(keys))
	for _, key := range keys {
		attrTypes[key] = types.StringType
	}
	return attrTypes
}

type tenantDto struct {
	TenantId                         string `json:"tenantId,omitempty"`
	State                            string `json:"state,omitempty"`
	Location                         string `json:"location,omitempty"`
	AadCountryGeo                    string `json:"aadCountryGeo,omitempty"`
	DataStorageGeo                   string `json:"dataStorageGeo,omitempty"`
	AadDataBoundary                  string `json:"aadDataBoundary,omitempty"`
	FedRAMPHighCertificationRequired bool   `json:"fedRAMPHighCertificationRequired,omitempty"`
}

type powerPlatformSettingsDto struct {
	Search                 *searchSettingsDto               `json:"search,omitempty"`
	TeamsIntegration       *teamsIntegrationSettingsDto     `json:"teamsIntegration,omitempty"`
	PowerApps              *powerAppsSettingsDto            `json:"powerApps,omitempty"`
	PowerAutomate          *powerAutomateSettingsDto        `json:"powerAutomate,omitempty"`
	Environments           *environmentSettingsDto          `json:"environments,omitempty"`
	Governance             *governanceSettingsDto           `json:"governance,omitempty"`
	Licensing              *licenseSettingsDto              `json:"licensing,omitempty"`
	PowerPages             *powerPagesSettingsDto           `json:"powerPages,omitempty"`
	Champions              *championSettingsDto             `json:"champions,omitempty"`
	Intelligence           *intelligenceSettingsDto         `json:"intelligence,omitempty"`
	ModelExperimentation   *modelExperimentationSettingsDto `json:"modelExperimentation,omitempty"`
	CatalogSettings        *catalogSettingsDto              `json:"catalogSettings,omitempty"`
	UserManagementSettings *userManagementSettingsDto       `json:"userManagementSettings,omitempty"`
}

type teamsIntegrationSettingsDto struct {
	ShareWithColleaguesUserLimit *int64 `json:"shareWithColleaguesUserLimit,omitempty"`
}

type powerAutomateSettingsDto struct {
	DisableCopilot          *bool `json:"disableCopilot,omitempty"`
	DisableCopilotWithBing  *bool `json:"disableCopilotWithBing,omitempty"`
	AllowUseOfHostedBrowser *bool `json:"enableComputerUseSharedMachines,omitempty"`
	DisableFlowResubmission *bool `json:"disableFlowRunResubmission,omitempty"`
}

type powerAppsSettingsDto struct {
	DisableCopilot                       *bool `json:"disableCopilot,omitempty"`
	DisableShareWithEveryone             *bool `json:"disableShareWithEveryone,omitempty"`
	EnableGuestsToMake                   *bool `json:"enableGuestsToMake,omitempty"`
	DisableMakerMatch                    *bool `json:"disableMakerMatch,omitempty"`
	DisableUnusedLicenseAssignment       *bool `json:"disableUnusedLicenseAssignment,omitempty"`
	DisableCreateFromImage               *bool `json:"disableCreateFromImage,omitempty"`
	DisableCreateFromFigma               *bool `json:"disableCreateFromFigma,omitempty"`
	DisableConnectionSharingWithEveryone *bool `json:"disableConnectionSharingWithEveryone,omitempty"`
	EnableCanvasAppInsights              *bool `json:"enableCanvasAppInsights,omitempty"`
}

type environmentSettingsDto struct {
	DisablePreferredDataLocationForTeamsEnvironment *bool `json:"disablePreferredDataLocationForTeamsEnvironment,omitempty"`
}

type searchSettingsDto struct {
	DisableDocsSearch      *bool `json:"disableDocsSearch,omitempty"`
	DisableCommunitySearch *bool `json:"disableCommunitySearch,omitempty"`
	DisableBingVideoSearch *bool `json:"disableBingVideoSearch,omitempty"`
}

type governanceSettingsDto struct {
	AdditionalAdminDigestEmailRecipients               *string            `json:"additionalAdminDigestEmailRecipients,omitempty"`
	DisableAdminDigest                                 *bool              `json:"disableAdminDigest,omitempty"`
	DisableDeveloperEnvironmentCreationByNonAdminUsers *bool              `json:"disableDeveloperEnvironmentCreationByNonAdminUsers,omitempty"`
	EnableDefaultEnvironmentRouting                    *bool              `json:"enableDefaultEnvironmentRouting,omitempty"`
	EnvironmentRoutingAllMakers                        *bool              `json:"environmentRoutingAllMakers,omitempty"`
	EnvironmentRoutingTargetEnvironmentGroupId         *string            `json:"environmentRoutingTargetEnvironmentGroupId,omitempty"`
	EnvironmentRoutingTargetSecurityGroupId            *string            `json:"environmentRoutingTargetSecurityGroupId,omitempty"`
	Policy                                             *policySettingsDto `json:"policy,omitempty"`
}

type policySettingsDto struct {
	EnableDesktopFlowDataPolicyManagement *bool `json:"enableDesktopFlowDataPolicyManagement,omitempty"`
}

type licenseSettingsDto struct {
	DisableBillingPolicyCreationByNonAdminUsers     *bool  `json:"disableBillingPolicyCreationByNonAdminUsers,omitempty"`
	EnableTenantCapacityReportForEnvironmentAdmins  *bool  `json:"enableTenantCapacityReportForEnvironmentAdmins,omitempty"`
	StorageCapacityConsumptionWarningThreshold      *int64 `json:"storageCapacityConsumptionWarningThreshold,omitempty"`
	EnableTenantLicensingReportForEnvironmentAdmins *bool  `json:"enableTenantLicensingReportForEnvironmentAdmins,omitempty"`
	DisableUseOfUnassignedAIBuilderCredits          *bool  `json:"disableUseOfUnassignedAIBuilderCredits,omitempty"`
	ApplyAutoClaimToOnlyManagedEnvironments         *bool  `json:"applyAutoClaimToOnlyManagedEnvironments,omitempty"`
	ApplyPAutoAutoClaimToOnlyManagedEnvironments    *bool  `json:"applyPAutoAutoClaimToOnlyManagedEnvironments,omitempty"`
}

type powerPagesSettingsDto struct {
}

type championSettingsDto struct {
	DisableChampionsInvitationReachout   *bool `json:"disableChampionsInvitationReachout,omitempty"`
	DisableSkillsMatchInvitationReachout *bool `json:"disableSkillsMatchInvitationReachout,omitempty"`
}

type intelligenceSettingsDto struct {
	DisableCopilot            *bool `json:"disableCopilot,omitempty"`
	EnableOpenAiBotPublishing *bool `json:"enableOpenAiBotPublishing,omitempty"`
	BasicCopilotFeedback      *bool `json:"disableCopilotFeedback,omitempty"`
	AdditionalCopilotFeedback *bool `json:"disableCopilotFeedbackMetadata,omitempty"`
}

type modelExperimentationSettingsDto struct {
	EnableModelDataSharing *bool `json:"enableModelDataSharing,omitempty"`
	DisableDataLogging     *bool `json:"disableDataLogging,omitempty"`
}

type catalogSettingsDto struct {
	PowerCatalogAudienceSetting *string `json:"powerCatalogAudienceSetting,omitempty"`
}

type userManagementSettingsDto struct {
	EnableDeleteDisabledUserinAllEnvironments *bool `json:"enableDeleteDisabledUserinAllEnvironments,omitempty"`
}

type tenantSettingsDto struct {
	WalkMeOptOut                                   *bool `json:"walkMeOptOut,omitempty"`
	DisableNewsletterSendout                       *bool `json:"disableNewsletterSendout,omitempty"`
	DisableEnvironmentCreationByNonAdminUsers      *bool `json:"disableEnvironmentCreationByNonAdminUsers,omitempty"`
	DisablePortalsCreationByNonAdminUsers          *bool `json:"disablePortalsCreationByNonAdminUsers,omitempty"`
	DisableTrialEnvironmentCreationByNonAdminUsers *bool `json:"disableTrialEnvironmentCreationByNonAdminUsers,omitempty"`
	DisableCapacityAllocationByEnvironmentAdmins   *bool `json:"disableCapacityAllocationByEnvironmentAdmins,omitempty"`
	DisableSupportTicketsVisibleByAllUsers         *bool `json:"disableSupportTicketsVisibleByAllUsers,omitempty"`

	DisableNPSCommentsReachout   *bool `json:"disableNPSCommentsReachout,omitempty"`
	DisableSurveyFeedback        *bool `json:"disableSurveyFeedback,omitempty"`
	DisableUserInitiatedFeedback *bool `json:"disableUserInitiatedFeedback,omitempty"`
	DisableSurveyScreenshots     *bool `json:"disableSurveyScreenshots,omitempty"`

	PowerPlatform *powerPlatformSettingsDto `json:"powerPlatform,omitempty"`
}

func convertFromTenantSettingsModel(ctx context.Context, tenantSettings TenantSettingsResourceModel) (tenantSettingsDto, error) {
	tenantSettingsDto := tenantSettingsDto{}

	tenantSettingsDto.WalkMeOptOut = helpers.BoolPointer(tenantSettings.WalkMeOptOut)
	tenantSettingsDto.DisableNewsletterSendout = helpers.BoolPointer(tenantSettings.DisableNewsletterSendout)
	tenantSettingsDto.DisableEnvironmentCreationByNonAdminUsers = helpers.BoolPointer(tenantSettings.DisableEnvironmentCreationByNonAdminUsers)
	tenantSettingsDto.DisablePortalsCreationByNonAdminUsers = helpers.BoolPointer(tenantSettings.DisablePortalsCreationByNonAdminUsers)
	tenantSettingsDto.DisableTrialEnvironmentCreationByNonAdminUsers = helpers.BoolPointer(tenantSettings.DisableTrialEnvironmentCreationByNonAdminUsers)
	tenantSettingsDto.DisableCapacityAllocationByEnvironmentAdmins = helpers.BoolPointer(tenantSettings.DisableCapacityAllocationByEnvironmentAdmins)
	tenantSettingsDto.DisableSupportTicketsVisibleByAllUsers = helpers.BoolPointer(tenantSettings.DisableSupportTicketsVisibleByAllUsers)

	if !tenantSettings.PowerPlatform.IsNull() && !tenantSettings.PowerPlatform.IsUnknown() {
		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &powerPlatformSettingsDto{}
		}
		powerPlatformAttributes := tenantSettings.PowerPlatform.Attributes()
		convertProductFeedbackModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
		convertSearchModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
		convertTeamsIntegrationModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
		convertPowerAppsModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
		convertPowerAutomateModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
		convertEnvironmentsModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
		convertGovernanceModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
		convertLicensingModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
		convertPowerPagesModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
		convertChampionsModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
		convertIntelligenceModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
		convertModelExperimentationModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
		convertCatalogSettingsModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
		convertUserManagementSettingsModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
	}
	return tenantSettingsDto, nil
}

func convertProductFeedbackModel(ctx context.Context, powerPlatformAttributes map[string]attr.Value, tenantSettingsDto *tenantSettingsDto) {
	productFeedbackObject := powerPlatformAttributes["product_feedback"]
	if productFeedbackObject != nil && !productFeedbackObject.IsNull() && !productFeedbackObject.IsUnknown() {
		objectValue, ok := productFeedbackObject.(basetypes.ObjectValue)
		if !ok {
			tflog.Debug(ctx, "failed to convert product_feedback to ObjectValue", map[string]any{
				"expected_type": "basetypes.ObjectValue",
				"actual_type":   fmt.Sprintf("%T", productFeedbackObject),
			})
			return // Skip conversion if type assertion fails
		}

		var productFeedbackSettings ProductFeedbackSettings
		objectValue.As(ctx, &productFeedbackSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &powerPlatformSettingsDto{}
		}

		if !productFeedbackSettings.DisableUserSurveyFeedback.IsNull() && !productFeedbackSettings.DisableUserSurveyFeedback.IsUnknown() {
			tenantSettingsDto.DisableSurveyFeedback = productFeedbackSettings.DisableUserSurveyFeedback.ValueBoolPointer()
		}
		if !productFeedbackSettings.DisableMicrosoftFollowUp.IsNull() && !productFeedbackSettings.DisableMicrosoftFollowUp.IsUnknown() {
			tenantSettingsDto.DisableNPSCommentsReachout = productFeedbackSettings.DisableMicrosoftFollowUp.ValueBoolPointer()
		}
		if !productFeedbackSettings.DisableAttachments.IsNull() && !productFeedbackSettings.DisableAttachments.IsUnknown() {
			tenantSettingsDto.DisableSurveyScreenshots = productFeedbackSettings.DisableAttachments.ValueBoolPointer()
		}
		if !productFeedbackSettings.DisableMicrosoftSurveysSend.IsNull() && !productFeedbackSettings.DisableMicrosoftSurveysSend.IsUnknown() {
			tenantSettingsDto.DisableUserInitiatedFeedback = productFeedbackSettings.DisableMicrosoftSurveysSend.ValueBoolPointer()
		}
	}
}

func convertSearchModel(ctx context.Context, powerPlatformAttributes map[string]attr.Value, tenantSettingsDto *tenantSettingsDto) {
	searchObject := powerPlatformAttributes["search"]
	if searchObject != nil && !searchObject.IsNull() && !searchObject.IsUnknown() {
		objectValue, ok := searchObject.(basetypes.ObjectValue)
		if !ok {
			tflog.Debug(ctx, "failed to convert search to ObjectValue", map[string]any{
				"expected_type": "basetypes.ObjectValue",
				"actual_type":   fmt.Sprintf("%T", searchObject),
			})
			return // Skip conversion if type assertion fails
		}

		var searchSettings SearchSettingsModel
		objectValue.As(ctx, &searchSettings, basetypes.ObjectAsOptions{
			UnhandledNullAsEmpty:    true,
			UnhandledUnknownAsEmpty: true,
		})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &powerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.Search = &searchSettingsDto{}
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
}

func convertTeamsIntegrationModel(ctx context.Context, powerPlatformAttributes map[string]attr.Value, tenantSettingsDto *tenantSettingsDto) {
	teamIntegrationObject := powerPlatformAttributes["teams_integration"]
	if teamIntegrationObject != nil && !teamIntegrationObject.IsNull() && !teamIntegrationObject.IsUnknown() {
		objectValue, ok := teamIntegrationObject.(basetypes.ObjectValue)
		if !ok {
			tflog.Debug(ctx, "failed to convert teams_integration to ObjectValue", map[string]any{
				"expected_type": "basetypes.ObjectValue",
				"actual_type":   fmt.Sprintf("%T", teamIntegrationObject),
			})
			return // Skip conversion if type assertion fails
		}

		var teamsIntegrationSettings TeamsIntegrationSettings
		objectValue.As(ctx, &teamsIntegrationSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &powerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.TeamsIntegration = &teamsIntegrationSettingsDto{}
		if !teamsIntegrationSettings.ShareWithColleaguesUserLimit.IsNull() && !teamsIntegrationSettings.ShareWithColleaguesUserLimit.IsUnknown() {
			tenantSettingsDto.PowerPlatform.TeamsIntegration.ShareWithColleaguesUserLimit = teamsIntegrationSettings.ShareWithColleaguesUserLimit.ValueInt64Pointer()
		}
	}
}

func convertPowerAppsModel(ctx context.Context, powerPlatformAttributes map[string]attr.Value, tenantSettingsDto *tenantSettingsDto) {
	powerAppsObject := powerPlatformAttributes["power_apps"]
	if powerAppsObject != nil && !powerAppsObject.IsNull() && !powerAppsObject.IsUnknown() {
		objectValue, ok := powerAppsObject.(basetypes.ObjectValue)
		if !ok {
			tflog.Debug(ctx, "failed to convert power_apps to ObjectValue", map[string]any{
				"expected_type": "basetypes.ObjectValue",
				"actual_type":   fmt.Sprintf("%T", powerAppsObject),
			})
			return // Skip conversion if type assertion fails
		}

		var powerAppsSettings PowerAppsSettings
		objectValue.As(ctx, &powerAppsSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &powerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.PowerApps = &powerAppsSettingsDto{}
		if !powerAppsSettings.DisableShareWithEveryone.IsNull() && !powerAppsSettings.DisableShareWithEveryone.IsUnknown() {
			tenantSettingsDto.PowerPlatform.PowerApps.DisableShareWithEveryone = powerAppsSettings.DisableShareWithEveryone.ValueBoolPointer()
		}
		if !powerAppsSettings.EnableGuestsToMake.IsNull() && !powerAppsSettings.EnableGuestsToMake.IsUnknown() {
			tenantSettingsDto.PowerPlatform.PowerApps.EnableGuestsToMake = powerAppsSettings.EnableGuestsToMake.ValueBoolPointer()
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
		if !powerAppsSettings.EnableCanvasAppInsights.IsNull() && !powerAppsSettings.EnableCanvasAppInsights.IsUnknown() {
			tenantSettingsDto.PowerPlatform.PowerApps.EnableCanvasAppInsights = powerAppsSettings.EnableCanvasAppInsights.ValueBoolPointer()
		}
		if !powerAppsSettings.DisableCopilot.IsNull() && !powerAppsSettings.DisableCopilot.IsUnknown() {
			tenantSettingsDto.PowerPlatform.PowerApps.DisableCopilot = powerAppsSettings.DisableCopilot.ValueBoolPointer()
		}
	}
}

func convertPowerAutomateModel(ctx context.Context, powerPlatformAttributes map[string]attr.Value, tenantSettingsDto *tenantSettingsDto) {
	powerAutomateObject := powerPlatformAttributes["power_automate"]
	if powerAutomateObject != nil && !powerAutomateObject.IsNull() && !powerAutomateObject.IsUnknown() {
		objectValue, ok := powerAutomateObject.(basetypes.ObjectValue)
		if !ok {
			tflog.Debug(ctx, "failed to convert power_automate to ObjectValue", map[string]any{
				"expected_type": "basetypes.ObjectValue",
				"actual_type":   fmt.Sprintf("%T", powerAutomateObject),
			})
			return // Skip conversion if type assertion fails
		}

		var powerAutomateSettings PowerAutomateSettings
		objectValue.As(ctx, &powerAutomateSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &powerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.PowerAutomate = &powerAutomateSettingsDto{}
		if !powerAutomateSettings.DisableCopilot.IsNull() && !powerAutomateSettings.DisableCopilot.IsUnknown() {
			tenantSettingsDto.PowerPlatform.PowerAutomate.DisableCopilot = powerAutomateSettings.DisableCopilot.ValueBoolPointer()
		}
		if !powerAutomateSettings.DisableCopilotWithBing.IsNull() && !powerAutomateSettings.DisableCopilotWithBing.IsUnknown() {
			tenantSettingsDto.PowerPlatform.PowerAutomate.DisableCopilotWithBing = powerAutomateSettings.DisableCopilotWithBing.ValueBoolPointer()
		}
		if !powerAutomateSettings.AllowUseOfHostedBrowser.IsNull() && !powerAutomateSettings.AllowUseOfHostedBrowser.IsUnknown() {
			tenantSettingsDto.PowerPlatform.PowerAutomate.AllowUseOfHostedBrowser = powerAutomateSettings.AllowUseOfHostedBrowser.ValueBoolPointer()
		}
		if !powerAutomateSettings.DisableFlowResubmission.IsNull() && !powerAutomateSettings.DisableFlowResubmission.IsUnknown() {
			tenantSettingsDto.PowerPlatform.PowerAutomate.DisableFlowResubmission = powerAutomateSettings.DisableFlowResubmission.ValueBoolPointer()
		}
	}
}

func convertEnvironmentsModel(ctx context.Context, powerPlatformAttributes map[string]attr.Value, tenantSettingsDto *tenantSettingsDto) {
	environmentsObject := powerPlatformAttributes["environments"]
	if environmentsObject != nil && !environmentsObject.IsNull() && !environmentsObject.IsUnknown() {
		objectValue, ok := environmentsObject.(basetypes.ObjectValue)
		if !ok {
			tflog.Debug(ctx, "failed to convert environments to ObjectValue", map[string]any{
				"expected_type": "basetypes.ObjectValue",
				"actual_type":   fmt.Sprintf("%T", environmentsObject),
			})
			return // Skip conversion if type assertion fails
		}

		var environmentsSettings EnvironmentsSettings
		objectValue.As(ctx, &environmentsSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &powerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.Environments = &environmentSettingsDto{}
		if !environmentsSettings.DisablePreferredDataLocationForTeamsEnvironment.IsNull() && !environmentsSettings.DisablePreferredDataLocationForTeamsEnvironment.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Environments.DisablePreferredDataLocationForTeamsEnvironment = environmentsSettings.DisablePreferredDataLocationForTeamsEnvironment.ValueBoolPointer()
		}
	}
}

func convertGovernanceModel(ctx context.Context, powerPlatformAttributes map[string]attr.Value, tenantSettingsDto *tenantSettingsDto) {
	governanceObject := powerPlatformAttributes["governance"]
	if governanceObject != nil && !governanceObject.IsNull() && !governanceObject.IsUnknown() {
		objectValue, ok := governanceObject.(basetypes.ObjectValue)
		if !ok {
			tflog.Debug(ctx, "failed to convert governance to ObjectValue", map[string]any{
				"expected_type": "basetypes.ObjectValue",
				"actual_type":   fmt.Sprintf("%T", governanceObject),
			})
			return // Skip conversion if type assertion fails
		}

		var governanceSettings GovernanceSettings
		objectValue.As(ctx, &governanceSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &powerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.Governance = &governanceSettingsDto{}

		if helpers.IsKnown(governanceSettings.WeeklyDigestEmailRecipients) {
			value := strings.Join(helpers.SetToStringSlice(governanceSettings.WeeklyDigestEmailRecipients), ";")
			tenantSettingsDto.PowerPlatform.Governance.AdditionalAdminDigestEmailRecipients = &value
		}
		if !governanceSettings.DisableAdminDigest.IsNull() && !governanceSettings.DisableAdminDigest.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Governance.DisableAdminDigest = governanceSettings.DisableAdminDigest.ValueBoolPointer()
		}
		if !governanceSettings.DisableDeveloperEnvironmentCreationByNonAdminUsers.IsNull() && !governanceSettings.DisableDeveloperEnvironmentCreationByNonAdminUsers.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Governance.DisableDeveloperEnvironmentCreationByNonAdminUsers = governanceSettings.DisableDeveloperEnvironmentCreationByNonAdminUsers.ValueBoolPointer()
		}
		if !governanceSettings.EnableDefaultEnvironmentRouting.IsNull() && !governanceSettings.EnableDefaultEnvironmentRouting.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Governance.EnableDefaultEnvironmentRouting = governanceSettings.EnableDefaultEnvironmentRouting.ValueBoolPointer()
		}
		if !governanceSettings.EnvironmentRoutingAllMakers.IsNull() && !governanceSettings.EnvironmentRoutingAllMakers.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Governance.EnvironmentRoutingAllMakers = governanceSettings.EnvironmentRoutingAllMakers.ValueBoolPointer()
		}

		if !governanceSettings.EnvironmentRoutingTargetEnvironmentGroupId.IsNull() && !governanceSettings.EnvironmentRoutingTargetEnvironmentGroupId.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Governance.EnvironmentRoutingTargetEnvironmentGroupId = governanceSettings.EnvironmentRoutingTargetEnvironmentGroupId.ValueStringPointer()
		}

		if !governanceSettings.EnvironmentRoutingTargetSecurityGroupId.IsNull() && !governanceSettings.EnvironmentRoutingTargetSecurityGroupId.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Governance.EnvironmentRoutingTargetSecurityGroupId = governanceSettings.EnvironmentRoutingTargetSecurityGroupId.ValueStringPointer()
		}

		policyObject := governanceSettings.Policy
		if !policyObject.IsNull() && !policyObject.IsUnknown() {
			var policySettings PolicySettings
			policyObject.As(ctx, &policySettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

			tenantSettingsDto.PowerPlatform.Governance.Policy = &policySettingsDto{}
			if !policySettings.EnableDesktopFlowDataPolicyManagement.IsNull() && !policySettings.EnableDesktopFlowDataPolicyManagement.IsUnknown() {
				tenantSettingsDto.PowerPlatform.Governance.Policy.EnableDesktopFlowDataPolicyManagement = policySettings.EnableDesktopFlowDataPolicyManagement.ValueBoolPointer()
			}
		}
	}
}

func convertLicensingModel(ctx context.Context, powerPlatformAttributes map[string]attr.Value, tenantSettingsDto *tenantSettingsDto) {
	licensingObject := powerPlatformAttributes["licensing"]
	if licensingObject != nil && !licensingObject.IsNull() && !licensingObject.IsUnknown() {
		objectValue, ok := licensingObject.(basetypes.ObjectValue)
		if !ok {
			tflog.Debug(ctx, "failed to convert licensing to ObjectValue", map[string]any{
				"expected_type": "basetypes.ObjectValue",
				"actual_type":   fmt.Sprintf("%T", licensingObject),
			})
			return // Skip conversion if type assertion fails
		}

		var licensingSettings LicensingSettings
		objectValue.As(ctx, &licensingSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &powerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.Licensing = &licenseSettingsDto{}
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
		if !licensingSettings.ApplyAutoClaimPowerAppsToOnlyManagedEnvironments.IsNull() && !licensingSettings.ApplyAutoClaimPowerAppsToOnlyManagedEnvironments.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Licensing.ApplyAutoClaimToOnlyManagedEnvironments = licensingSettings.ApplyAutoClaimPowerAppsToOnlyManagedEnvironments.ValueBoolPointer()
		}
		if !licensingSettings.ApplyAutoClaimPowerAutomateToOnlyManagedEnvironments.IsNull() && !licensingSettings.ApplyAutoClaimPowerAutomateToOnlyManagedEnvironments.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Licensing.ApplyPAutoAutoClaimToOnlyManagedEnvironments = licensingSettings.ApplyAutoClaimPowerAutomateToOnlyManagedEnvironments.ValueBoolPointer()
		}
	}
}

func convertPowerPagesModel(ctx context.Context, powerPlatformAttributes map[string]attr.Value, tenantSettingsDto *tenantSettingsDto) {
	powerPagesObject := powerPlatformAttributes["power_pages"]
	if powerPagesObject != nil && !powerPagesObject.IsNull() && !powerPagesObject.IsUnknown() {
		objectValue, ok := powerPagesObject.(basetypes.ObjectValue)
		if !ok {
			tflog.Debug(ctx, "failed to convert power_pages to ObjectValue", map[string]any{
				"expected_type": "basetypes.ObjectValue",
				"actual_type":   fmt.Sprintf("%T", powerPagesObject),
			})
			return // Skip conversion if type assertion fails
		}

		var powerPagesSettings PowerPagesSettings
		objectValue.As(ctx, &powerPagesSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &powerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.PowerPages = &powerPagesSettingsDto{}
	}
}

func convertChampionsModel(ctx context.Context, powerPlatformAttributes map[string]attr.Value, tenantSettingsDto *tenantSettingsDto) {
	championsObject := powerPlatformAttributes["champions"]
	if championsObject != nil && !championsObject.IsNull() && !championsObject.IsUnknown() {
		objectValue, ok := championsObject.(basetypes.ObjectValue)
		if !ok {
			tflog.Debug(ctx, "failed to convert champions to ObjectValue", map[string]any{
				"expected_type": "basetypes.ObjectValue",
				"actual_type":   fmt.Sprintf("%T", championsObject),
			})
			return // Skip conversion if type assertion fails
		}

		var championsSettings ChampionsSettings
		objectValue.As(ctx, &championsSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &powerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.Champions = &championSettingsDto{}
		if !championsSettings.DisableChampionsInvitationReachout.IsNull() && !championsSettings.DisableChampionsInvitationReachout.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Champions.DisableChampionsInvitationReachout = championsSettings.DisableChampionsInvitationReachout.ValueBoolPointer()
		}
		if !championsSettings.DisableSkillsMatchInvitationReachout.IsNull() && !championsSettings.DisableSkillsMatchInvitationReachout.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Champions.DisableSkillsMatchInvitationReachout = championsSettings.DisableSkillsMatchInvitationReachout.ValueBoolPointer()
		}
	}
}

func convertIntelligenceModel(ctx context.Context, powerPlatformAttributes map[string]attr.Value, tenantSettingsDto *tenantSettingsDto) {
	intelligenceObject := powerPlatformAttributes["intelligence"]
	if intelligenceObject != nil && !intelligenceObject.IsNull() && !intelligenceObject.IsUnknown() {
		objectValue, ok := intelligenceObject.(basetypes.ObjectValue)
		if !ok {
			return // Skip conversion if type assertion fails
		}

		var intelligenceSettings IntelligenceSettings
		objectValue.As(ctx, &intelligenceSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if tenantSettingsDto.PowerPlatform == nil {
			tenantSettingsDto.PowerPlatform = &powerPlatformSettingsDto{}
		}
		tenantSettingsDto.PowerPlatform.Intelligence = &intelligenceSettingsDto{}
		if !intelligenceSettings.DisableCopilot.IsNull() && !intelligenceSettings.DisableCopilot.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Intelligence.DisableCopilot = intelligenceSettings.DisableCopilot.ValueBoolPointer()
		}
		if !intelligenceSettings.EnableOpenAiBotPublishing.IsNull() && !intelligenceSettings.EnableOpenAiBotPublishing.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Intelligence.EnableOpenAiBotPublishing = intelligenceSettings.EnableOpenAiBotPublishing.ValueBoolPointer()
		}
		if !intelligenceSettings.BasicCopilotFeedback.IsNull() && !intelligenceSettings.BasicCopilotFeedback.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Intelligence.BasicCopilotFeedback = intelligenceSettings.BasicCopilotFeedback.ValueBoolPointer()
		}
		if !intelligenceSettings.AdditionalCopilotFeedback.IsNull() && !intelligenceSettings.AdditionalCopilotFeedback.IsUnknown() {
			tenantSettingsDto.PowerPlatform.Intelligence.AdditionalCopilotFeedback = intelligenceSettings.AdditionalCopilotFeedback.ValueBoolPointer()
		}
	}
}

func convertModelExperimentationModel(ctx context.Context, powerPlatformAttributes map[string]attr.Value, tenantSettingsDto *tenantSettingsDto) {
	modelExperimentationObject := powerPlatformAttributes["model_experimentation"]
	if modelExperimentationObject != nil && !modelExperimentationObject.IsNull() && !modelExperimentationObject.IsUnknown() {
		objectValue, ok := modelExperimentationObject.(basetypes.ObjectValue)
		if !ok {
			tflog.Debug(ctx, "failed to convert model_experimentation to ObjectValue", map[string]any{
				"expected_type": "basetypes.ObjectValue",
				"actual_type":   fmt.Sprintf("%T", modelExperimentationObject),
			})
			return // Skip conversion if type assertion fails
		}

		var modelExperimentationSettings ModelExperimentationSettings
		objectValue.As(ctx, &modelExperimentationSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		tenantSettingsDto.PowerPlatform.ModelExperimentation = &modelExperimentationSettingsDto{}
		if !modelExperimentationSettings.EnableModelDataSharing.IsNull() && !modelExperimentationSettings.EnableModelDataSharing.IsUnknown() {
			tenantSettingsDto.PowerPlatform.ModelExperimentation.EnableModelDataSharing = modelExperimentationSettings.EnableModelDataSharing.ValueBoolPointer()
		}
		if !modelExperimentationSettings.DisableDataLogging.IsNull() && !modelExperimentationSettings.DisableDataLogging.IsUnknown() {
			tenantSettingsDto.PowerPlatform.ModelExperimentation.DisableDataLogging = modelExperimentationSettings.DisableDataLogging.ValueBoolPointer()
		}
	}
}

func convertCatalogSettingsModel(ctx context.Context, powerPlatformAttributes map[string]attr.Value, tenantSettingsDto *tenantSettingsDto) {
	catalogSettingsObject := powerPlatformAttributes["catalog_settings"]
	if catalogSettingsObject != nil && !catalogSettingsObject.IsNull() && !catalogSettingsObject.IsUnknown() {
		objectValue, ok := catalogSettingsObject.(basetypes.ObjectValue)
		if !ok {
			tflog.Debug(ctx, "failed to convert catalog_settings to ObjectValue", map[string]any{
				"expected_type": "basetypes.ObjectValue",
				"actual_type":   fmt.Sprintf("%T", catalogSettingsObject),
			})
			return // Skip conversion if type assertion fails
		}

		var catalogSettings CatalogSettingsSettings
		objectValue.As(ctx, &catalogSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		tenantSettingsDto.PowerPlatform.CatalogSettings = &catalogSettingsDto{}
		if !catalogSettings.PowerCatalogAudienceSetting.IsNull() && !catalogSettings.PowerCatalogAudienceSetting.IsUnknown() {
			tenantSettingsDto.PowerPlatform.CatalogSettings.PowerCatalogAudienceSetting = catalogSettings.PowerCatalogAudienceSetting.ValueStringPointer()
		}
	}
}

func convertUserManagementSettingsModel(ctx context.Context, powerPlatformAttributes map[string]attr.Value, tenantSettingsDto *tenantSettingsDto) {
	userManagementSettingsObject := powerPlatformAttributes["user_management_settings"]
	if userManagementSettingsObject != nil && !userManagementSettingsObject.IsNull() && !userManagementSettingsObject.IsUnknown() {
		objectValue, ok := userManagementSettingsObject.(basetypes.ObjectValue)
		if !ok {
			tflog.Debug(ctx, "failed to convert user_management_settings to ObjectValue", map[string]any{
				"expected_type": "basetypes.ObjectValue",
				"actual_type":   fmt.Sprintf("%T", userManagementSettingsObject),
			})
			return // Skip conversion if type assertion fails
		}

		var userManagementSettings UserManagementSettings
		objectValue.As(ctx, &userManagementSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		tenantSettingsDto.PowerPlatform.UserManagementSettings = &userManagementSettingsDto{}
		if !userManagementSettings.EnableDeleteDisabledUserinAllEnvironments.IsNull() && !userManagementSettings.EnableDeleteDisabledUserinAllEnvironments.IsUnknown() {
			tenantSettingsDto.PowerPlatform.UserManagementSettings.EnableDeleteDisabledUserinAllEnvironments = userManagementSettings.EnableDeleteDisabledUserinAllEnvironments.ValueBoolPointer()
		}
	}
}

func convertFromTenantSettingsDto[T TenantSettingsDataSourceModel | TenantSettingsResourceModel](tenantSettingsDto tenantSettingsDto, timeout timeouts.Value) (T, basetypes.ObjectValue, error) {
	objTypePowerPlatformSettings, objValuePowerPlatformSettings := convertPowerPlatformSettings(tenantSettingsDto)

	tenantSettingsProperties := map[string]attr.Type{
		"walk_me_opt_out":                                       types.BoolType,
		"disable_nps_comments_reachout":                         types.BoolType,
		"disable_newsletter_sendout":                            types.BoolType,
		"disable_environment_creation_by_non_admin_users":       types.BoolType,
		"disable_portals_creation_by_non_admin_users":           types.BoolType,
		"disable_trial_environment_creation_by_non_admin_users": types.BoolType,
		"disable_capacity_allocation_by_environment_admins":     types.BoolType,
		"disable_support_tickets_visible_by_all_users":          types.BoolType,
		"power_platform":                                        objTypePowerPlatformSettings,
	}

	tenantSettingsValues := map[string]attr.Value{
		"walk_me_opt_out":                                       types.BoolPointerValue(tenantSettingsDto.WalkMeOptOut),
		"disable_nps_comments_reachout":                         types.BoolPointerValue(tenantSettingsDto.DisableNPSCommentsReachout),
		"disable_newsletter_sendout":                            types.BoolPointerValue(tenantSettingsDto.DisableNewsletterSendout),
		"disable_environment_creation_by_non_admin_users":       types.BoolPointerValue(tenantSettingsDto.DisableEnvironmentCreationByNonAdminUsers),
		"disable_portals_creation_by_non_admin_users":           types.BoolPointerValue(tenantSettingsDto.DisablePortalsCreationByNonAdminUsers),
		"disable_trial_environment_creation_by_non_admin_users": types.BoolPointerValue(tenantSettingsDto.DisableTrialEnvironmentCreationByNonAdminUsers),
		"disable_capacity_allocation_by_environment_admins":     types.BoolPointerValue(tenantSettingsDto.DisableCapacityAllocationByEnvironmentAdmins),
		"disable_support_tickets_visible_by_all_users":          types.BoolPointerValue(tenantSettingsDto.DisableSupportTicketsVisibleByAllUsers),
		"power_platform":                                        objValuePowerPlatformSettings,
	}

	objValue := types.ObjectValueMust(tenantSettingsProperties, tenantSettingsValues)

	var result T
	switch any(result).(type) {
	case TenantSettingsDataSourceModel:
		dsModel := TenantSettingsDataSourceModel{
			Timeouts:                 timeout,
			WalkMeOptOut:             types.BoolPointerValue(tenantSettingsDto.WalkMeOptOut),
			DisableNewsletterSendout: types.BoolPointerValue(tenantSettingsDto.DisableNewsletterSendout),
			DisableEnvironmentCreationByNonAdminUsers:      types.BoolPointerValue(tenantSettingsDto.DisableEnvironmentCreationByNonAdminUsers),
			DisablePortalsCreationByNonAdminUsers:          types.BoolPointerValue(tenantSettingsDto.DisablePortalsCreationByNonAdminUsers),
			DisableTrialEnvironmentCreationByNonAdminUsers: types.BoolPointerValue(tenantSettingsDto.DisableTrialEnvironmentCreationByNonAdminUsers),
			DisableCapacityAllocationByEnvironmentAdmins:   types.BoolPointerValue(tenantSettingsDto.DisableCapacityAllocationByEnvironmentAdmins),
			DisableSupportTicketsVisibleByAllUsers:         types.BoolPointerValue(tenantSettingsDto.DisableSupportTicketsVisibleByAllUsers),
			PowerPlatform:                                  objValuePowerPlatformSettings,
		}
		typedResult, ok := any(dsModel).(T)
		if !ok {
			return result, objValue, errors.New("failed to convert TenantSettingsDataSourceModel to generic type T")
		}
		return typedResult, objValue, nil
	case TenantSettingsResourceModel:
		resModel := TenantSettingsResourceModel{
			Timeouts:                 timeout,
			Id:                       types.StringValue(""),
			WalkMeOptOut:             types.BoolPointerValue(tenantSettingsDto.WalkMeOptOut),
			DisableNewsletterSendout: types.BoolPointerValue(tenantSettingsDto.DisableNewsletterSendout),
			DisableEnvironmentCreationByNonAdminUsers:      types.BoolPointerValue(tenantSettingsDto.DisableEnvironmentCreationByNonAdminUsers),
			DisablePortalsCreationByNonAdminUsers:          types.BoolPointerValue(tenantSettingsDto.DisablePortalsCreationByNonAdminUsers),
			DisableTrialEnvironmentCreationByNonAdminUsers: types.BoolPointerValue(tenantSettingsDto.DisableTrialEnvironmentCreationByNonAdminUsers),
			DisableCapacityAllocationByEnvironmentAdmins:   types.BoolPointerValue(tenantSettingsDto.DisableCapacityAllocationByEnvironmentAdmins),
			DisableSupportTicketsVisibleByAllUsers:         types.BoolPointerValue(tenantSettingsDto.DisableSupportTicketsVisibleByAllUsers),
			PowerPlatform:                                  objValuePowerPlatformSettings,
		}
		typedResult, ok := any(resModel).(T)
		if !ok {
			return result, objValue, errors.New("failed to convert TenantSettingsResourceModel to generic type T")
		}
		return typedResult, objValue, nil
	default:
		return result, objValue, fmt.Errorf("unexpected type %T", result)
	}
}

func convertPowerPlatformSettings(tenantSettingsDto tenantSettingsDto) (basetypes.ObjectType, basetypes.ObjectValue) {
	productFeedbackObjectType, productFeedbackObjectValue := convertProductFeedbackSettings(tenantSettingsDto)
	searchSettingsObjectType, searchSettingsObjectValue := convertSearchSettings(tenantSettingsDto)
	teamsIntegrationObjectType, teamsIntegrationObjectValue := convertTeamsIntegrationSettings(tenantSettingsDto)
	powerAppsObjectType, powerAppsObjectValue := convertPowerAppsSettings(tenantSettingsDto)
	powerAutomateObjectType, powerAutomateObjectValue := convertPowerAutomateSettings(tenantSettingsDto)
	environmentsObjectType, environmentsObjectValue := convertEnvironmentSettings(tenantSettingsDto)
	governanceSettingsObjectType, governanceSettingsObjectValue := convertGovernanceSettings(tenantSettingsDto)
	licensingSettingsObjectType, licensingSettingsObjectValue := convertLicensingSettings(tenantSettingsDto)
	powerPagesSettingsObjectType, powerPagesSettingsObjectValue := convertPowerPagesSettings(tenantSettingsDto)
	championsSettingsObjectType, championsSettingsObjectValue := convertChampionsSettings(tenantSettingsDto)
	intelligenceSettingsObjectType, intelligenceSettingsObjectValue := convertIntelligenceSettings(tenantSettingsDto)
	modelExperimentationSettingsObjectType, modelExperimentationSettingsObjectValue := convertModelExperimentationSettings(tenantSettingsDto)
	catalogSettingsObjectType, catalogSettingsObjectValue := convertCatalogSettings(tenantSettingsDto)
	userManagementSettingsObjectType, userManagementSettingsObjectValue := convertUserManagementSettings(tenantSettingsDto)

	attrTypesPowerPlatformObject := map[string]attr.Type{
		"product_feedback":         productFeedbackObjectType,
		"search":                   searchSettingsObjectType,
		"teams_integration":        teamsIntegrationObjectType,
		"power_apps":               powerAppsObjectType,
		"power_automate":           powerAutomateObjectType,
		"environments":             environmentsObjectType,
		"governance":               governanceSettingsObjectType,
		"licensing":                licensingSettingsObjectType,
		"power_pages":              powerPagesSettingsObjectType,
		"champions":                championsSettingsObjectType,
		"intelligence":             intelligenceSettingsObjectType,
		"model_experimentation":    modelExperimentationSettingsObjectType,
		"catalog_settings":         catalogSettingsObjectType,
		"user_management_settings": userManagementSettingsObjectType,
	}

	if tenantSettingsDto.PowerPlatform == nil {
		return types.ObjectType{AttrTypes: attrTypesPowerPlatformObject}, types.ObjectNull(attrTypesPowerPlatformObject)
	}
	attrValuesPowerPlatformObject := map[string]attr.Value{
		"product_feedback":         productFeedbackObjectValue,
		"search":                   searchSettingsObjectValue,
		"teams_integration":        teamsIntegrationObjectValue,
		"power_apps":               powerAppsObjectValue,
		"power_automate":           powerAutomateObjectValue,
		"environments":             environmentsObjectValue,
		"governance":               governanceSettingsObjectValue,
		"licensing":                licensingSettingsObjectValue,
		"power_pages":              powerPagesSettingsObjectValue,
		"champions":                championsSettingsObjectValue,
		"intelligence":             intelligenceSettingsObjectValue,
		"model_experimentation":    modelExperimentationSettingsObjectValue,
		"catalog_settings":         catalogSettingsObjectValue,
		"user_management_settings": userManagementSettingsObjectValue,
	}
	return types.ObjectType{AttrTypes: attrTypesPowerPlatformObject}, types.ObjectValueMust(attrTypesPowerPlatformObject, attrValuesPowerPlatformObject)
}

func convertUserManagementSettings(tenantSettingsDto tenantSettingsDto) (basetypes.ObjectType, basetypes.ObjectValue) {
	attrTypes := convertBoolSettingsMap([]string{"enable_delete_disabled_user_in_all_environments"})
	if tenantSettingsDto.PowerPlatform == nil || tenantSettingsDto.PowerPlatform.UserManagementSettings == nil {
		return convertSimpleSettings(attrTypes, nil)
	}
	attrValues := map[string]attr.Value{
		"enable_delete_disabled_user_in_all_environments": types.BoolPointerValue(tenantSettingsDto.PowerPlatform.UserManagementSettings.EnableDeleteDisabledUserinAllEnvironments),
	}
	return convertSimpleSettings(attrTypes, attrValues)
}

func convertCatalogSettings(tenantSettingsDto tenantSettingsDto) (basetypes.ObjectType, basetypes.ObjectValue) {
	attrTypes := convertStringSettingsMap([]string{"power_catalog_audience_setting"})
	if tenantSettingsDto.PowerPlatform == nil || tenantSettingsDto.PowerPlatform.CatalogSettings == nil {
		return convertSimpleSettings(attrTypes, nil)
	}
	attrValues := map[string]attr.Value{
		"power_catalog_audience_setting": types.StringPointerValue(tenantSettingsDto.PowerPlatform.CatalogSettings.PowerCatalogAudienceSetting),
	}
	return convertSimpleSettings(attrTypes, attrValues)
}

func convertModelExperimentationSettings(tenantSettingsDto tenantSettingsDto) (basetypes.ObjectType, basetypes.ObjectValue) {
	attrTypes := convertBoolSettingsMap([]string{"enable_model_data_sharing", "disable_data_logging"})
	if tenantSettingsDto.PowerPlatform == nil || tenantSettingsDto.PowerPlatform.ModelExperimentation == nil {
		return convertSimpleSettings(attrTypes, nil)
	}
	attrValues := map[string]attr.Value{
		"enable_model_data_sharing": types.BoolPointerValue(tenantSettingsDto.PowerPlatform.ModelExperimentation.EnableModelDataSharing),
		"disable_data_logging":      types.BoolPointerValue(tenantSettingsDto.PowerPlatform.ModelExperimentation.DisableDataLogging),
	}
	return convertSimpleSettings(attrTypes, attrValues)
}

func convertIntelligenceSettings(tenantSettingsDto tenantSettingsDto) (basetypes.ObjectType, basetypes.ObjectValue) {
	attrTypes := convertBoolSettingsMap([]string{"disable_copilot", "allow_copilot_authors_publish_when_ai_features_are_enabled", "basic_copilot_feedback", "additional_copilot_feedback"})
	if tenantSettingsDto.PowerPlatform == nil || tenantSettingsDto.PowerPlatform.Intelligence == nil {
		return convertSimpleSettings(attrTypes, nil)
	}
	attrValues := map[string]attr.Value{
		"disable_copilot": types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Intelligence.DisableCopilot),
		"allow_copilot_authors_publish_when_ai_features_are_enabled": types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Intelligence.EnableOpenAiBotPublishing),
		"basic_copilot_feedback":                                     types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Intelligence.BasicCopilotFeedback),
		"additional_copilot_feedback":                                types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Intelligence.AdditionalCopilotFeedback),
	}
	return convertSimpleSettings(attrTypes, attrValues)
}

func convertChampionsSettings(tenantSettingsDto tenantSettingsDto) (basetypes.ObjectType, basetypes.ObjectValue) {
	attrTypes := convertBoolSettingsMap([]string{"disable_champions_invitation_reachout", "disable_skills_match_invitation_reachout"})
	if tenantSettingsDto.PowerPlatform == nil || tenantSettingsDto.PowerPlatform.Champions == nil {
		return convertSimpleSettings(attrTypes, nil)
	}
	attrValues := map[string]attr.Value{
		"disable_champions_invitation_reachout":    types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Champions.DisableChampionsInvitationReachout),
		"disable_skills_match_invitation_reachout": types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Champions.DisableSkillsMatchInvitationReachout),
	}
	return convertSimpleSettings(attrTypes, attrValues)
}

func convertPowerPagesSettings(tenantSettingsDto tenantSettingsDto) (basetypes.ObjectType, basetypes.ObjectValue) {
	attrTypesPowerPagesProperties := map[string]attr.Type{}

	if tenantSettingsDto.PowerPlatform == nil || tenantSettingsDto.PowerPlatform.PowerPages == nil {
		return types.ObjectType{AttrTypes: attrTypesPowerPagesProperties}, types.ObjectNull(attrTypesPowerPagesProperties)
	}
	attrValuesPowerPagesProperties := map[string]attr.Value{}
	return types.ObjectType{AttrTypes: attrTypesPowerPagesProperties}, types.ObjectValueMust(attrTypesPowerPagesProperties, attrValuesPowerPagesProperties)
}

func convertLicensingSettings(tenantSettingsDto tenantSettingsDto) (basetypes.ObjectType, basetypes.ObjectValue) {
	attrTypesLicencingProperties := map[string]attr.Type{
		"disable_billing_policy_creation_by_non_admin_users":           types.BoolType,
		"enable_tenant_capacity_report_for_environment_admins":         types.BoolType,
		"storage_capacity_consumption_warning_threshold":               types.Int64Type,
		"enable_tenant_licensing_report_for_environment_admins":        types.BoolType,
		"disable_use_of_unassigned_ai_builder_credits":                 types.BoolType,
		"apply_auto_claim_power_apps_to_only_managed_environments":     types.BoolType,
		"apply_auto_claim_power_automate_to_only_managed_environments": types.BoolType,
	}

	if tenantSettingsDto.PowerPlatform == nil || tenantSettingsDto.PowerPlatform.Licensing == nil {
		return types.ObjectType{AttrTypes: attrTypesLicencingProperties}, types.ObjectNull(attrTypesLicencingProperties)
	}
	attrValuesLicencingProperties := map[string]attr.Value{
		"disable_billing_policy_creation_by_non_admin_users":           types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Licensing.DisableBillingPolicyCreationByNonAdminUsers),
		"enable_tenant_capacity_report_for_environment_admins":         types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Licensing.EnableTenantCapacityReportForEnvironmentAdmins),
		"storage_capacity_consumption_warning_threshold":               types.Int64PointerValue(tenantSettingsDto.PowerPlatform.Licensing.StorageCapacityConsumptionWarningThreshold),
		"enable_tenant_licensing_report_for_environment_admins":        types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Licensing.EnableTenantLicensingReportForEnvironmentAdmins),
		"disable_use_of_unassigned_ai_builder_credits":                 types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Licensing.DisableUseOfUnassignedAIBuilderCredits),
		"apply_auto_claim_power_apps_to_only_managed_environments":     types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Licensing.ApplyAutoClaimToOnlyManagedEnvironments),
		"apply_auto_claim_power_automate_to_only_managed_environments": types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Licensing.ApplyPAutoAutoClaimToOnlyManagedEnvironments),
	}
	return types.ObjectType{AttrTypes: attrTypesLicencingProperties}, types.ObjectValueMust(attrTypesLicencingProperties, attrValuesLicencingProperties)
}

func convertGovernanceSettings(tenantSettingsDto tenantSettingsDto) (basetypes.ObjectType, basetypes.ObjectValue) {
	attrTypesPolicyProperties := map[string]attr.Type{
		"enable_desktop_flow_data_policy_management": types.BoolType,
	}

	attrTypesGovernanceProperties := map[string]attr.Type{
		"weekly_digest_email_recipients":                            types.SetType{ElemType: types.StringType},
		"disable_admin_digest":                                      types.BoolType,
		"disable_developer_environment_creation_by_non_admin_users": types.BoolType,
		"enable_default_environment_routing":                        types.BoolType,
		"environment_routing_all_makers":                            types.BoolType,
		"environment_routing_target_environment_group_id":           customtypes.UUIDType{},
		"environment_routing_target_security_group_id":              customtypes.UUIDType{},
		"policy": types.ObjectType{AttrTypes: attrTypesPolicyProperties},
	}

	if tenantSettingsDto.PowerPlatform == nil || tenantSettingsDto.PowerPlatform.Governance == nil {
		return types.ObjectType{AttrTypes: attrTypesGovernanceProperties}, types.ObjectNull(attrTypesGovernanceProperties)
	}
	var objValuePolicyProperties basetypes.ObjectValue
	if tenantSettingsDto.PowerPlatform.Governance.Policy == nil {
		objValuePolicyProperties = types.ObjectNull(attrTypesPolicyProperties)
	} else {
		objValuePolicyProperties = types.ObjectValueMust(attrTypesPolicyProperties, map[string]attr.Value{
			"enable_desktop_flow_data_policy_management": types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Governance.Policy.EnableDesktopFlowDataPolicyManagement),
		})
	}

	weeklyDigestList := []attr.Value{}
	if tenantSettingsDto.PowerPlatform.Governance.AdditionalAdminDigestEmailRecipients != nil && *tenantSettingsDto.PowerPlatform.Governance.AdditionalAdminDigestEmailRecipients != "" {
		for _, email := range strings.Split(*tenantSettingsDto.PowerPlatform.Governance.AdditionalAdminDigestEmailRecipients, ";") {
			weeklyDigestList = append(weeklyDigestList, types.StringValue(email))
		}
	}

	var weeklyDigestValue attr.Value
	if tenantSettingsDto.PowerPlatform.Governance.AdditionalAdminDigestEmailRecipients == nil {
		// Field not present in API response - return null
		weeklyDigestValue = types.SetNull(types.StringType)
	} else {
		// Field present (even if empty string) - return set (possibly empty)
		weeklyDigestValue = types.SetValueMust(types.StringType, weeklyDigestList)
	}

	attrValuesGovernanceProperties := map[string]attr.Value{
		"weekly_digest_email_recipients":                            weeklyDigestValue,
		"disable_admin_digest":                                      types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Governance.DisableAdminDigest),
		"disable_developer_environment_creation_by_non_admin_users": types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Governance.DisableDeveloperEnvironmentCreationByNonAdminUsers),
		"enable_default_environment_routing":                        types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Governance.EnableDefaultEnvironmentRouting),
		"environment_routing_all_makers":                            types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Governance.EnvironmentRoutingAllMakers),
		"environment_routing_target_environment_group_id":           customtypes.NewUUIDPointerValue(tenantSettingsDto.PowerPlatform.Governance.EnvironmentRoutingTargetEnvironmentGroupId),
		"environment_routing_target_security_group_id":              customtypes.NewUUIDPointerValue(tenantSettingsDto.PowerPlatform.Governance.EnvironmentRoutingTargetSecurityGroupId),
		"policy": objValuePolicyProperties,
	}
	return types.ObjectType{AttrTypes: attrTypesGovernanceProperties}, types.ObjectValueMust(attrTypesGovernanceProperties, attrValuesGovernanceProperties)
}

func convertEnvironmentSettings(tenantSettingsDto tenantSettingsDto) (basetypes.ObjectType, basetypes.ObjectValue) {
	attrTypesEnvironmentsProperties := map[string]attr.Type{
		"disable_preferred_data_location_for_teams_environment": types.BoolType,
	}

	if tenantSettingsDto.PowerPlatform == nil || tenantSettingsDto.PowerPlatform.Environments == nil {
		return types.ObjectType{AttrTypes: attrTypesEnvironmentsProperties}, types.ObjectNull(attrTypesEnvironmentsProperties)
	}
	attrValuesEnvironmentsProperties := map[string]attr.Value{
		"disable_preferred_data_location_for_teams_environment": types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Environments.DisablePreferredDataLocationForTeamsEnvironment),
	}
	return types.ObjectType{AttrTypes: attrTypesEnvironmentsProperties}, types.ObjectValueMust(attrTypesEnvironmentsProperties, attrValuesEnvironmentsProperties)
}

func convertPowerAutomateSettings(tenantSettingsDto tenantSettingsDto) (basetypes.ObjectType, basetypes.ObjectValue) {
	attrTypesPowerAutomateProperties := map[string]attr.Type{
		"disable_copilot":             types.BoolType,
		"disable_copilot_with_bing":   types.BoolType,
		"allow_use_of_hosted_browser": types.BoolType,
		"disable_flow_resubmission":   types.BoolType,
	}

	if tenantSettingsDto.PowerPlatform == nil || tenantSettingsDto.PowerPlatform.PowerAutomate == nil {
		return types.ObjectType{AttrTypes: attrTypesPowerAutomateProperties}, types.ObjectNull(attrTypesPowerAutomateProperties)
	}
	attrValuesPowerAutomateProperties := map[string]attr.Value{
		"disable_copilot":             types.BoolPointerValue(tenantSettingsDto.PowerPlatform.PowerAutomate.DisableCopilot),
		"disable_copilot_with_bing":   types.BoolPointerValue(tenantSettingsDto.PowerPlatform.PowerAutomate.DisableCopilotWithBing),
		"allow_use_of_hosted_browser": types.BoolPointerValue(tenantSettingsDto.PowerPlatform.PowerAutomate.AllowUseOfHostedBrowser),
		"disable_flow_resubmission":   types.BoolPointerValue(tenantSettingsDto.PowerPlatform.PowerAutomate.DisableFlowResubmission),
	}
	return types.ObjectType{AttrTypes: attrTypesPowerAutomateProperties}, types.ObjectValueMust(attrTypesPowerAutomateProperties, attrValuesPowerAutomateProperties)
}

func convertPowerAppsSettings(tenantSettingsDto tenantSettingsDto) (basetypes.ObjectType, basetypes.ObjectValue) {
	attrTypesPowerAppsProperties := map[string]attr.Type{
		"disable_copilot":                          types.BoolType,
		"disable_share_with_everyone":              types.BoolType,
		"enable_guests_to_make":                    types.BoolType,
		"disable_maker_match":                      types.BoolType,
		"disable_unused_license_assignment":        types.BoolType,
		"disable_create_from_image":                types.BoolType,
		"disable_create_from_figma":                types.BoolType,
		"disable_connection_sharing_with_everyone": types.BoolType,
		"enable_canvas_app_insights":               types.BoolType,
	}

	if tenantSettingsDto.PowerPlatform == nil || tenantSettingsDto.PowerPlatform.PowerApps == nil {
		return types.ObjectType{AttrTypes: attrTypesPowerAppsProperties}, types.ObjectNull(attrTypesPowerAppsProperties)
	}
	attrValuesPowerAppsProperties := map[string]attr.Value{
		"disable_copilot":                          types.BoolPointerValue(tenantSettingsDto.PowerPlatform.PowerApps.DisableCopilot),
		"disable_share_with_everyone":              types.BoolPointerValue(tenantSettingsDto.PowerPlatform.PowerApps.DisableShareWithEveryone),
		"enable_guests_to_make":                    types.BoolPointerValue(tenantSettingsDto.PowerPlatform.PowerApps.EnableGuestsToMake),
		"disable_maker_match":                      types.BoolPointerValue(tenantSettingsDto.PowerPlatform.PowerApps.DisableMakerMatch),
		"disable_unused_license_assignment":        types.BoolPointerValue(tenantSettingsDto.PowerPlatform.PowerApps.DisableUnusedLicenseAssignment),
		"disable_create_from_image":                types.BoolPointerValue(tenantSettingsDto.PowerPlatform.PowerApps.DisableCreateFromImage),
		"disable_create_from_figma":                types.BoolPointerValue(tenantSettingsDto.PowerPlatform.PowerApps.DisableCreateFromFigma),
		"disable_connection_sharing_with_everyone": types.BoolPointerValue(tenantSettingsDto.PowerPlatform.PowerApps.DisableConnectionSharingWithEveryone),
		"enable_canvas_app_insights":               types.BoolPointerValue(tenantSettingsDto.PowerPlatform.PowerApps.EnableCanvasAppInsights),
	}
	return types.ObjectType{AttrTypes: attrTypesPowerAppsProperties}, types.ObjectValueMust(attrTypesPowerAppsProperties, attrValuesPowerAppsProperties)
}

func convertTeamsIntegrationSettings(tenantSettingsDto tenantSettingsDto) (basetypes.ObjectType, basetypes.ObjectValue) {
	attrTypesTeamsIntegrationProperties := map[string]attr.Type{
		"share_with_colleagues_user_limit": types.Int64Type,
	}
	if tenantSettingsDto.PowerPlatform == nil || tenantSettingsDto.PowerPlatform.TeamsIntegration == nil {
		return types.ObjectType{AttrTypes: attrTypesTeamsIntegrationProperties}, types.ObjectNull(attrTypesTeamsIntegrationProperties)
	}
	attrValuesTeamsIntegrationProperties := map[string]attr.Value{
		"share_with_colleagues_user_limit": types.Int64PointerValue(tenantSettingsDto.PowerPlatform.TeamsIntegration.ShareWithColleaguesUserLimit),
	}
	return types.ObjectType{AttrTypes: attrTypesTeamsIntegrationProperties}, types.ObjectValueMust(attrTypesTeamsIntegrationProperties, attrValuesTeamsIntegrationProperties)
}

func convertProductFeedbackSettings(tenantSettingsDto tenantSettingsDto) (basetypes.ObjectType, basetypes.ObjectValue) {
	attrTypesProductFeedbackProperties := map[string]attr.Type{
		"disable_user_survey_feedback":   types.BoolType,
		"disable_microsoft_follow_up":    types.BoolType,
		"disable_attachments":            types.BoolType,
		"disable_microsoft_surveys_send": types.BoolType,
	}

	attrValuesProductFeedbackProperties := map[string]attr.Value{
		"disable_user_survey_feedback":   types.BoolPointerValue(tenantSettingsDto.DisableSurveyFeedback),
		"disable_microsoft_follow_up":    types.BoolPointerValue(tenantSettingsDto.DisableNPSCommentsReachout),
		"disable_attachments":            types.BoolPointerValue(tenantSettingsDto.DisableSurveyScreenshots),
		"disable_microsoft_surveys_send": types.BoolPointerValue(tenantSettingsDto.DisableUserInitiatedFeedback),
	}
	return types.ObjectType{AttrTypes: attrTypesProductFeedbackProperties}, types.ObjectValueMust(attrTypesProductFeedbackProperties, attrValuesProductFeedbackProperties)
}

func convertSearchSettings(tenantSettingsDto tenantSettingsDto) (basetypes.ObjectType, basetypes.ObjectValue) {
	attrTypesSearchProperties := map[string]attr.Type{
		"disable_docs_search":       types.BoolType,
		"disable_community_search":  types.BoolType,
		"disable_bing_video_search": types.BoolType,
	}

	if tenantSettingsDto.PowerPlatform == nil || tenantSettingsDto.PowerPlatform.Search == nil {
		return types.ObjectType{AttrTypes: attrTypesSearchProperties}, types.ObjectNull(attrTypesSearchProperties)
	}
	attrValuesSearchProperties := map[string]attr.Value{
		"disable_docs_search":       types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Search.DisableDocsSearch),
		"disable_community_search":  types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Search.DisableCommunitySearch),
		"disable_bing_video_search": types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Search.DisableBingVideoSearch),
	}

	return types.ObjectType{AttrTypes: attrTypesSearchProperties}, types.ObjectValueMust(attrTypesSearchProperties, attrValuesSearchProperties)
}
