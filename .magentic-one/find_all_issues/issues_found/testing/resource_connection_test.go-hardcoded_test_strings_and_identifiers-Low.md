# Title

Hardcoded Test Strings and Identifiers

##

/workspaces/terraform-provider-power-platform/internal/services/connection/resource_connection_test.go

## Problem

Both tests rely on hardcoded string identifiers and resource names (e.g., `00000000-0000-0000-0000-000000000000` for `environment_id`, `"shared_azureopenai"`, and HCL field names). This approach can make it more difficult to adapt the tests for broader scenarios or to run in parallel with other similar tests.

## Impact

Hardcoded values reduce the reusability and flexibility of the tests and may lead to brittle tests in more complex test rigs. Reuse of identifiers might cause side effects when tests are executed in parallel or with differing external conditions.

Severity: Low

## Location

In both `Config` and regexes for HTTP mocks, as well as resource names.

```go
environment_id = "00000000-0000-0000-0000-000000000000"
name = "shared_azureopenai"
// Regex and HCL with above hardcoded values
```

## Fix

Leverage randomly generated values for IDs (where possible) or unique names, possibly using helper functions. For example, use `mocks.TestName()` or `t.Name()` to generate unique names and avoid collisions.

```go
envID := "00000000-0000-0000-0000-000000000000" // Use a helper to randomize for each test run
resourceName := fmt.Sprintf("azure_openai_conn_%s", t.Name())
```
This practice makes your tests more robust and flexible to changes in test data or parallel running scenarios.
