// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package environment_settings

type environmentSettingsValueDto struct {
	Value []environmentSettingsDto `json:"value"`
}

type environmentSettingsDto struct {
	MaxUploadFileSize                        *int64  `json:"maxuploadfilesize,omitempty"`
	PluginTraceLogSetting                    *int64  `json:"plugintracelogsetting,omitempty"`
	IsAuditEnabled                           *bool   `json:"isauditenabled,omitempty"`
	IsUserAccessAuditEnabled                 *bool   `json:"isuseraccessauditenabled,omitempty"`
	IsReadAuditEnabled                       *bool   `json:"isreadauditenabled,omitempty"`
	AuditRetentionPeriodV2                   *int32  `json:"auditretentionperiodv2,omitempty"`
	BoundDashboardDefaultCardExpanded        *bool   `json:"bounddashboarddefaultcardexpanded,omitempty"`
	OrganizationId                           *string `json:"organizationid,omitempty"`
	PowerAppsComponentFrameworkForCanvasApps *bool   `json:"iscustomcontrolsincanvasappsenabled,omitempty"`

	EnableIpBasedFirewallRule            *bool   `json:"enableipbasedfirewallrule,omitempty"`     //enable disable whole thing
	AllowedIpRangeForFirewall            *string `json:"allowediprangeforfirewall,omitempty"`     //ip ranges
	AllowedServiceTagsForFirewall        *string `json:"allowedservicetagsforfirewall,omitempty"` //service tags
	AllowApplicationUserAccess           *bool   `json:"allowapplicationuseraccess,omitempty"`
	AllowMicrosoftTrustedServiceTags     *bool   `json:"allowmicrosofttrustedservicetags,omitempty"`
	EnableIpBasedFirewallRuleInAuditMode *bool   `json:"enableipbasedfirewallruleinauditmode,omitempty"`
	ReverseProxyIpAddresses              *string `json:"reverseproxyipaddresses,omitempty"`
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
