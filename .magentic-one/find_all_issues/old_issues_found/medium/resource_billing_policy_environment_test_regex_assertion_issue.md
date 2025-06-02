# Title

Improper assertions using regular expressions in test checks

## Path

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy_environment_test.go

## Problem

In the test cases, assertions validate resource attributes using regular expressions (e.g., `regexp.MustCompile(helpers.GuidRegex)`). Although these validations check if the attributes resemble GUIDs or specific patterns, they do not verify whether the value actually corresponds to a valid resource or matches the expected state within the Terraform environment.

## Impact

Assertions based solely on regular expressions can give a false sense of correctness. A value matching the GUID format does not guarantee that it represents a valid or correctly provisioned resource, leading to potential false positive test results. Severity: **Medium**

## Location

Below is an example where regular expression assertions are used improperly:

### Code Issue

```go
resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", regexp.MustCompile(helpers.GuidRegex)),
resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.0", regexp.MustCompile(helpers.GuidRegex)),
resource.TestMatchResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.1", regexp.MustCompile(helpers.GuidRegex)),
```

## Fix

Enhance the assertions to verify the actual values or state of the resources instead of depending solely on regular expressions. For example, use the Terraform output or mock responses to validate correctness.

### Code Example

```go

// Replace TestMatchResourceAttr with TestCheckResourceAttr using actual expected values:
resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "billing_policy_id", "expected-billing-policy-id"),
resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.0", "expected-environment-id-1"),
resource.TestCheckResourceAttr("powerplatform_billing_policy_environment.pay_as_you_go_policy_envs", "environments.1", "expected-environment-id-2"),
```

This ensures that the tests verify concrete expected values, making them more robust and meaningful.
