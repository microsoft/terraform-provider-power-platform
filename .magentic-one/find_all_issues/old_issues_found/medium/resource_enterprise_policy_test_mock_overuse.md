# Title

Potential Overuse of Mock Responses in Unit Tests

##

`/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/resource_enterprise_policy_test.go`

## Problem

The unit test registers a large number of HTTP mock responders for various endpoints. Although mocking is useful for isolating tests from external systems, excessive mocking can make tests difficult to understand and maintain. Each mock responder introduces complexity, making debugging harder when tests fail.

## Impact

This overuse of mocks makes tests fragile and increases maintenance overhead whenever the mocked APIs change. It also diminishes the readability of the test file as developers need to understand the purpose and setup of each mock. Severity: **Medium**

## Location

The issue resides in the repeated use of `httpmock.RegisterResponder` in the `TestUnitTestEnterpisePolicyResource_Validate_Create` function.

---

## Code Issue

```go
httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
	func(req *http.Request) (*http.Response, error) {
		id := httpmock.MustGetSubmatch(req, 1)
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/Validate_Create/get_environment_%s.json", id)).String()), nil
	},
)

httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
	func(req *http.Request) (*http.Response, error) {
		resp := httpmock.NewStringResponse(http.StatusAccepted, "")
		resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
		return resp, nil
	},
)

// Other similar mock responders...
```

---

## Fix

Reduce the number of static mock responders by consolidating mock endpoints where feasible or introducing shared test helpers to streamline setup:

Create a mock initialization helper:
```go
func setupMockResponders() {
	httpmock.RegisterResponder("GET", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			id := httpmock.MustGetSubmatch(req, 1)
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File(fmt.Sprintf("tests/Validate_Create/get_environment_%s.json", id)).String()), nil
		},
	)

	httpmock.RegisterResponder("DELETE", `=~^https://api\.bap\.microsoft\.com/providers/Microsoft\.BusinessAppPlatform/scopes/admin/environments/([\d-]+)\z`,
		func(req *http.Request) (*http.Response, error) {
			resp := httpmock.NewStringResponse(http.StatusAccepted, "")
			resp.Header.Add("Location", "https://europe.api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/lifecycleOperations/b03e1e6d-73db-4367-90e1-2e378bf7e2fc?api-version=2023-06-01")
			return resp, nil
		},
	)

	// Consolidate further mock registrations here as needed.
}
```

Call this helper function in the test setup:
```go
func TestUnitTestEnterprisePolicyResource_Validate_Create(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()
	setupMockResponders()
	// Rest of the test code...
}
```

This approach improves readability, centralizes mock setup, and reduces chances of duplication or conflicting responders.
