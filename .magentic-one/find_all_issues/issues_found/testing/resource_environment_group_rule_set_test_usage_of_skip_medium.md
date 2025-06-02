# Usage of t.Skip hides incomplete acceptance test

##

/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/resource_environment_group_rule_set_test.go

## Problem

In the acceptance tests (`TestAccEnvironmentGroupRuleSetResource_Validate_Create` and `TestAccEnvironmentGroupRuleSetResource_Validate_Update`), the first line is `t.Skip("creating rule sets with SP is NOT yet supported")`. This means these tests are never actually executed. While it's valid to skip unsupported flows, the presence of these test cases can give a false impression of complete test coverage for the acceptance layer. There should be a clear indication in the codebase or in developer documentation that these are intentionally not run as opposed to being forgotten or unimplemented.

## Impact

Severity: Medium

Having incomplete/disabled acceptance tests could lead teams or reviewers mistakenly assuming the provider is fully tested in a real environment, which could lead to regression when the functionality is eventually enabled or used.

## Location

Lines with `t.Skip(...)` in acceptance tests:
- Near line: function TestAccEnvironmentGroupRuleSetResource_Validate_Create
- Near line: function TestAccEnvironmentGroupRuleSetResource_Validate_Update

## Code Issue

```go
func TestAccEnvironmentGroupRuleSetResource_Validate_Create(t *testing.T) {
	t.Skip("creating rule sets with SP is NOT yet supported")
	resource.Test(t, resource.TestCase{
	  ...
	})
}
...
func TestAccEnvironmentGroupRuleSetResource_Validate_Update(t *testing.T) {
	t.Skip("creating rule sets with SP is NOT yet supported")
	resource.Test(t, resource.TestCase{
	  ...
	})
}
```

## Fix

Add a TODO or clearer comment documenting why the tests are skipped, and optionally use Go build tags or t.Helper() with a gating mechanism that makes explicit when the acceptance tests should be re-enabled. 

Alternatively, document in the file header that acceptance tests are placeholders and are intentionally skipped, to avoid confusion by contributors and reviewers.

```go
// TODO: Remove t.Skip() and re-enable test when service principal creation is supported

func TestAccEnvironmentGroupRuleSetResource_Validate_Create(t *testing.T) {
	// Acceptance tests are intentionally skipped until rule set creation is supported via SP 
	t.Skip("creating rule sets with SP is NOT yet supported. This is a placeholder acceptance test.")
	...
}
```
