# Title

Missing Assertions for HTTP Mock Response Properties

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings_test.go

## Problem

While HTTP mock responders are used to simulate API responses, the tests do not validate critical response properties, such as proper HTTP response status codes or specific output structure from the mocks. This results in the tests merely asserting that a responder is invoked but failing to validate if the response aligns with expected behavior.

## Impact

- **Reduced Test Accuracy**: Inconsistent response properties may go unnoticed, leading to undetected errors during API interaction simulations.
- **Test Drift**: Future updates to the API or mock structure might invalidate test relevance.
- **High Severity**: This compromises the reliability of API-related unit tests, as results from mocks do not ensure consistent interaction expectations.

## Location

Located in both the `TestUnitTestTenantSettingsResource_Validate_Create` and `TestUnitTestTenantSettingsResource_Validate_Update` functions.

## Code Issue

Mock responder lacks assertion on critical response properties:

```go
httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2020-08-01",
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resources/Validate_Create/get_tenant.json").String()), nil
	})
```

## Fix

Add assertions for HTTP response status codes, headers, and mock payload structure after the responder's invocation. While running mocks, use `Assert` functions or even manual checks, matching expectations to reality.

```go
func TestUnitTestTenantSettingsResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// Add an assertion to ensure proper response structure or status code
	httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2020-08-01",
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resources/Validate_Create/get_tenant.json").String())
			if resp.StatusCode != http.StatusOK {
				t.Errorf("Expected status OK but got %v", resp.StatusCode)
			}
			return resp, nil
		})
}
```