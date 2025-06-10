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
