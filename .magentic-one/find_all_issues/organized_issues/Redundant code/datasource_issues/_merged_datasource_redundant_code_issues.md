# Data Source Redundant Code Issues

This document consolidates all redundant code issues found in data source components of the Terraform Power Platform provider.

## ISSUE 1

**File:** `/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/datasource_dlp_policy_test.go`

### Problem

In the `TestUnitDlpPolicyDataSource_Validate_Read`, the test checks for `"policies.0.business_connectors.#"` twice in succession. This is redundant and provides no additional value, making the test harder to maintain.

### Impact

This decreases maintainability and clarity, especially for someone maintaining or reviewing these tests. It can be confusing and suggests a possible copy-paste or logic error. Severity: low.

### Location

Lines 38-40 (approximate)

### Code Issue

```go
resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.business_connectors.#", "0"),
resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.business_connectors.#", "0"),
```

### Fix

Remove the duplicate assertion for `"policies.0.business_connectors.#"`:

```go
resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.business_connectors.#", "0"),
```

## ISSUE 2

**File:** `/workspaces/terraform-provider-power-platform/internal/services/application/datasource_environment_application_packages_test.go`

### Problem

In `TestUnitEnvironmentApplicationPackagesDataSource_Validate_No_Dataverse()`, the test step sets both an `ExpectError` and an empty `Check` via `resource.ComposeAggregateTestCheckFunc()`. Specifying an empty check function is redundant when only error validation is needed.

### Impact

Severity: Low

This adds a minor amount of unnecessary cognitive overhead to the test; the purpose of the test step would be clearer if non-applicable code was omitted.

### Location

```go
ExpectError: regexp.MustCompile("No Dataverse exists in environment"),
Check:       resource.ComposeAggregateTestCheckFunc(),
```

### Fix

You may remove the `Check` entirely when the test's only assertion is that an error occurs:

```go
ExpectError: regexp.MustCompile("No Dataverse exists in environment"),
// Remove `Check`
```

Or simply avoid the empty function where no checks are required.

## ISSUE 3

**File:** `/workspaces/terraform-provider-power-platform/internal/services/application/datasource_environment_application_packages_test.go`

### Problem

The same set of HTTPMock responders is re-registered in multiple test functions, introducing repetitive and redundant code. This makes maintaining the tests unnecessarily difficult and increases the risk of inconsistencies.

### Impact

Severity: Medium

Duplicated code is harder to maintain and synchronize. Bugs or updates in mock setup require multiple coordinated changes, raising risk for inconsistent test results.

### Location

All three main unit test functions, e.g.:

```go
httpmock.Activate()
defer httpmock.DeactivateAndReset()

mocks.ActivateEnvironmentHttpMocks()

httpmock.RegisterResponder("GET", `https://api.powerplatform.com/...`,
    func(req *http.Request) (*http.Response, error) {
        // ...
    })
```

### Fix

Extract HTTPMock responder setup into a helper function that is reused in each test function.

```go
func setupEnvironmentApplicationPackagesMocks(testVariant string) {
    httpmock.RegisterResponder(..., // use testVariant for test file selection
        func(req *http.Request) (*http.Response, error) {
            // ...
        })
    // Add the other responders
}

// Call this in the tests:
setupEnvironmentApplicationPackagesMocks("Validate_Read")
```

## ISSUE 4

**File:** `/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query.go`

### Problem

In the `Read` function, after setting `state.Output`, you set the state again and append the diagnostics. However, you immediately check for errors and return, which is redundant because diagnostics are already appended, and in Terraform Plugin SDKs, the diagnostics should generally be returned if present, but it's best practice to append diagnostics in-place and return only if a critical condition is hit (not just always immediately after).

### Impact

This statement is not a bug, but is minorly misleading in idiomatic Go/Terraform Plugin SDK code, potentially confusing future maintainers or causing over-defensive exit flow. It is of low severity.

### Location

At the end of the `Read` method:

### Code Issue

```go
 diags := resp.State.Set(ctx, &state)
 resp.Diagnostics.Append(diags...)
 if resp.Diagnostics.HasError() {
  return
 }
```

### Fix

Just append the diagnostics. Returning from `Read` after appending diagnostics that might not actually be errors is an overly defensive code pattern. Instead, you may append and only return early in more critical flows. For clarity, this block could be trimmed as:

```go
 resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
```

If you want to retain the `HasError` early return, consider adding a comment and restructuring for clarity. If this is code style for all your providers, it is not strictly wrong, but be aware that the plugin will handle the diagnostics array, and returning is usually needed only for control flow exit, not always after every set.

## ISSUE 5

**File:** `/workspaces/terraform-provider-power-platform/internal/services/solution_checker_rules/datasource_solution_checker_rules_test.go`

### Problem

There is repeated/redundant configuration in both `TestAccSolutionCheckerRulesDataSource_Validate_Read` and `TestUnitSolutionCheckerRulesDataSource_Validate_Read` regarding data source declaration and the use of TestCheckResourceAttr or TestMatchResourceAttr for similar attributes. While details differ, the pattern is duplicated, making maintainability harder.

### Impact

**Low severity** – Repetition makes the suite harder to maintain; missing DRY (Don't Repeat Yourself) principles can lead to drift between test intent and actual coverage, divergence, and extra effort on updates.

### Location

Affects:

- Both test functions
- Steps and aggregate check setup in each

### Code Issue

```go
// Both test functions re-declare very similar Config sections
// Both use similar ComposeAggregateTestCheckFunc blocks referencing nearly the exact same attributes
```

### Fix

Use shared helper functions to construct repeated configuration snippets and to assemble common attribute checking logic.

```go
func testCheckerRuleConfig(environmentID string) string {
    return fmt.Sprintf(`
    data "powerplatform_solution_checker_rules" "test" {
        environment_id = "%s"
    }
    `, environmentID)
}

func testCheckerRuleChecks() resource.TestCheckFunc {
    return resource.ComposeAggregateTestCheckFunc(
        resource.TestCheckResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.code", "meta-remove-dup-reg"),
        // ...repeat as needed
    )
}
```

---

## To finish the task you have to

1. Run linter and fix any issues
2. Run UnitTest and fix any of failing ones
3. Generate docs
4. Run Changie

## Changie Instructions

Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```

Where:

- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number

Apply this fix to the whole codebase
