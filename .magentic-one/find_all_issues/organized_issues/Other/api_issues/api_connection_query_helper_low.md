# Title

Inefficient Use of url.Values: Always Initialize Even if Used Once per Request

##

/workspaces/terraform-provider-power-platform/internal/services/connection/api_connection.go

## Problem

In each method, `values := url.Values{}` is created even in cases where it is only used to add a couple of fixed query parameters, and then immediately encoded to the URL's `RawQuery`. Inlining the creation and setting could reduce code verbosity, and a helper function could further improve maintainability.

## Impact

This is a minor maintainability and readability issue, possibly leading to repetitive boilerplate code for setting URL parameters. Severity: Low.

## Location

Repeated pattern, for example:

```go
values := url.Values{}
values.Add("api-version", "1")
values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
apiUrl.RawQuery = values.Encode()
```

## Code Issue

```go
values := url.Values{}
values.Add("api-version", "1")
values.Add("$filter", fmt.Sprintf("environment eq '%s'", environmentId))
apiUrl.RawQuery = values.Encode()
```

## Fix

Introduce a utility/helper function (e.g., `buildQuery`) to centralize query construction:

```go
func buildQuery(environmentId string) string {
    values := url.Values{}
    values.Add(apiVersionKey, apiVersionVal)
    values.Add(filterKey, fmt.Sprintf("environment eq '%s'", environmentId))
    return values.Encode()
}
```
Then in each method:

```go
apiUrl.RawQuery = buildQuery(environmentId)
```

---

*This file should be saved in:*
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/api_connection_query_helper_low.md
