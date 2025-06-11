# Config Constants Issues - Merged Issues

## ISSUE 1

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


---

## ISSUE 2

# Title 

Function Naming - StringPtr

##

/workspaces/terraform-provider-power-platform/internal/helpers/config.go

## Problem

The `StringPtr` naming does not follow Go idiomatic naming, which typically suggests `StringPointer`.

## Impact

Low severity. Minor impact, mostly readability and consistency.

## Location

Line 81-83

## Code Issue

```go
func StringPtr(s string) *string {
	return &s
}
```

## Fix

Rename the function to `StringPointer` for consistency with Go naming conventions for such helpers.

```go
func StringPointer(s string) *string {
	return &s
}
```


---

## ISSUE 3

# Title

Inconsistent Naming Conventions for Constant Identifiers

##
/workspaces/terraform-provider-power-platform/internal/constants/constants.go

## Problem

Some constant identifiers use a naming pattern that is inconsistent or potentially confusing. For example, `DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES` is misleading because the value is actually in `time.Duration` (nanoseconds), not strictly an integer count of minutes. Additionally, a few constants use redundant or unclear suffixes, such as `*_URL` (where some are hostnames, not URLs), or mix capitalization patterns for similar concepts.

## Impact

Low severity. The code will still run as expected, but inconsistent or misleading naming can decrease readability and maintainability, and may confuse developers about the expected units or usage of a constant.

## Location

Example:
```go
const (
	DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES = 20 * time.Minute
)
```

Other cases:
- `*_URL` sometimes points to hosts, sometimes to full URLs

## Code Issue

```go
const (
	DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES = 20 * time.Minute
)
```
and elsewhere similar `*_URL` or `*_DOMAIN` distinctions.

## Fix

Adopt clear, consistent naming to communicate intent and type/units:

- For the `DEFAULT_RESOURCE_OPERATION_TIMEOUT_IN_MINUTES`, use `DEFAULT_RESOURCE_OPERATION_TIMEOUT` (no unit in name if storing as a duration).
- For hosts, use `*_HOST` or `*_DOMAIN`. Use `*_URL` for complete URLs only.
- Use consistent suffixes and casing.

Example:

```go
const (
	DEFAULT_RESOURCE_OPERATION_TIMEOUT = 20 * time.Minute
	PUBLIC_ADMIN_POWER_PLATFORM_HOST   = "api.admin.powerplatform.microsoft.com"
	PUBLIC_OAUTH_AUTHORITY_URL        = "https://login.microsoftonline.com/"
)
```

---

**Save location:**  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/constants.go-naming-low.md


---

# To finish the task you have to 
1. Run linter and fix any issues 
2. Run UnitTest and fix any of failing ones
3. Generate docs 
4. Run Changie

# Changie Instructions
Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```
Where:
- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
