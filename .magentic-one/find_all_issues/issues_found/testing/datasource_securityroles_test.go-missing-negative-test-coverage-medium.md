# Title

Missing Test Coverage for Error Scenarios (Edge Cases)

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/datasource_securityroles_test.go

## Problem

The current test suite focuses on the happy path and a specific single error scenario (`No Dataverse exists in environment`). However, it lacks extensive testing for other error cases, such as handling malformed responses, HTTP failures (e.g., 500-series errors), missing attributes in role objects, or connectivity issues.

## Impact

Insufficient negative test coverage can allow critical bugs or regressions to enter production, especially when new code changes are made. The severity is **medium** as it could lead to undetected failures in production if error paths are not exercised by tests.

## Location

Throughout all test functions—none check for:
- non-200/Accepted HTTP status codes,
- corrupted/invalid mock JSON files,
- missing/invalid fields in returned data,
- API connection errors.

## Code Issue

```go
// Example from the file (only one error scenario tested)
{
    Config: `
    resource "powerplatform_environment" "env" { ... }
    data "powerplatform_security_roles" "all" {
        environment_id = powerplatform_environment.env.id
    }`,
    ExpectError: regexp.MustCompile(`No Dataverse exists in environment`),
    Check:       resource.ComposeTestCheckFunc(),
},
// Other cases lack ExpectError/negative test paths entirely
```

## Fix

Add additional unit tests for error scenarios, simulating API failures (500, 404), invalid JSON, and missing/invalid response fields.

```go
func TestUnitSecurityDataSource_InvalidResponses(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()
    mocks.ActivateEnvironmentHttpMocks()

    // Simulate a 500 error on environment call
    httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments/00000000-0000-0000-0000-000000000001?...`,
        func(req *http.Request) (*http.Response, error) {
            return httpmock.NewStringResponse(http.StatusInternalServerError, ""), nil
        })

    resource.Test(t, resource.TestCase{
        IsUnitTest:               true,
        ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: `
                    data "powerplatform_security_roles" "all" {
                        environment_id = "00000000-0000-0000-0000-000000000001"
                    }`,
                ExpectError: regexp.MustCompile(`500`), // Adjust as needed for error message
            },
        },
    })
}

// Repeat for other error scenarios: invalid JSON, missing fields, etc.
```

This increases robustness and ensures all critical paths are exercised in testing.

---

This issue will be saved as  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/testing/datasource_securityroles_test.go-missing-negative-test-coverage-medium.md`.

If any further issues are detected, I will continue.
