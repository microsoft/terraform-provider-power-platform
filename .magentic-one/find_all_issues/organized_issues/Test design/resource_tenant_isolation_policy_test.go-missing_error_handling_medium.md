# Title

Missing Error Handling in HTTP Mock Registration Logic

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy_test.go

## Problem

Within the test functions (such as `TestUnitTenantIsolationPolicyResource_Validate_Update` and similar), the code repeatedly calls various `httpmock.RegisterResponder` and `httpmock.RegisterResponder` nested registrations without checking or handling errors returned by the registration function. The `RegisterResponder` function can return an error if the registration fails (e.g., due to misconfiguration), but its return value is ignored.

## Impact

By not handling errors from `RegisterResponder`, tests may become unreliable—the test setup may silently fail, causing later assertions to fail without a clear cause. This can increase debugging time and result in misleading test results. The severity is **medium**, as it affects reliability, but the tests can still run under normal circumstances assuming no registration error.

## Location

Throughout multiple test and setup functions, e.g.:

```go
httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
    func(req *http.Request) (*http.Response, error) {
        // ...
    })
```

## Code Issue

```go
httpmock.RegisterResponder("PUT", fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
    func(req *http.Request) (*http.Response, error) {
        // function body omitted for brevity
    })
```

## Fix

Capture the error returned by `RegisterResponder` and fail the test immediately if an error occurs during mock setup, which is standard best practice in Go tests:

```go
err := httpmock.RegisterResponder(
    "PUT",
    fmt.Sprintf("https://api.bap.microsoft.com/providers/PowerPlatform.Governance/v1/tenants/%s/tenantIsolationPolicy", testTenantID),
    func(req *http.Request) (*http.Response, error) {
        // function body omitted for brevity
    },
)
if err != nil {
    t.Fatalf("failed to register HTTP mock: %v", err)
}
```

Repeat this pattern for all `RegisterResponder` calls inside test functions to ensure setup errors are immediately visible.
