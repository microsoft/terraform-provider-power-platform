# Title

Hardcoded Attributes in Tests

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy_test.go

## Problem

The resource configuration in tests includes hardcoded attributes such as `location`, `status`, and `resource_group`. This approach reduces test flexibility and makes future updates difficult.

## Impact

Hardcoded values directly impact test scalability, reusability, and reduced adaptability for different environments. Severity: Medium.

## Location

For example, `location = "europe"` found in the configuration block starting on line 121.

## Code Issue

```go
resource "powerplatform_billing_policy" "pay_as_you_go" {
    location = "europe"
}
```

## Fix

Utilize variables or test-specific configuration files to define values like `location`. Incorporate dynamic values from mocks or environment-specific placeholders.

```go
resource "powerplatform_billing_policy" "pay_as_you_go" {
    location = mock.Location()
}
```

This ensures the tests can adapt to varying environments dynamically and can be maintained efficiently.
