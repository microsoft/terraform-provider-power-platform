# Issue: Test Quality—No Negative or Edge Case Checks

## 
/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/resource_enterprise_policy_test.go

## Problem

There is only a positive ("happy path") test and no test steps to cover invalid inputs, API failures, or edge cases.

## Impact

Missing checks for errors or unexpected input leads to less robust tests, so regressions or failure modes may go undetected. Severity: medium.

## Location

The `resource.TestCase` definition, lines 52-94. Only one `resource.TestStep` is present, all for success.

## Code Issue

```go
Steps: []resource.TestStep{
    {
        Config: `...`,
        Check: resource.ComposeAggregateTestCheckFunc(
            ...
        ),
    },
},
```

## Fix

Add further `resource.TestStep` entries for negative cases or edge conditions. For example, testing with an invalid `policy_type` or simulating HTTP failures.

```go
Steps: []resource.TestStep{
    // Existing positive test step...
    {
        Config: `...valid config...`,
        Check: resource.ComposeAggregateTestCheckFunc(
            // ...
        ),
    },
    // Negative case: invalid policy_type
    {
        Config: `
            resource "powerplatform_environment" "dev" {
                display_name = "x"
                description = "desc"
                cadence     = "Frequent"
                location    = "europe"
                environment_type = "Sandbox"
            }
            resource "powerplatform_enterprise_policy" "bad_policy" {
                environment_id = powerplatform_environment.dev.id
                system_id      = "/regions/europe/providers/Microsoft.PowerPlatform/enterprisePolicies/1234"
                policy_type    = "INVALID_POLICY"
            }
        `,
        ExpectError: regexp.MustCompile("expected policy_type to be one of"),
    },
}
```
