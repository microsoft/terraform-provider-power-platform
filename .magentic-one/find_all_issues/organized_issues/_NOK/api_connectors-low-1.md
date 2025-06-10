# Magic Strings for API Version and Filters

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/api_connectors.go

## Problem

The code uses hardcoded, "magic" strings for query parameters such as `"2019-05-01"` (API version), `"showApisWithToS"`, and the filter `"environment eq '~Default'"`. These values are not documented or given semantic meaning, and are repeated in one place without any abstraction.

## Impact

Using magic strings directly in the code reduces readability, increases the risk of typos, and makes upgrades or refactoring harder (e.g., upgrading API versions in the future). Severity: **low**.

## Location

Within the `GetConnectors` method:

```go
values.Add("api-version", "2019-05-01")
values.Add("showApisWithToS", "true")
values.Add("hideDlpExemptApis", "true")
values.Add("showAllDlpEnforceableApis", "true")
values.Add("$filter", "environment eq '~Default'")
```

## Code Issue

```go
values.Add("api-version", "2019-05-01")
values.Add("showApisWithToS", "true")
values.Add("hideDlpExemptApis", "true")
values.Add("showAllDlpEnforceableApis", "true")
values.Add("$filter", "environment eq '~Default'")
```

## Fix

Define constants at the top of the file for these values:

```go
const (
	apiVersion                  = "2019-05-01"
	showApisWithToS             = "showApisWithToS"
	hideDlpExemptApis           = "hideDlpExemptApis"
	showAllDlpEnforceableApis   = "showAllDlpEnforceableApis"
	filterDefaultEnvironment    = "environment eq '~Default'"
)

...

values.Add("api-version", apiVersion)
values.Add(showApisWithToS, "true")
values.Add(hideDlpExemptApis, "true")
values.Add(showAllDlpEnforceableApis, "true")
values.Add("$filter", filterDefaultEnvironment)
```

This improves code readability and maintainability.
