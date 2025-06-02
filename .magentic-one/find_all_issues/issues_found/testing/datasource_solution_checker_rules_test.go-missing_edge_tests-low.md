# Title

Missing Test Coverage for Edge Attribute Values

##

/workspaces/terraform-provider-power-platform/internal/services/solution_checker_rules/datasource_solution_checker_rules_test.go

## Problem

The tests only cover nominal/regular attribute values and do not exercise edge or boundary values (for example: empty strings except for `how_to_fix`, non-string types, unusually large numbers for `severity`, or unrecognized values).

## Impact

**Low severity** – Not exploring edge cases in test data leaves room for subtle bugs and unhandled conditions in the provider logic, reducing reliability in rare or unusual production scenarios.

## Location

All attribute assertions in both test functions.

## Code Issue

```go
// data.powerplatform_solution_checker_rules.test", "rules.0.how_to_fix", ""  // only attribute where edge is tested
// Other attributes are always given reasonable expected values, never empty/unusual/bad data
```

## Fix

Add test cases for boundary/edge conditions. Example ideas:
- Rules with missing or empty fields
- Extremely long descriptions
- Values outside normal expected range (e.g. severity = -1 or 100)
- Bad URL for guidance_url

```go
resource.TestCheckResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.summary", "")
resource.TestCheckResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.severity", "-1")
resource.TestCheckResourceAttr("data.powerplatform_solution_checker_rules.test", "rules.0.guidance_url", "not-a-url")
```
