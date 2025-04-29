# Title

Redundant and Repetitive Logic in Test Configuration

## Path

`/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy_test.go`

## Problem

There is a significant amount of repetitive configuration and resource setup logic across various test cases (`Config`, `resource "powerplatform_data_loss_prevention_policy" "my_policy"`). Instead of reusing and parameterizing this logic, the same code is written multiple times with slight modifications.

## Impact

This issue reduces code maintainability, increases the likelihood of errors when updating tests, and makes extending the test suite cumbersome. The severity is **medium**, as it directly affects development efficiency and code quality.

## Location

Found in multiple test cases, such as:
- `TestUnitDataLossPreventionPolicyResource_Validate_Update`
- `TestUnitDataLossPreventionPolicyResource_Validate_Create`
- `TestAccDataLossPreventionPolicyResource_Validate_Create`

## Code Issue

```go
Config: `
resource "powerplatform_data_loss_prevention_policy" "my_policy" {
    display_name = "Block All Policy_1"
    default_connectors_classification = "Blocked"
    environment_type = "AllEnvironments"
    environments = []

    business_connectors = []
    non_business_connectors = []
    blocked_connectors = []

    custom_connectors_patterns = toset([
        {
            order = 1
            host_url_pattern = "*"
            data_group = "Ignore"
        }
    ])
}`
```

## Fix

Create test helper functions or reusable constants to centralize configuration setup. This will reduce redundancy and improve code maintainability.

### Example Fix

```go
// Define a reusable helper function for creating test configuration
func createTestPolicyConfig(displayName, classification, envType string) string {
    return fmt.Sprintf(`
        resource "powerplatform_data_loss_prevention_policy" "my_policy" {
            display_name = "%s"
            default_connectors_classification = "%s"
            environment_type = "%s"
            environments = []

            business_connectors = []
            non_business_connectors = []
            blocked_connectors = []

            custom_connectors_patterns = toset([
                {
                    order = 1
                    host_url_pattern = "*"
                    data_group = "Ignore"
                }
            ])
        }`, displayName, classification, envType)
}

// Usage in the test case
Config: createTestPolicyConfig("Block All Policy_1", "Blocked", "AllEnvironments"),
```