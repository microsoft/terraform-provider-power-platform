# Title

Unnecessary Duplicate Attribute Check in Test Case

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/datasource_dlp_policy_test.go

## Problem

In the `TestUnitDlpPolicyDataSource_Validate_Read`, the test checks for `"policies.0.business_connectors.#"` twice in succession. This is redundant and provides no additional value, making the test harder to maintain.

## Impact

This decreases maintainability and clarity, especially for someone maintaining or reviewing these tests. It can be confusing and suggests a possible copy-paste or logic error. Severity: low.

## Location

Lines 38-40 (approximate)

## Code Issue

```go
resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.business_connectors.#", "0"),
resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.business_connectors.#", "0"),
```

## Fix

Remove the duplicate assertion for `"policies.0.business_connectors.#"`:

```go
resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.business_connectors.#", "0"),
```
