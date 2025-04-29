# Title

Duplicate Resource Attribute Check

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/datasource_dlp_policy_test.go

## Problem

The following attribute is being checked multiple consecutive times in one of the test configurations, which results in redundancy:

```go
resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.business_connectors.#", "0"),
resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.business_connectors.#", "0"),
```

## Impact

This redundancy unnecessarily bloats the test file and could confuse readers or developers maintaining the codebase. Severity is low for this issue but cleaning it up improves code clarity.

## Location

Line 37-38 in the file `/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/datasource_dlp_policy_test.go`.

## Code Issue

```go
resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.business_connectors.#", "0"),
resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.business_connectors.#", "0"),
```

## Fix

Removing the duplicate line results in cleaner tests:

```go
resource.TestCheckResourceAttr("data.powerplatform_data_loss_prevention_policies.all", "policies.0.business_connectors.#", "0"),
```