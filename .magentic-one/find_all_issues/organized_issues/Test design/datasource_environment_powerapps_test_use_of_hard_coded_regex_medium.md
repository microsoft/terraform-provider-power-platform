# Use of Hard-Coded Regular Expressions Instead of Constants

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/datasource_environment_powerapps_test.go

## Problem

The following regular expression is hard-coded directly into the test function:

```go
resource.TestMatchResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.0.created_time", regexp.MustCompile(`^\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+([+-][0-2]\d:[0-5]\d|Z)$`))
```

While some regex patterns are referenced from helpers (such as `helpers.GuidRegex`, etc.), the `created_time` regex is written in-line rather than defined as a constant or in a shared location. This makes the code less maintainable and harder to reuse or update in the future.

## Impact

- **Maintainability**: Medium – Using inline regular expressions reduces test maintainability and increases the risk of duplication or inconsistencies if the pattern is ever used elsewhere.
- **Readability**: Slightly reduces readability due to dense regex inline with test logic.

## Location

```go
resource.TestMatchResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.0.created_time", regexp.MustCompile(`^\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+([+-][0-2]\d:[0-5]\d|Z)$`)),
```

## Code Issue

```go
resource.TestMatchResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.0.created_time", regexp.MustCompile(`^\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+([+-][0-2]\d:[0-5]\d|Z)$`)),
```

## Fix

Define a constant in a shared helpers location for the datetime format and reference it here:

```go
// helpers/regex.go or similar
package helpers

const DateTimeRegex = `^\d{4}-[01]\d-[0-3]\dT[0-2]\d:[0-5]\d:[0-5]\d\.\d+([+-][0-2]\d:[0-5]\d|Z)$`

// In your test file
resource.TestMatchResourceAttr("data.powerplatform_environment_powerapps.all", "powerapps.0.created_time", regexp.MustCompile(helpers.DateTimeRegex)),
```
