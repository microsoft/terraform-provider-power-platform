# Redundant Regular Expressions for Simple Attribute Checks

##

/workspaces/terraform-provider-power-platform/internal/services/languages/datasource_languages_test.go

## Problem

The tests use regular expressions to validate attributes that could be checked with more straightforward assertions, reducing clarity and potentially hiding logic errors.

## Impact

This can make the tests harder to read and debug, a low severity readability and maintainability issue.

## Location

```go
resource.TestMatchResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.#", regexp.MustCompile(`^[1-9]\d*$`)),
```

## Code Issue

```go
resource.TestMatchResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.#", regexp.MustCompile(`^[1-9]\d*$`)),
```

## Fix

If you expect a specific count and value, use `TestCheckResourceAttr`:

```go
resource.TestCheckResourceAttr("data.powerplatform_languages.all_languages_for_unitedstates", "languages.#", "45"),
```

For string attributes, consider equality if you know the expected value; use regex only if flexibility is required.
