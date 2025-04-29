# Title

Unexported Struct Definitions Lead to Limited Accessibility

# Path

`/workspaces/terraform-provider-power-platform/internal/services/environment_settings/dto.go`

# Problem

All the struct definitions (`environmentSettingsValueDto, environmentSettingsDto, environmentIdDto, environmentIdPropertiesDto, linkedEnvironmentIdMetadataDto`) are unexported (start with lowercase), which makes them inaccessible outside the `environment_settings` package. This could become an issue if these DTOs are meant to be marshaled or unmarshaled in other parts of the application.

# Impact

1. Data Transfer Objects (DTOs) cannot be effectively reused or extended in other packages, limiting functionality.
2. Might lead to redundant struct definitions in other parts of the application.

**Severity**: **High**

# Location

```go
type environmentSettingsValueDto struct {
  Value []environmentSettingsDto `json:"value"`
}

type environmentSettingsDto struct {
  ...
}

type environmentIdDto struct {
  Id         string                     `json:"id"`
  Name       string                     `json:"name"`
  Properties environmentIdPropertiesDto `json:"properties"`
}

type environmentIdPropertiesDto struct {
  ...
}

type linkedEnvironmentIdMetadataDto struct {
  InstanceURL string
}
```

# Fix

To fix this, capitalize the struct names to export them.

```go
type EnvironmentSettingsValueDto struct {
  Value []EnvironmentSettingsDto `json:"value"`
}

type EnvironmentSettingsDto struct {
  MaxUploadFileSize                        *int64  `json:"maxuploadfilesize,omitempty"`
  PluginTraceLogSetting                    *int64  `json:"plugintracelogsetting,omitempty"`
  IsAuditEnabled                           *bool   `json:"isauditenabled,omitempty"`
  IsUserAccessAuditEnabled                 *bool   `json:"isuseraccessauditenabled,omitempty"`
  IsReadAuditEnabled                       *bool   `json:"isreadauditenabled,omitempty"`
  AuditRetentionPeriodV2                   *int32  `json:"auditretentionperiodv2,omitempty"`
  BoundDashboardDefaultCardExpanded        *bool   `json:"bounddashboarddefaultcardexpanded,omitempty"`
  OrganizationId                           *string `json:"organizationid,omitempty"`
  PowerAppsComponentFrameworkForCanvasApps *bool   `json:"iscustomcontrolsincanvasappsenabled,omitempty"`
  EnableIpBasedCookieBinding               *bool   `json:"enableipbasedcookiebinding,omitempty"`
  EnableIpBasedFirewallRule                *bool   `json:"enableipbasedfirewallrule,omitempty"`
  AllowedIpRangeForFirewall                *string `json:"allowediprangeforfirewall,omitempty"`
  AllowedServiceTagsForFirewall            *string `json:"allowedservicetagsforfirewall,omitempty"`
  AllowApplicationUserAccess               *bool   `json:"allowapplicationuseraccess,omitempty"`
  AllowMicrosoftTrustedServiceTags         *bool   `json:"allowmicrosofttrustedservicetags,omitempty"`
  EnableIpBasedFirewallRuleInAuditMode     *bool   `json:"enableipbasedfirewallruleinauditmode,omitempty"`
  ReverseProxyIpAddresses                  *string `json:"reverseproxyipaddresses,omitempty"`
}

type EnvironmentIdDto struct {
  Id         string                      `json:"id"`
  Name       string                      `json:"name"`
  Properties EnvironmentIdPropertiesDto `json:"properties"`
}

type EnvironmentIdPropertiesDto struct {
  LinkedEnvironmentMetadata LinkedEnvironmentIdMetadataDto `json:"linkedEnvironmentMetadata"`
}

type LinkedEnvironmentIdMetadataDto struct {
  InstanceURL string
}
```

This fix ensures that the structs are accessible in other packages while still retaining their functionality.

---