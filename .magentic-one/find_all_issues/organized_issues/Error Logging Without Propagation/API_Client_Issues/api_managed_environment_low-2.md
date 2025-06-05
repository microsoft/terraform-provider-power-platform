# Inefficient Repeated Construction of URL Query Parameters

##

/workspaces/terraform-provider-power-platform/internal/services/managed_environment/api_managed_environment.go

## Problem

Both `EnableManagedEnvironment` and `DisableManagedEnvironment` construct almost identical `apiUrl` and query parameter blocks, duplicating the logic. Similarly, `FetchSolutionCheckerRules` builds another, slightly different, API URL with minor variation. These can be encapsulated to avoid code repetition, reduce error risk, and improve maintainability.

## Impact

- **Low severity**
- Code duplication leads to higher maintenance overhead.
- Increases risk of divergence and subtle bugs if updates are made in one place but not the other.

## Location

- `EnableManagedEnvironment`  
- `DisableManagedEnvironment`  
- `FetchSolutionCheckerRules`

## Code Issue

```go
apiUrl := &url.URL{
    Scheme: constants.HTTPS,
    Host:   client.Api.GetConfig().Urls.BapiUrl,
    Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/environments/%s/governanceConfiguration", environmentId),
}
values := url.Values{}
values.Add("api-version", "2021-04-01")
apiUrl.RawQuery = values.Encode()
```
(similar logic elsewhere)

## Fix

Extract the URL construction into a helper method:

```go
func bapiGovernanceUrl(baseHost, environmentId string) string {
    apiUrl := &url.URL{
        Scheme: constants.HTTPS,
        Host:   baseHost,
        Path:   fmt.Sprintf("/providers/Microsoft.BusinessAppPlatform/environments/%s/governanceConfiguration", environmentId),
    }
    values := url.Values{}
    values.Add("api-version", "2021-04-01")
    apiUrl.RawQuery = values.Encode()
    return apiUrl.String()
}
```
Then reuse:

```go
apiUrl := bapiGovernanceUrl(client.Api.GetConfig().Urls.BapiUrl, environmentId)
```

---

This will be saved to:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/api_managed_environment_low.md`
