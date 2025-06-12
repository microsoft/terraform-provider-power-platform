# Title

Redundant/Repeated Code in Test Definitions

##

/workspaces/terraform-provider-power-platform/internal/services/solution_checker_rules/datasource_solution_checker_rules_test.go

## Problem

There is repeated/redundant configuration in both `TestAccSolutionCheckerRulesDataSource_Validate_Read` and `TestUnitSolutionCheckerRulesDataSource_Validate_Read` regarding data source declaration and the use of TestCheckResourceAttr or TestMatchResourceAttr for similar attributes. While details differ, the pattern is duplicated, making maintainability harder.

## Impact

**Low severity** – Repetition makes the suite harder to maintain; missing DRY (Don't Repeat Yourself) principles can lead to drift between test intent and actual coverage, divergence, and extra effort on updates.

## Location

Affects:
- Both test functions
- Steps and aggregate check setup in each

## Code Issue

```go
// Both test functions re-declare very similar Config sections
// Both use similar ComposeAggregateTestCheckFunc blocks referencing nearly the exact same attributes
```

## Fix

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
