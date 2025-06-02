# Title

Reliance on Magic Strings and Hardcoded IDs in Test Setup

##

internal/services/data_record/datasource_data_record_test.go

## Problem

Tests frequently use repeated, hardcoded UUIDs and other literal values (e.g., `"00000000-0000-0000-0000-000000000001"`) for environment IDs, entity IDs, and resource names. These "magic strings" recur many times, both in test resource configs and mock HTTP endpoints, increasing the chance of copy-paste errors and making it harder to update or refactor test infrastructure.

## Impact

- Difficult refactoring due to multiple places to update values (low-to-medium impact).
- Higher risk of inconsistencies and accidental mismatches.
- Obfuscates meaning and intent, reducing test readability and maintainability.

## Location

Examples:
- In HCL configuration strings (e.g., `environment_id = "00000000-0000-0000-0000-000000000001"`)
- In mock URLs (e.g., `"https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts?...`)
- In JSON paths, responders, etc.

## Code Issue

```go
// test HCL snippet
environment_id = "00000000-0000-0000-0000-000000000001"
// endpoint registration
httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/contacts?$select=...", ...)
// repeated throughout the file
```

## Fix

Replace repeated magic strings with clearly-named constants at the top level of the test file, or with a function to generate standard/randomized test UUIDs. For example:

```go
const testEnvID = "00000000-0000-0000-0000-000000000001"
// Then use testEnvID everywhere instead of repeated literals
```
Or:
```go
func testEnvID() string { return \"00000000-0000-0000-0000-000000000001\" }
// Or return a randomized value if possible
```
This makes maintaining/generating test configuration and updating test resource IDs easier, and improves code readability and robustness.

Save as a structure/readability/maintainability issue.
