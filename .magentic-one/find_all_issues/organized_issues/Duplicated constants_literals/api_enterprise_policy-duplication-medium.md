# Function Length/Structure & Code Duplication

##

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/api_enterprise_policy.go

## Problem

Both `LinkEnterprisePolicy` and `UnLinkEnterprisePolicy` share significant duplicated logic (URL construction, request building, API call, waiting). This makes maintenance harder and increases risk of bugs. Separation of concerns is weak; logging, error handling, and API calls are mixed.

## Impact

Severity: **medium**. Readability and maintainability are reduced; increases technical debt and future maintenance risk.

## Location

Functions `LinkEnterprisePolicy` and `UnLinkEnterprisePolicy`.

## Code Issue

```go
// LinkEnterprisePolicy
apiUrl := &url.URL{ ... }
values := url.Values{}
values.Add("api-version", "2019-10-01")
apiUrl.RawQuery = values.Encode()
linkEnterprosePolicyDto := linkEnterprosePolicyDto{
	SystemId: systemId,
}
apiResponse, err := client.Api.Execute(...)
// UnLinkEnterprisePolicy
// (Repeated code)
```

## Fix

Refactor to extract common code (e.g., URL construction, request execution) into helper functions or methods.

```go
func buildEnterprisePolicyURL(baseURL, action, envId, envType string) string {
	return fmt.Sprintf("https://%s/providers/Microsoft.BusinessAppPlatform/environments/%s/enterprisePolicies/%s/%s?api-version=2019-10-01", 
		baseURL, envId, envType, action)
}
// Then use
linkURL := buildEnterprisePolicyURL(client.Api.GetConfig().Urls.BapiUrl, "link", environmentId, environmentType)
unlinkURL := buildEnterprisePolicyURL(client.Api.GetConfig().Urls.BapiUrl, "unlink", environmentId, environmentType)
// And so on.
```
