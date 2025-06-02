# Title

Test Coverage: No Negative or Edge Case Tests Except for One

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/datasource_environment_settings_test.go

## Problem

Most tests check only the successful (happy path) state of the data source, except a single test that checks for a missing Dataverse (`Validate_No_Dataverse`). There are no tests for other error conditions, such as malformed responses, permission errors, or partial returns.

## Impact

Insufficient negative and edge case test coverage can lead to undetected bugs or regressions in error handling. Risk is medium when dealing with cloud resources and APIs which may respond in various failure modes in real-world usage.

## Location

General pattern across the file. Specifically, `TestAccTestEnvironmentSettingsDataSource_Validate_Read` and `TestUnitTestEnvironmentSettingsDataSource_Validate_Read` are happy-path only.

## Code Issue

```go
resource.Test(t, resource.TestCase{
    ...
    Steps: []resource.TestStep{
        {
            Config: `...`,

            Check: resource.ComposeAggregateTestCheckFunc(
                // Only checks for positive conditions, no ExpectError or error cases
            ),
        },
    },
})
```

## Fix

Add negative/edge-case tests, for example:
- Missing/invalid environment_id.
- Auth failures in the mock responder.
- Unexpected/malformed JSON.
- Partial data (e.g., certain nested attributes missing).

You can use `ExpectError`, `ExpectNonEmptyPlan`, or custom checks to validate error conditions. E.g.:

```go
// When environment_id is invalid
{
    Config: `
        data "powerplatform_environment_settings" "settings" {
            environment_id = "invalid-guid"
        }
    `,
    ExpectError: regexp.MustCompile("invalid environment_id"),
},
```

---

This file will be saved as:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/testing/datasource_environment_settings_test.go_limited_negative_tests_medium.md`
