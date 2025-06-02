# Title

Insufficient Use of Dependency Injection

## Path

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/resource_environment_settings_test.go

## Problem

The mocks and HTTP requests within the test methods are using static, direct calls (`httpmock.RegisterResponder`, etc.), making them tightly coupled with the mock framework. Dependent configurations or behaviors cannot be overridden or dynamically adjusted at runtime. This limits the flexibility of tests.

## Impact

Medium — Tight coupling may reduce the scalability and reusability of the test suite. As the project grows, managing tightly coupled static calls can lead to significant overhead and debugging complexity.

## Location

Instances such as:

```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01`,
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resources/Validate_Create_Empty_Settings/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
	})
```

## Code Issue

```go
httpmock.RegisterResponder("GET", `<endpoint_path>`,
    func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusOK, "<dynamic_mock_response>"), nil
    })
```

## Fix

Introduce dependency injection to load test configuration dynamically rather than relying directly on static mocks. For example:

```go
func registerMockResponse(endpoint string, responseFilePath string) {
	httpmock.RegisterResponder("GET", endpoint,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(responseFilePath).String()), nil
		})
}

// Usage
registerMockResponse(
	"https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?api-version=2023-06-01",
	"tests/resources/Validate_Create_Empty_Settings/get_environment_00000000-0000-0000-0000-000000000001.json",
)
```