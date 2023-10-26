package powerplatform

type TenantSettingsDto struct {
	WalkMeOptOut                                   bool `json:"walkMeOptOut"`
	DisableNPSCommentsReachout                     bool `json:"disableNPSCommentsReachout"`
	DisableNewsletterSendout                       bool `json:"disableNewsletterSendout"`
	DisableEnvironmentCreationByNonAdminUsers      bool `json:"disableEnvironmentCreationByNonAdminUsers"`
	DisablePortalsCreationByNonAdminUsers          bool `json:"disablePortalsCreationByNonAdminUsers"`
	DisableSurveyFeedback                          bool `json:"disableSurveyFeedback"`
	DisableTrialEnvironmentCreationByNonAdminUsers bool `json:"disableTrialEnvironmentCreationByNonAdminUsers"`
	DisableCapacityAllocationByEnvironmentAdmins   bool `json:"disableCapacityAllocationByEnvironmentAdmins"`
	DisableSupportTicketsVisibleByAllUsers         bool `json:"disableSupportTicketsVisibleByAllUsers"`
	PowerPlatform                                  struct {
		Search struct {
			DisableDocsSearch      bool `json:"disableDocsSearch"`
			DisableCommunitySearch bool `json:"disableCommunitySearch"`
			DisableBingVideoSearch bool `json:"disableBingVideoSearch"`
		} `json:"search"`
		TeamsIntegration struct {
			ShareWithColleaguesUserLimit int64 `json:"shareWithColleaguesUserLimit"`
		} `json:"teamsIntegration"`
		PowerApps struct {
			DisableShareWithEveryone             bool `json:"disableShareWithEveryone"`
			EnableGuestsToMake                   bool `json:"enableGuestsToMake"`
			DisableMembersIndicator              bool `json:"disableMembersIndicator"`
			DisableMakerMatch                    bool `json:"disableMakerMatch"`
			DisableUnusedLicenseAssignment       bool `json:"disableUnusedLicenseAssignment"`
			DisableCreateFromImage               bool `json:"disableCreateFromImage"`
			DisableCreateFromFigma               bool `json:"disableCreateFromFigma"`
			DisableConnectionSharingWithEveryone bool `json:"disableConnectionSharingWithEveryone"`
		} `json:"powerApps"`
		PowerAutomate struct {
			DisableCopilot bool `json:"disableCopilot"`
		} `json:"powerAutomate"`
		Environments struct {
			DisablePreferredDataLocationForTeamsEnvironment bool `json:"disablePreferredDataLocationForTeamsEnvironment"`
		} `json:"environments"`
		Governance struct {
			DisableAdminDigest                                 bool `json:"disableAdminDigest"`
			DisableDeveloperEnvironmentCreationByNonAdminUsers bool `json:"disableDeveloperEnvironmentCreationByNonAdminUsers"`
			EnableDefaultEnvironmentRouting                    bool `json:"enableDefaultEnvironmentRouting"`
			Policy                                             struct {
				EnableDesktopFlowDataPolicyManagement bool `json:"enableDesktopFlowDataPolicyManagement"`
			} `json:"policy"`
		} `json:"governance"`
		Licensing struct {
			DisableBillingPolicyCreationByNonAdminUsers     bool  `json:"disableBillingPolicyCreationByNonAdminUsers"`
			EnableTenantCapacityReportForEnvironmentAdmins  bool  `json:"enableTenantCapacityReportForEnvironmentAdmins"`
			StorageCapacityConsumptionWarningThreshold      int64 `json:"storageCapacityConsumptionWarningThreshold"`
			EnableTenantLicensingReportForEnvironmentAdmins bool  `json:"enableTenantLicensingReportForEnvironmentAdmins"`
			DisableUseOfUnassignedAIBuilderCredits          bool  `json:"disableUseOfUnassignedAIBuilderCredits"`
		} `json:"licensing"`
		PowerPages struct {
		} `json:"powerPages"`
		Champions struct {
			DisableChampionsInvitationReachout   bool `json:"disableChampionsInvitationReachout"`
			DisableSkillsMatchInvitationReachout bool `json:"disableSkillsMatchInvitationReachout"`
		} `json:"champions"`
		Intelligence struct {
			DisableCopilot            bool `json:"disableCopilot"`
			EnableOpenAiBotPublishing bool `json:"enableOpenAiBotPublishing"`
		} `json:"intelligence"`
		ModelExperimentation struct {
			EnableModelDataSharing bool `json:"enableModelDataSharing"`
			DisableDataLogging     bool `json:"disableDataLogging"`
		} `json:"modelExperimentation"`
		CatalogSettings struct {
			PowerCatalogAudienceSetting string `json:"powerCatalogAudienceSetting"`
		} `json:"catalogSettings"`
	} `json:"powerPlatform"`
}
