# Title

Repeated Test Setup Code is Not Modular

## Path

`/workspaces/terraform-provider-power-platform/internal/provider/provider_test.go`

## Problem

The test setup code, such as activating HTTP mocks, registering responders, and test execution (`test.Test()`), is duplicated across multiple unit tests. This violates the DRY (Don't Repeat Yourself) principle.

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()

mocks.ActivateEnvironmentHttpMocks()

httpmock.RegisterRegexpResponder(
	"GET",
	regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments.*api-version=2023-06-01$`),
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../services/environment/tests/datasource/Validate_Read/get_environments.json").String()), nil
	},
)
```

## Impact

- **Code Duplication**: Leads to inconsistencies if setup logic changes.
- **Reduced Maintainability**: Developers need to make manual adjustments across tests, increasing maintenance overhead.
- **Medium Severity**: Indirectly impacts the scalability and maintainability of the test file.

## Location

- Lines 102–117 (`TestUnitPowerPlatformProvider_Validate_Telemetry_Optout_Is_False`)
- Lines 128–143 (`TestUnitPowerPlatformProvider_Validate_Telemetry_Optout_Is_True`)

## Code Issue

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()

mocks.ActivateEnvironmentHttpMocks()

httpmock.RegisterRegexpResponder(
	"GET",
	regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments.*api-version=2023-06-01$`),
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../services/environment/tests/datasource/Validate_Read/get_environments.json").String()), nil
	},
)
```

## Fix

Extract the repeated setup code into a helper function or utility method that encapsulates common setup logic. This approach reduces duplication and improves maintainability.

```go
func setupTestEnvironment(httpMockConfig func()) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	if httpMockConfig != nil {
		httpMockConfig()
	}
}

// Usage in test functions:
func TestUnitPowerPlatformProvider_Validate_Telemetry_Optout_Is_False(t *testing.T) {
	setupTestEnvironment(func() {
		httpmock.RegisterRegexpResponder(
			"GET",
			regexp.MustCompile(`^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments.*api-version=2023-06-01$`),
			func(req *http.Request) (*http.Response, error) {
				return httpmock.NewStringResponse(http.StatusOK, httpmock.File("../services/environment/tests/datasource/Validate_Read/get_environments.json").String()), nil
			},
		)
	})

	test.Test(t, test.TestCase{
		// Test case setup...
	})
}
```