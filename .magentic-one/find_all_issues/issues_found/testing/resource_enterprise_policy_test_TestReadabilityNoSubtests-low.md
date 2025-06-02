# Issue: Test Readability—No Subtests Used

## 
/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/resource_enterprise_policy_test.go

## Problem

The test function is monolithic. Using subtests (with `t.Run`) could clarify separate behaviors and make failures easier to diagnose.

## Impact

Single large tests are harder to debug, and it's more difficult to add new variants or granular checks. Severity: low.

## Location

Whole test function.

## Code Issue

```go
func TestUnitTestEnterpisePolicyResource_Validate_Create(t *testing.T) {
    // ... all logic in one function
}
```

## Fix

Split into subtests using `t.Run`, especially if you expand scenarios with negative or edge cases:

```go
func TestUnitTestEnterprisePolicyResource_Validate_Create(t *testing.T) {
    t.Run("Valid config creates and links policy", func(t *testing.T) {
        // happy path
        // ... current setup and assertions ...
    })

    t.Run("Invalid policy type returns error", func(t *testing.T) {
        // ... test that triggers an error and checks ExpectError ...
    })

    // More subtests as needed...
}
```
