# Use of Inline HCL Instead of Test Fixtures or Helper Functions

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy_test.go

## Problem

The HCL resource configurations are constructed via multi-line inline strings within test steps. This makes them hard to manage and error-prone for complex test data.

## Impact

- **Maintainability**: Duplication across different steps and tests.
- **Readability**: Noise from string concatenation logic and indentation.
- **Extensibility**: Difficult to add new scenarios with only minimal differences.

**Severity: Low**

## Location

All `Config: ` assignments in test steps.

## Code Issue

```go
Config: `
resource "powerplatform_billing_policy" "pay_as_you_go" {
    name     = "payAsYouGoBillingPolicyExample"
    # ...
}`
```

## Fix

Refactor resource configurations into reusable functions or load from external files.

```go
func billingPolicyConfig(name, location, status, rg, sub string) string {
    return fmt.Sprintf(`
resource "powerplatform_billing_policy" "pay_as_you_go" {
    name     = "%s"
    location = "%s"
    status   = "%s"
    billing_instrument = {
      resource_group  = "%s"
      subscription_id = "%s"
    }
}`, name, location, status, rg, sub)
}
```

Or, use separate files with `os.ReadFile("testdata/policy_create.hcl")`.
