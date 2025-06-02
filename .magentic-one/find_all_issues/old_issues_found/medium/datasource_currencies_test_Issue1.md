# Title

Improper String Validation for Configuration Attributes

##

/workspaces/terraform-provider-power-platform/internal/services/currencies/datasource_currencies_test.go

## Problem

The `regexp.MustCompile(helpers.StringRegex)` is used to perform string validation in multiple test checks, but `helpers.StringRegex` is not visible in this file. It is critical to confirm that `helpers.StringRegex` is properly verified and specific enough. Otherwise, overly permissive validation could lead to false positives or undetected issues in the tests.

## Impact

If the regex defined in `helpers.StringRegex` is overly permissive, it may result in incorrect test passes, reducing the reliability of the test suite and leading to issues in further development and deployment stages.

Severity: Medium

## Location

`TestAccCurrenciesDataSource_Validate_Read` and `TestUnitCurrenciesDataSource_Validate_Read` test cases.

## Code Issue

```go
resource.TestMatchResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "currencies.0.id", regexp.MustCompile(helpers.StringRegex)),
resource.TestMatchResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "currencies.0.name", regexp.MustCompile(helpers.StringRegex)),
resource.TestMatchResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "currencies.0.display_name", regexp.MustCompile(helpers.StringRegex)),
resource.TestMatchResourceAttr("data.powerplatform_currencies.all_currencies_for_unitedstates", "currencies.0.locale_id", regexp.MustCompile(helpers.StringRegex)),
```

## Fix

Ensure that `helpers.StringRegex` is properly designed to validate specific conditions for the attributes being tested, such as proper format for names, IDs, or locale identifiers. You may also add comments documenting the constraints enforced by this regex.

```go
// Verify that helpers.StringRegex ensures specific constraints for attributes
func TestStringRegex() {
    assert.Regexp(t, helpers.StringRegex, "ExpectedValidString")
    assert.NotRegexp(t, helpers.StringRegex, "UnexpectedInvalidString")
}
```
