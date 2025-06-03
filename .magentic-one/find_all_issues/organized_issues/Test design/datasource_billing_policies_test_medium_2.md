# Title

Dynamic Test Data Generation in Acceptance Test Reduces Reproducibility

## 

/workspaces/terraform-provider-power-platform/internal/services/licensing/datasource_billing_policies_test.go

## Problem

In the `TestAccBillingPoliciesDataSource_Validate_Read` function, input data for resource names and values is generated dynamically at runtime using calls like `mocks.TestName()`, `strconv.Itoa(rand.IntN(9999))`, and `strings.ReplaceAll(mocks.TestName(), "_", "")`. This causes every test run to generate slightly different resource names, potentially making it difficult to reproduce a test or diagnose failures by copying the same test configuration.

## Impact

Severity: **medium**

Dynamic test data makes tests less reproducible because the actual resources used in each run are different. Diagnosing or re-running failed tests using logs or outputs from previous runs is more difficult. In rare cases, if resources are not properly cleaned up, or naming collisions are possible, this can even cause test flakes.

## Location

First test function, lines ~19–53 (test config code):

## Code Issue

```go
resource "azapi_resource" "rg_example" {
	type      = "Microsoft.Resources/resourceGroups@2021-04-01"
	location  = "East US"
	name      = "power-platform-billing-` + mocks.TestName() + strconv.Itoa(rand.IntN(9999)) + `"
}

resource "powerplatform_billing_policy" "pay_as_you_go" {
	name     = "` + strings.ReplaceAll(mocks.TestName(), "_", "") + `"
	location = "unitedstates"
	status   = "Enabled"
	...
}
```

## Fix

For reproducibility and easier diagnostics, prefer using deterministic values, or limit dynamic generation to cases where global uniqueness is required (and log generated values for debugging). A better pattern is to assign specific, static names for all test resources unless there's a strong reason for randomization:

```go
resource "azapi_resource" "rg_example" {
	type      = "Microsoft.Resources/resourceGroups@2021-04-01"
	location  = "East US"
	name      = "power-platform-billing-test"
}

resource "powerplatform_billing_policy" "pay_as_you_go" {
	name     = "payasyougotestpolicy"
	location = "unitedstates"
	status   = "Enabled"
	...
}
```
If isolation is required, consider generating a seed based from the `t.Name()` string and logging the random choices so you can match failures to specific input. If cleanup is a risk, make sure random resource names are properly destroyed after tests.
