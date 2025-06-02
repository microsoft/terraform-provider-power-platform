# No Assertion Comments or Grouping in resource.ComposeAggregateTestCheckFunc

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy_test.go

## Problem

Assertions are listed one after another with minimal distinction, making it difficult to see which assertion validates what scenario in the resource.

## Impact

- **Readability**: Difficult to trace which assertion relates to which resource.
- **Debuggability**: Harder to track test failures to actual intent.

**Severity: Low**

## Location

All `Check: resource.ComposeAggregateTestCheckFunc(...)` invocations.

## Code Issue

```go
Check: resource.ComposeAggregateTestCheckFunc(
    resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", "00000000-0000-0000-0000-000000000001"),
    resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "name", "payAsYouGoBillingPolicyExample"),
    // ... more attributes
),
```

## Fix

Group assertions by logical sections and add comments, or, better yet, split into mini helper check functions for better intent expression.

```go
Check: resource.ComposeAggregateTestCheckFunc(
    // ID Checks
    resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "id", "expected-id"),
    // Core Attribute Checks
    resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "name", "expected-name"),
    resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "status", "Enabled"),
    // Billing Instrument Checks
    resource.TestCheckResourceAttr("powerplatform_billing_policy.pay_as_you_go", "billing_instrument.resource_group", "expected-group"),
),
```
