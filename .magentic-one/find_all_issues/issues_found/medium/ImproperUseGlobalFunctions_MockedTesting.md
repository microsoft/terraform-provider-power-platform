# Title

Improper Use of Global Functions Without Mocked Testing

##

/workspaces/terraform-provider-power-platform/internal/services/environment/resource_environment_test.go

## Problem

The code uses global functions like `httpmock.RegisterResponder` without properly mocking dependencies. This can lead to brittle tests as external dependencies may change unexpectedly, causing test failures.

## Impact

Using non-mocked global functions in testing can result in unreliable tests that depend on external factors such as API behavior or network connectivity. The severity of this issue is medium because it affects test stability but does not break the logic.

## Location

Several places in the file `/resource_environment_test.go`:

Example:
1. Line: `httpmock.RegisterResponder` call is utilized without dependency injection or proper isolation of `HTTPMock`.

## Code Issue

```go
	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})
```

## Fix

Replace external HTTP calls with mocked ones and use dependency injection for `HTTPMock`. This ensures the tests remain reliable and isolated from actual external functionality.

```go
func mockHTTPCalls() {
	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/00000000-0000-0000-0000-000000000001?api-version=2023-06-01")
			return resp, nil
		})
}

// Call this in test setup
mockHTTPCalls()
```

By centralizing the mocking into a clear setup call, you make the dependency injection clear and the test isolations robust, improving reliability.