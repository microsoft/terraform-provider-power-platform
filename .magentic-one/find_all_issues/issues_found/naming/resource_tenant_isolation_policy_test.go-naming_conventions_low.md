# Title

Test Function Naming Does Not Follow Go Naming Conventions

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy_test.go

## Problem

Several test functions, such as `TestUnitTenantIsolationPolicyResource_Validate_Create` and `TestAccTenantIsolationPolicy_Validate_Create`, use highly verbose and non-idiomatic naming. The Go convention for test function names is to be more concise, using underscores sparingly, and to use descriptive but succinct names since the context ("test") is already clear, and details should go in comments.

## Impact

Non-idiomatic naming reduces the readability and approachability of the test suite for Go developers, making it harder to quickly understand and run tests. It may also make integration with some Go test tools more awkward. Severity: **low**, as it does not prevent the code from running, but detracts from codebase maintainability.

## Location

Function definitions throughout the file, such as:

```go
func TestUnitTenantIsolationPolicyResource_Validate_Create(t *testing.T) {
    // ...
}
```

## Code Issue

```go
func TestUnitTenantIsolationPolicyResource_Validate_Create(t *testing.T) {
    // ...
}
```

## Fix

Rename the test functions to follow Go conventions; for example, focus on succinctly describing the subject and scenario:

```go
func TestTenantIsolationPolicy_CreateUnit(t *testing.T) {
    // ...
}
```

or

```go
func TestTenantIsolationPolicy_Create(t *testing.T) {
    // ...
}
```

Add scenario details in descriptive comments above the function body. Apply similar changes for all test functions in the file.
