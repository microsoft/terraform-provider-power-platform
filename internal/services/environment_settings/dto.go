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
	BoundDashboardDefaultCardExpanded        *bool   `json:"bounddashboarddefaultcardexpanded,omitempty"`
	OrganizationId                           *string `json:"organizationid,omitempty"`
	PowerAppsComponentFrameworkForCanvasApps *bool   `json:"iscustomcontrolsincanvasappsenabled,omitempty"`
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
