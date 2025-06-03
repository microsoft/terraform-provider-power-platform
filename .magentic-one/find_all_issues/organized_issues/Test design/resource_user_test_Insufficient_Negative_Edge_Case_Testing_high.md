# Title

Insufficient Negative/Edge Case Testing

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user_test.go

## Problem

The tests focus primarily on positive scenarios (creation, update, attribute checks). There is no structured negative or edge case testing for the resource, such as invalid IDs, missing required fields, or backend error simulation.

## Impact

Potential regressions or missed bugs in error handling and validation logic. Low test confidence in real-world edge and failure cases. Severity: high.

## Location

All test bodies—no negative resource steps.

## Code Issue

```go
// Only positive cases exist for resource creation/update.
resource.TestCheckResourceAttr(...)
// No scenarios for invalid configs, missing required fields, or API error conditions.
```

## Fix

Add resource test steps (or dedicated test functions) for invalid configs (e.g., missing `environment_id`, invalid `aad_id`), simulating API failure responses (non-200 codes), and assert the right error or panic occurs.

```go
t.Run("invalid config", func(t *testing.T) {
    resource.TestStep{
        Config: `
        resource "powerplatform_user" "fail_user" {
            // Missing required environment_id
            aad_id = "some"
        }`,
        ExpectError: regexp.MustCompile("environment_id is required"),
    }
})
```
