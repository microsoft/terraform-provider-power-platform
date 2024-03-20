// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

type EnvironmentSettingsValueDto struct {
	Value []EnvironmentSettingsDto `json:"value"`
}

type EnvironmentSettingsDto struct {
	MaxUploadFileSize                 int64 `json:"maxuploadfilesize"`
	PluginTraceLogSetting             int64 `json:"plugintracelogsetting"`
	IsAuditEnabled                    bool  `json:"IsAuditEnabled"`
	IsUserAccessAuditEnabled          bool  `json:"isuseraccessauditenabled"`
	IsReadAuditEnabled                bool  `json:"isreadauditenabled"`
	BoundDashboardDefaultCardExpanded bool  `json:"bounddashboarddefaultcardexpanded"`
}

type EnvironmentIdDto struct {
	Id         string                     `json:"id"`
	Name       string                     `json:"name"`
	Properties EnvironmentIdPropertiesDto `json:"properties"`
}

type EnvironmentIdPropertiesDto struct {
	LinkedEnvironmentMetadata LinkedEnvironmentIdMetadataDto `json:"linkedEnvironmentMetadata"`
}

type LinkedEnvironmentIdMetadataDto struct {
	InstanceURL string
}
