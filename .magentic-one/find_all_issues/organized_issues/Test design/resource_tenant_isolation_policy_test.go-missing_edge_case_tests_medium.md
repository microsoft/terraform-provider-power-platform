# Title

Test Coverage: Missing Negative and Edge Case Handling

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy_test.go

## Problem

The test file primarily covers standard create, update, and delete paths—with one explicit error case for invalid tenant IDs. However, there is limited or no testing for additional important negative and edge cases, such as:

- Simulating and verifying behavior when mock files are missing or malformed.
- HTTP responses with unexpected codes (e.g., `500`, `403`, timeouts).
- Network errors and other external failure conditions.
- API returns content that does not match the expected schema (malformed body).
- Testing policy creation with overlapping or duplicate tenants.

## Impact

Not covering edge cases increases the risk that the code will break in unexpected production scenarios, reducing the reliability and robustness of the resource provider. This is an issue of **medium** severity, especially for infrastructure automation providers.

## Location

Test suite scope/coverage (all test functions in the file):

```go
// Most test functions simulate normal or single validation error paths, but do not simulate network, server, or file infrastructure failures.
```

## Code Issue

```go
// NO direct invocation or assertion of error handling when, e.g., `httpmock.File()` is called with a missing file,
// or when API returns 500, 403, or malformed JSON responses.
```

## Fix

Augment the test suite to add functions covering additional error, negative, and corner cases. For example:

```go
func TestTenantIsolationPolicy_HandleHTTP500Error(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    // Setup all required responders (create, get, etc.)
    err := httpmock.RegisterResponder(
        "GET",
        fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
        func(req *http.Request) (*http.Response, error) {
            return httpmock.NewStringResponse(http.StatusInternalServerError, "Internal Server Error"), nil
        },
    )
    if err != nil {
        t.Fatalf("Failed to register responder: %v", err)
    }

    // ...set up test and verify failure/expected error...
}
```

Create similar functions to test for file not found, network timeouts, permission denied (`403`), and malformed response bodies.
