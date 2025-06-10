# Title

Possible Improvement: Use of Constants for Repeated URL Query Keys

##

/workspaces/terraform-provider-power-platform/internal/services/connection/api_connection.go

## Problem

The `api-version` and `$filter` strings are repeatedly used to add query parameters to URLs, but are hardcoded each time. Defining these as package-level constants would improve maintainability and prevent typos.

## Impact

Improves code maintainability, reduces the risk of copy-paste mistakes, and clarifies intent. Severity: Low.

## Location

For example, in many methods:

```go
values := url.Values{}
values.Add("api-version", "1")
values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
```

## Code Issue

```go
values.Add("api-version", "1")
values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
```

## Fix

Define package-level constants:

```go
const (
    apiVersionKey = "api-version"
    apiVersionVal = "1"
    filterKey     = "$filter"
)
```
and use:

```go
values.Add(apiVersionKey, apiVersionVal)
values.Add(filterKey, fmt.Sprintf("environment eq '%s'", environmentId))
```

---

*This file should be saved in:*
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/api_connection_query_constants_low.md
