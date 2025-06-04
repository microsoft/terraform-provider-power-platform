# Maintainability: Duplicate Code for Building API URLs

##

/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/api_environment_group_rule_set.go

## Problem

Every method manually builds the API URL object, sets scheme, host, path, and adds the version query parameter. This results in repetitive boilerplate and increases risk of errors and inconsistencies.

## Impact

Hurts maintainability and increases likelihood of bugs, especially if the base pattern for URL structure or version ever needs to change. Refactoring is more difficult. Severity: Medium.

## Location

Repeated in every client method:
```go
apiUrl := &url.URL{ ... }
values := url.Values{}
values.Add("api-version", ...)
apiUrl.RawQuery = values.Encode()
```

## Code Issue

```go
apiUrl := &url.URL{
    Scheme: constants.HTTPS,
    Host:   helpers.BuildTenantHostUri(tenantDto.TenantId, client.Api.GetConfig().Urls.PowerPlatformUrl),
    Path:   ...,
}
values := url.Values{}
values.Add("api-version", "2021-10-01-preview")
apiUrl.RawQuery = values.Encode()
```

## Fix

Refactor common URL-building logic into a helper method to DRY up the codebase:

```go
func buildEnvironmentGroupRuleSetURL(tenantID, baseHost, path, version string) *url.URL {
    apiUrl := &url.URL{
        Scheme: constants.HTTPS,
        Host:   helpers.BuildTenantHostUri(tenantID, baseHost),
        Path:   path,
    }
    values := url.Values{}
    values.Add("api-version", version)
    apiUrl.RawQuery = values.Encode()
    return apiUrl
}
// Use this in each method.
```

---

This issue will be saved in:
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/api_environment_group_rule_set_structure_medium_boilerplate.md
