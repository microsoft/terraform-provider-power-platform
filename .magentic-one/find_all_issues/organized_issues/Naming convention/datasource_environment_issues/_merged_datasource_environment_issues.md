# Datasource Environment Issues - Merged Issues

## ISSUE 1

# Inconsistent Use of Test Name Helper Utility

##

/workspaces/terraform-provider-power-platform/internal/services/application/datasource_environment_application_packages_test.go

## Problem

In one acceptance test the display name is set using `mocks.TestName()`, in others it's hardcoded to "displayname". This leads to inconsistent naming and could impact test deduplication, parallelization, and diagnostics.

## Impact

Severity: Low

Inconsistencies may lead to confusion during debugging, make it difficult to search logs, and may affect the behavior of tests if uniqueness is required.

## Location

For example:

```go
display_name      = "` + mocks.TestName() + `"
```
vs
```go
display_name                              = "displayname"
```

## Fix

Adopt the same approach for all tests. Preferably use the utility for unique, informative naming:

```go
display_name = "` + mocks.TestName() + `"
```


---

# To finish the task you have to 
1. Run linter and fix any issues 
2. Run UnitTest and fix any of failing ones
3. Generate docs 
4. Run Changie

# Changie Instructions
Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```
Where:
- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
