# Title

Inconsistent naming: ProviderConfigUrls vs ProviderConfigUrls fields

##

/workspaces/terraform-provider-power-platform/internal/config/config.go

## Problem

The struct `ProviderConfigUrls` and its fields use inconsistent naming conventions: some fields end with `Url`, some with `Scope`, and some with other variations (e.g., `LicensingUrl`, `PowerAppsAdvisor`, `PowerAppsAdvisorScope`). This inconsistency can confuse maintainers and lead to errors in usage or further development.

## Impact

Severity: Low

This is a low severity issue, as it does not affect correctness, but it does reduce code readability, maintainability, and consistency in the API exposed to consumers.

## Location

Struct definition and field names:

```go
type ProviderConfigUrls struct {
    AdminPowerPlatformUrl string
    BapiUrl               string
    PowerAppsUrl          string
    PowerAppsScope        string
    PowerPlatformUrl      string
    PowerPlatformScope    string
    LicensingUrl          string
    PowerAppsAdvisor      string
    PowerAppsAdvisorScope string
    AnalyticsScope        string
}
```

## Code Issue

```go
type ProviderConfigUrls struct {
    AdminPowerPlatformUrl string
    BapiUrl               string
    PowerAppsUrl          string
    PowerAppsScope        string
    PowerPlatformUrl      string
    PowerPlatformScope    string
    LicensingUrl          string
    PowerAppsAdvisor      string
    PowerAppsAdvisorScope string
    AnalyticsScope        string
}
```

## Fix

Align naming conventions. For example, append `Url` to all URLs and use `Scope` for scopes. If `PowerAppsAdvisor` is a URL, rename to `PowerAppsAdvisorUrl`:

```go
type ProviderConfigUrls struct {
    AdminPowerPlatformUrl      string
    BapiUrl                    string
    PowerAppsUrl               string
    PowerAppsScope             string
    PowerPlatformUrl           string
    PowerPlatformScope         string
    LicensingUrl               string
    PowerAppsAdvisorUrl        string
    PowerAppsAdvisorScope      string
    AnalyticsScope             string
}
```

This improves clarity and reduces confusion for future maintainers.
