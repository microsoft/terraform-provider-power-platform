# Title

Skipped Acceptance Test Without Tracking or Clear Reason

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/datasource_dlp_policy_test.go

## Problem

The function `TestAccDlpPolicyDataSource_Validate_Read` starts with `t.Skip("Skipping due to inconsistent connectors results")`, effectively disabling the acceptance test. There is no clear tracking, TODO, or documentation regarding a plan to fix or re-enable this important test.

## Impact

Disabling acceptance tests for extended periods leads to decreased test coverage and possible undetected regressions. Lacking clear reference or tracking, it can be forgotten and never re-enabled. Severity: medium.

## Location

Function: `TestAccDlpPolicyDataSource_Validate_Read`

## Code Issue

```go
func TestAccDlpPolicyDataSource_Validate_Read(t *testing.T) {
    t.Skip("Skipping due to inconsistent connectors results")
    ...
}
```

## Fix

Add a `// TODO:` comment referencing a ticket, GitHub issue, or tracking link so this skip is intentional and findable. Plan to enable the test, or provide more details about the issue.

```go
func TestAccDlpPolicyDataSource_Validate_Read(t *testing.T) {
    // TODO: https://github.com/<repo>/issues/<id> - Skipped due to inconsistent connectors results. Re-enable when fixed.
    t.Skip("Skipping due to inconsistent connectors results")
    ...
}
```
