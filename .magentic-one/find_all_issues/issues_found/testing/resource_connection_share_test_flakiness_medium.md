# Title

Potential test flakiness and performance due to resource-heavy acceptance test

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_share_test.go

## Problem

The `TestAccConnectionsShareResource_Validate_Create` function appears to be a full acceptance test which creates Azure Active Directory groups, users, random passwords, Power Platform environments, connections, and shares. This degree of resource creation can lead to a long and potentially flaky test run, especially if run often or in CI/CD, due to real cloud dependencies and possible resource cleanup issues.

## Impact

- Makes test suite slow, potentially causing timeouts in CI.
- Raises the risk of rate limiting or hitting quota limits/financial costs with external providers.
- Test failures may occur unrelated to the code being tested (e.g., due to Azure service outages or delays).
- Difficulties in resource cleanup can leave orphan resources, creating further pollution/noise.

Severity: medium.

## Location

Entire function:

```go
func TestAccConnectionsShareResource_Validate_Create(t *testing.T) { ... }
```

## Code Issue

```go
func TestAccConnectionsShareResource_Validate_Create(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
        ExternalProviders: map[string]resource.ExternalProvider{
            ...
        },
        Steps: []resource.TestStep{ ... },
    })
}
```

## Fix

- Acceptance tests should be tagged or filtered so they do not run by default as part of unit/integration test suites (`go test -short`).
- Ensure that all created resources are reliably destroyed with a `resource.TestStep` or equivalent teardown logic.
- Consider splitting exhaustive scenarios into lighter unit tests using HTTP mocks wherever possible for most logic verification, and restrict acceptance tests to a minimum essential use case.
- If using in CI/CD, restrict the frequency of these acceptance tests and ensure they run in parallel with quotas in mind.

Example: 

```go
func TestAccConnectionsShareResource_Validate_Create(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping acceptance test in short mode.")
    }
    // ... rest of test logic
}
```

---

This issue will be saved under `/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/testing/resource_connection_share_test_flakiness_medium.md`.
