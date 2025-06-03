# Title

Lack of Test Case Separation and Readability in Test Functions

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/datasource_securityroles_test.go

## Problem

The test functions in this file are overly long and contain setup, mock definitions, resource test logic, and complex inline configuration definitions all in one place. This makes it hard to isolate test setup from the actual assertions and to identify what each test is doing at a glance.

## Impact

This impacts code maintainability and readability. When tests fail, it is harder to diagnose the issue, and when updating or extending the tests, changes can introduce unintended bugs. The severity is **medium**, as it hinders developer productivity and increases the risk of breaking tests during refactoring.

## Location

Throughout the file, especially in functions like `TestUnitSecurityDataSource_Validate_Read` and `TestUnitSecurityDataSource_Validate_No_Dataverse`.

## Code Issue

```go
func TestUnitSecurityDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	mocks.ActivateEnvironmentHttpMocks()

	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/security_roles/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})

	// (more mock responders...)

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// (steps...)
		},
	})
}
```

## Fix

Separate the setup logic, HTTP mock registration, and the actual test logic into helper functions. This not only makes each test function shorter and clearer but also reduces code duplication and improves maintainability.

```go
func setupMockRespondersForValidateRead() {
	httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?%24expand=permissions%2Cproperties.capacity%2Cproperties%2FbillingPolicy%2Cproperties%2FcopilotPolicies&api-version=2023-06-01`,
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/security_roles/Validate_Read/get_environment_00000000-0000-0000-0000-000000000001.json").String()), nil
		})
	httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/roles",
		func(req *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/security_roles/Validate_Read/get_security_roles.json").String()), nil
		})
}

func TestUnitSecurityDataSource_Validate_Read(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	mocks.ActivateEnvironmentHttpMocks()
	setupMockRespondersForValidateRead()

	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// (test steps remain as before)
		},
	})
}
```

This approach should be generalized for all test functions in the file. Consider extracting common test configuration strings to package-level constants for further clarity.

---

This issue will be saved as `/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/structure/datasource_securityroles_test.go-lack-of-test-case-separation-medium.md`.
Proceeding to next issue if found.
