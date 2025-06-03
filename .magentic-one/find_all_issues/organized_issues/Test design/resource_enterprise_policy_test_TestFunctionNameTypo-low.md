# Issue: Test Function Name Typo

## 
/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/resource_enterprise_policy_test.go

## Problem

The test function is named `TestUnitTestEnterpisePolicyResource_Validate_Create`. The word "Enterpise" appears to be a typo and should be "Enterprise".

## Impact

Incorrect naming reduces readability and consistency, making it harder to find and reference tests. Severity: low.

## Location

Line 13, function definition.

## Code Issue

```go
func TestUnitTestEnterpisePolicyResource_Validate_Create(t *testing.T) {
```

## Fix

Rename the function to use the correct spelling ("Enterprise"):

```go
func TestUnitTestEnterprisePolicyResource_Validate_Create(t *testing.T) {
```
