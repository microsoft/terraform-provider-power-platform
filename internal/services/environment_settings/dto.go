// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_settings

type environmentSettings struct {
	BackendSettings *environmentBackendSettingDto
	OrgSettings     *environmentOrgSettingsDto
}

type environmentBackendSettingsValueDto struct {
	ODataContext            string                         `json:"@odata.context"`
	SettingDetailCollection []environmentBackendSettingDto `json:"SettingDetailCollection"`
}

type environmentBackendSettingDto struct {
	Name     string `json:"Name"`
	Value    string `json:"Value"`
	DataType int    `json:"DataType"`
}

type environmentOrgSettingsValueDto struct {
	Value []environmentOrgSettingsDto `json:"value"`
}

type environmentOrgSettingsDto struct {
	MaxUploadFileSize                                    *int64  `json:"maxuploadfilesize,omitempty"`
	PluginTraceLogSetting                                *int64  `json:"plugintracelogsetting,omitempty"`
	IsAuditEnabled                                       *bool   `json:"isauditenabled,omitempty"`
	IsUserAccessAuditEnabled                             *bool   `json:"isuseraccessauditenabled,omitempty"`
	IsReadAuditEnabled                                   *bool   `json:"isreadauditenabled,omitempty"`
	AuditRetentionPeriodV2                               *int32  `json:"auditretentionperiodv2,omitempty"`
	BoundDashboardDefaultCardExpanded                    *bool   `json:"bounddashboarddefaultcardexpanded,omitempty"`
	OrganizationId                                       *string `json:"organizationid,omitempty"`
	PowerAppsComponentFrameworkForCanvasApps             *bool   `json:"iscustomcontrolsincanvasappsenabled,omitempty"`
	PowerAppsMakerBotEnabled                             *bool   `json:"powerappsmakerbotenabled,omitempty"`
	BlockAccessToSessionTranscriptsForCopilotStudio      *bool   `json:"blockaccesstosessiontranscriptsforcopilotstudio,omitempty"`
	BlockTranscriptRecordingForCopilotStudio             *bool   `json:"blocktranscriptrecordingforcopilotstudio,omitempty"`
	EnableCopilotStudioShareDataWithVivaInsights         *bool   `json:"enablecopilotstudiosharedatawithvivainsights,omitempty"`
	EnableCopilotStudioCrossGeoShareDataWithVivaInsights *bool   `json:"enablecopilotstudiocrossgeosharedatawithvivainsights,omitempty"`
	PaiPreviewScenarioEnabled                            *bool   `json:"paipreviewscenarioenabled,omitempty"`
	AiPromptsEnabled                                     *bool   `json:"aipromptsenabled,omitempty"`
	EnableIpBasedCookieBinding                           *bool   `json:"enableipbasedcookiebinding,omitempty"`
	EnableIpBasedFirewallRule                            *bool   `json:"enableipbasedfirewallrule,omitempty"`
	AllowedIpRangeForFirewall                            *string `json:"allowediprangeforfirewall,omitempty"`
	AllowedServiceTagsForFirewall                        *string `json:"allowedservicetagsforfirewall,omitempty"`
	AllowApplicationUserAccess                           *bool   `json:"allowapplicationuseraccess,omitempty"`
	AllowMicrosoftTrustedServiceTags                     *bool   `json:"allowmicrosofttrustedservicetags,omitempty"`
	EnableIpBasedFirewallRuleInAuditMode                 *bool   `json:"enableipbasedfirewallruleinauditmode,omitempty"`
	ReverseProxyIpAddresses                              *string `json:"reverseproxyipaddresses,omitempty"`
}

type environmentIdDto struct {
	Id         string                     `json:"id"`
	Name       string                     `json:"name"`
	Properties environmentIdPropertiesDto `json:"properties"`
}

type environmentIdPropertiesDto struct {
	LinkedEnvironmentMetadata linkedEnvironmentIdMetadataDto `json:"linkedEnvironmentMetadata"`
}

type linkedEnvironmentIdMetadataDto struct {
	InstanceURL string
}
