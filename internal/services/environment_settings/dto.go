// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_settings

import (
	"fmt"
)

const (
	ON                                = "On"
	OFF                               = "Off"
	DEFAULT                           = "Default"
	AUTO                              = "Auto"
	ALL_USERS                         = "AllUsers"
	USER_AS_FEATURE_BECOMES_AVAILABLE = "UserAsFeatureBecomesAvailable"
	NO_ONE                            = "NoOne"
)

const (
	EnableCopilotAnswerControl              = "EnableCopilotAnswerControl"
	AppCopilotEnabled                       = "appcopilotenabled"
	FormPredictEnabled                      = "FormPredictEnabled"
	FormPredictSmartPasteEnabledOnByDefault = "FormPredictSmartPasteEnabledOnByDefault"
	FormFillBarUXEnabled                    = "FormFillBarUXEnabled"
	NLGridSearchSetting                     = "NLGridSearchSetting"
	NLChartDataVisualizationSetting         = "NLChartDataVisualizationSetting"
)

type environmentSettings struct {
	BackendSettings *environmentBackendSettingsValueDto
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

func (environmentBackendSettings environmentBackendSettingsValueDto) GetValue(name string) string {
	for _, setting := range environmentBackendSettings.SettingDetailCollection {
		if setting.Name == name {
			return setting.Value
		}
	}
	return ""
}

func (environmentBackendSettings *environmentBackendSettingsValueDto) SetValue(name string, value any) error {
	for i, setting := range environmentBackendSettings.SettingDetailCollection {
		if setting.Name == name {
			switch v := value.(type) {
			case string:
				environmentBackendSettings.SettingDetailCollection[i].Value = v
			case bool:
				if v {
					environmentBackendSettings.SettingDetailCollection[i].Value = "true"
				} else {
					environmentBackendSettings.SettingDetailCollection[i].Value = "false"
				}
			default:
				return fmt.Errorf("unsupported value type: %T", v)
			}
			return nil
		}
	}
	return fmt.Errorf("setting not found: %s", name)
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
