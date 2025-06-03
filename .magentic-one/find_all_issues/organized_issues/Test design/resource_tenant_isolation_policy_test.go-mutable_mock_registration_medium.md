# Title

Possible Test Flakiness from Dynamic Mock Registration in Nested Handlers

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy_test.go

## Problem

Some unit test functions dynamically register new HTTP responders inside existing responder/handler functions (i.e., as a side effect of handling a test PUT or GET request). This pattern is used to simulate changing backend state across test steps. However, it can introduce race conditions and unpredictable test results if multiple requests are issued "out of order," or if parallel test execution is ever adopted.

## Impact

Dynamic stateful registration inside handlers can result in flakiness if the test order or concurrency model changes, or as the test grows more complex. Future maintenance or test authors may be unaware of the global, mutable nature of these responders. This is a **medium** severity issue due to potential confusion and flakiness in test behavior.

## Location

Example from `TestUnitTenantIsolationPolicyResource_Validate_Update`:

```go
// Step 1: First PUT creates policy with initial state.
httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
    func(req *http.Request) (*http.Response, error) {
        // After first PUT, register a new GET to return initial state.
        httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
            func(req *http.Request) (*http.Response, error) {
                return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Update/initial_policy.json").String()), nil
            })

        // Register a new PUT handler for the update operation.
        httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
            func(req *http.Request) (*http.Response, error) {
                // After second PUT, register a new GET to return updated state.
                httpmock.RegisterResponder("GET", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
                    func(req *http.Request) (*http.Response, error) {
                        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Update/updated_policy.json").String()), nil
                    })
                return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Update/updated_policy.json").String()), nil
            })
        return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/resource/Validate_Update/initial_policy.json").String()), nil
    })
```

## Code Issue

```go
func(req *http.Request) (*http.Response, error) {
    // After first PUT, register a new GET to return initial state.
    httpmock.RegisterResponder("GET", url, func(...) {...})

    // Register a new PUT handler for the update operation.
    httpmock.RegisterResponder("PUT", url, func(...) {
        // register further handlers here
    })

    // ...
}
```

## Fix

Structure test setup so that all responder state is defined before the test steps begin. Use variables or an in-memory map (static to the test function) to represent the mock server's state, and provide handlers that check state and return the appropriate response. This avoids mutable global registration during tests.

Example pattern:

```go
var state = "initial"
httpmock.RegisterResponder("PUT", url, func(req *http.Request) (*http.Response, error) {
    if state == "initial" {
        state = "updated"
        return httpmock.NewStringResponse(http.StatusOK, ...initialContent...), nil
    }
    // handle update, etc.
})
httpmock.RegisterResponder("GET", url, func(req *http.Request) (*http.Response, error) {
    if state == "initial" {
        return httpmock.NewStringResponse(http.StatusOK, ...initialContent...), nil
    }
    return httpmock.NewStringResponse(http.StatusOK, ...updatedContent...), nil
})
```

This pattern is easier to reason about and less brittle with respect to request order/randomness.
