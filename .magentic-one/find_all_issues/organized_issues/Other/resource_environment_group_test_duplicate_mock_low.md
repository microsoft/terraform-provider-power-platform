# Duplicate HTTP Mock Registration for Same Request

##

/workspaces/terraform-provider-power-platform/internal/services/environment_groups/resource_environment_group_test.go

## Problem

There are two registrations for the exact same HTTP GET request to:

```
https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24filter=properties%2FparentEnvironmentGroup%2Fid+eq+00000000-0000-0000-0000-000000000001&api-version=2021-04-01
```

This can lead to confusion over which responder is active or can inadvertently override a previous responder.

## Impact

This introduces maintainability and readability issues and can cause confusion as to which responder will actually be used. In `httpmock`, the last registered responder overrides previous entries for the same method+url. Severity: **low**

## Location

```go
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24filter=properties%2FparentEnvironmentGroup%2Fid+eq+00000000-0000-0000-0000-000000000001&api-version=2021-04-01",
		httpmock.NewStringResponder(http.StatusOK, `{"value":[]}`),
	)

// ...other registrations here...

	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24filter=properties%2FparentEnvironmentGroup%2Fid+eq+00000000-0000-0000-0000-000000000001&api-version=2021-04-01",
		httpmock.NewStringResponder(http.StatusOK, `{"value":[]}`),
	)
```

## Code Issue

```go
	httpmock.RegisterResponder("GET", "<...same-url...>",
		httpmock.NewStringResponder(http.StatusOK, `{"value":[]}`),
	)
	// (Repeated above and again below.)
```

## Fix

Remove the duplicate registration. Only one registration is needed for a given HTTP method and URL.

```go
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?%24filter=properties%2FparentEnvironmentGroup%2Fid+eq+00000000-0000-0000-0000-000000000001&api-version=2021-04-01",
		httpmock.NewStringResponder(http.StatusOK, `{"value":[]}`),
	)
```
