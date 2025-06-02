# Test Functions are Very Large & Duplicative

##

/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/resource_environment_group_rule_set_test.go

## Problem

Individual test functions, especially for acceptance and CRUD flows, are very long and difficult to read and maintain. There is considerable repeated configuration and verification code — for example, the Terraform resource configuration HCL is almost always repeated, and many `TestCheckResourceAttr` lines are duplicated across test steps. This increases the risk of errors/omissions and makes future additions or failure analysis harder.

## Impact

Severity: Low

Excessive duplication and long functions reduce maintainability, increase merge conflicts, and slow down test updates.

## Location

All test function bodies, e.g., both acceptance and unit CRUD tests.

## Fix

Refactor repeated code into helper constants/functions.

Example:

```go
const validRuleSetConfig = `
  resource "powerplatform_environment_group_rule_set" "example_group_rule_set" {
    // ...
  }
`

func checkAllAttributes() resource.TestCheckFunc {
  return resource.ComposeAggregateTestCheckFunc(
    resource.TestCheckResourceAttr("...", "...", "..."),
    // ...
  )
}

// Then use:
{
  Config: validRuleSetConfig,
  Check:  checkAllAttributes(),
}
```

This will reduce code length, help isolate test updates, and clarify intent.
