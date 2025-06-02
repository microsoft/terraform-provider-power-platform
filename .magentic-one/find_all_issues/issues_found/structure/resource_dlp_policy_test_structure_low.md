# Lack of Test Function Separation and Documentation

##

`/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy_test.go`

## Problem

All unit and acceptance tests are placed within the same file without clear separation, grouping, or sufficient function documentation. Test readability and maintainability could be significantly improved with proper documentation and better functional decomposition, especially given the complexity of the configurations used in this test suite.

## Impact

**Severity: Low**  
This issue decreases maintainability and makes it harder for collaborators to quickly identify what each test is intended to do or to expand/fix individual behaviors.

## Location

Throughout the file, all test functions lack doc comments and the test suite is lengthy without clear documentation between logical sections.

## Code Issue

```go
func TestUnitDataLossPreventionPolicyResource_Validate_Update(t *testing.T) {
	// (test setup and steps)
}

func TestUnitDataLossPreventionPolicyResource_Validate_Create(t *testing.T) {
	// (test setup and steps)
}

func TestAccDataLossPreventionPolicyResource_Validate_Create(t *testing.T) {
	// (test setup and steps)
}
```

## Fix

Use Go doc comments to document the intent of each test. Extract very large test setups/configurations into helpers/constants for clarity.

```go
// TestUnitDataLossPreventionPolicyResource_Validate_Update verifies updating a data loss prevention policy.
func TestUnitDataLossPreventionPolicyResource_Validate_Update(t *testing.T) {
	// Test logic...
}

// TestUnitDataLossPreventionPolicyResource_Validate_Create verifies creating a data loss prevention policy.
func TestUnitDataLossPreventionPolicyResource_Validate_Create(t *testing.T) {
	// Test logic...
}

// TestAccDataLossPreventionPolicyResource_Validate_Create is acceptance test for creating a data loss prevention policy.
// (this test is currently skipped due to API behavior)
func TestAccDataLossPreventionPolicyResource_Validate_Create(t *testing.T) {
	// Test logic...
}
```
