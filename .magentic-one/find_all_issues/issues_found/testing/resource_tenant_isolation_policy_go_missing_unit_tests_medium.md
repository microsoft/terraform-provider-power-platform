# Lack of Unit Tests for Resource Methods

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy.go

## Problem

The file implements resource logic, schema, and validation but lacks corresponding unit tests. There is no indication of tests for custom validation logic (`ValidateConfig`), error paths, resource Create/Read/Update/Delete logic, or their return conditions—especially the error or diagnostic paths.

## Impact

- **Severity: Medium**
- Reduces confidence in code correctness.
- Increases the risk of regressions.
- Poor test coverage can hide bugs, particularly due to the custom diagnostics/error handling in schema and CRUD operations.

## Location

_All CRUD methods, validation, and all custom logic. (No test coverage present for any of them, which would normally reside in a *_test.go file or be explicitly visible through interface injections/mocks.)_

## Code Issue

```go
// No test function or interface mock found for any of the following:
//   - ValidateConfig()
//   - Create()
//   - Read()
//   - Update()
//   - Delete()
//   - ImportState()
// And all error/diagnostic branches.
```

## Fix

Implement thorough unit tests for:

- The custom validators in `ValidateConfig`, especially error and success cases.
- The error paths and state transitions in Create, Update, Delete, Read, and ImportState.
- The happy-path and error branches for API client errors, missing IDs, and diagnostic reports.
- Consider dependency injection or mocks for API interactions.

```go
func TestValidateConfig_InvalidAllowedTenant(t *testing.T) {
    // Initialize required resource and config
    req := resource.ValidateConfigRequest{/* ...fill config with invalid AllowedTenant */} 
    resp := &resource.ValidateConfigResponse{}
    res := NewTenantIsolationPolicyResource()
    res.ValidateConfig(context.Background(), req, resp)
    if !resp.Diagnostics.HasError() {
        t.Fatal("Expected diagnostic error for invalid tenant configuration")
    }
}

// Repeat similar patterns for other methods and error conditions
```
