# Title

Hardcoded resource identifiers in unit tests

## Path

/workspaces/terraform-provider-power-platform/internal/services/licensing/resource_billing_policy_environment_test.go

## Problem

Unit test configurations hardcode resource identifiers such as billing policy IDs and environment IDs in API requests and assertions. For example, strings like `"00000000-0000-0000-0000-000000000000"` are used directly, which does not dynamically adjust or vary across test runs.

## Impact

Hardcoding resource identifiers restricts the flexibility and scalability of unit tests. Any future changes in resource identifiers or conventions might require widespread manual updates to all instances. Additionally, these tests might inadvertently test hardcoded values rather than the dynamic behavior expected for actual inputs. Severity: **High**

## Location

Below is an example of a code segment with hardcoded identifiers:

### Code Issue

```go
billing_policy_id = "00000000-0000-0000-0000-000000000000"
```

## Fix

Refactor the test code to dynamically generate or mock resource identifiers instead of hardcoding them. This can be achieved using helper functions or constants defined specifically for testing.

### Code Example

```go

// Define constants or helper functions to generate resource identifiers dynamically.
const TestBillingPolicyID = "test-billing-policy-id"

// Replace hardcoded IDs with the defined constants:
billing_policy_id = TestBillingPolicyID

// Alternatively, create a mock generator for IDs:
func GenerateTestResourceID(resourceType string) string {
    return fmt.Sprintf("test-%s-%s", resourceType, uuid.New().String())
}
billing_policy_id = GenerateTestResourceID("billing-policy")

```

Using generated or configurable identifiers instead of hardcoded values ensures that the tests remain adaptable across environments and refactor-friendly.