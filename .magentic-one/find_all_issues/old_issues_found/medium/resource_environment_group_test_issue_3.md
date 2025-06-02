# Issue 3

## Hardcoded Strings in HTTP Mock Responders

##

`/workspaces/terraform-provider-power-platform/internal/services/environment_groups/resource_environment_group_test.go`

## Problem

Several HTTP mock responders use hardcoded strings, such as URLs and JSON responses, directly in the code.

## Impact

Hardcoded values reduce maintainability and flexibility. If these values change in future API versions, it will require significant effort to locate and update them across the file. This can introduce bugs or inconsistencies in tests when updates occur. **Severity: Medium**

### Location

Example hardcoded sections:

```go
httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24filter=properties%2FparentEnvironmentGroup%2Fid+eq+00000000-0000-0000-0000-000000000001&api-version=2021-04-01",
	httpmock.NewStringResponder(http.StatusOK, `{"value":[]}`),
)
```

## Code Issue

```go
httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24filter=properties%2FparentEnvironmentGroup%2Fid+eq+00000000-0000-0000-0000-000000000001&api-version=2021-04-01",
	httpmock.NewStringResponder(http.StatusOK, `{"value":[]}`),
)
```

## Fix

Define constants for all URLs and responses to centralize these values. This makes them easier to update and maintain:

```go
const (
	apiURLForEnvironmentGroup = "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24filter=properties%2FparentEnvironmentGroup%2Fid+eq+00000000-0000-0000-0000-000000000001&api-version=2021-04-01"
	mockResponseForEnvironmentGroup = `{"value":[]}`
)

httpmock.RegisterResponder("GET", apiURLForEnvironmentGroup,
	httpmock.NewStringResponder(http.StatusOK, mockResponseForEnvironmentGroup),
)
```