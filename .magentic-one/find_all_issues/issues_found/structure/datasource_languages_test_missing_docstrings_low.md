# Lack of Documentation Comments for Test Functions

##

/workspaces/terraform-provider-power-platform/internal/services/languages/datasource_languages_test.go

## Problem

Test functions lack documentation comments, which hurts IDE/generated docs and makes grasping the test’s purpose harder.

## Impact

Readability/maintainability, low severity.

## Location

All test function declarations.

## Code Issue

```go
func TestAccLanguagesDataSource_Validate_Read(t *testing.T) {
...
func TestUnitLanguagesDataSource_Validate_Read(t *testing.T) {
```

## Fix

Add docstrings explaining what the test validates:

```go
// TestAccLanguagesDataSource_Validate_Read validates the data source against the live provider with region "unitedstates".
func TestAccLanguagesDataSource_Validate_Read(t *testing.T) { ... }
```
