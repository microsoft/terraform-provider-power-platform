# Title

Potential Redundancy in Variable Declaration

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy_test.go

## Problem

Variables such as `rgName` are generated multiple times across test blocks using repeated patterns. While not functionally incorrect, they introduce unnecessary redundancy that can be simplified for cleaner tests.

## Impact

Although the impact is low, redundant code can impact readability and future maintainability. Severity: Low.

## Location

Examples found at lines 15, 121, 175.

## Code Issue

```go
rgName := "power-platform-billing-" + mocks.TestName() + strconv.Itoa(rand.IntN(9999))
```

## Fix

Create a helper method or generator to standardize resource group naming conventions.

```go
func GenerateResourceGroupName(base string) string {
    return base + mocks.TestName() + strconv.Itoa(rand.IntN(9999))
}

// Usage in tests:
rgName := GenerateResourceGroupName("power-platform-billing-")
```
