# Title

Inadequate comments for logical flow explanation in test cases

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy_test.go

## Problem

Several functions and logic pieces in test cases lack detailed comments explaining their purpose or thought process. While basic comments are present, they do not sufficiently describe certain complex flows, such as the mechanics of state transitions in mocked responses.

## Impact

Developers and collaborators new to the codebase may find it harder to quickly comprehend the purpose of various test steps and logic flows without comprehensive comments. This impacts code maintainability, albeit with a low severity.

## Location

Throughout the file, e.g., near functions `TestAccTenantIsolationPolicy_Validate_Update` and `TestUnitTenantIsolationPolicyResource_Validate_Update`.

## Code Issue

```go
// Mock tenant endpoint that's called before CRUD operations.
httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2021-04-01",
// Excludes details about CRUD operation transitions and specific mocks for individual methods in test steps.
```

## Fix

Enhance comments with detailed annotations about test logic, expected results, and assumptions.

```go
// Mock tenant endpoint that's called before CRUD operations.
// This endpoint simulates the initial tenant retrieval process. It's required to initialize state
// transition testing for all CRUD test scenarios.
httpmock.RegisterResponder("GET", "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/tenant?api-version=2021-04-01",
```