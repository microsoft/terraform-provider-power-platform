# Title

Misuse of trailing commas in slice and map composite literals

##

internal/services/data_record/datasource_data_record_test.go

## Problem

Some composite literals for slices and maps within the Terraform HCL test configuration strings use trailing commas just before closing parentheses of elements, or trailing commas within slice elements. While Go itself may allow trailing commas, the HCL representation being inserted into the string template is syntactically incorrect when later interpolated, and this may cause test failures if the string is parsed.

## Impact

If the string containing HCL/config is parsed, the extra commas may cause parsing errors, leading to failing or flaky tests, or worse, false negatives. This has a medium impact as it may prevent acceptance or unit tests for this provider from running properly.

## Location

Several occurrences within dynamically-generated string literals, e.g.:

- In `expand = [...]` and other similar block assignments.
- Particularly in blocks under:  
  - `TestAccDataRecordDatasource_Validate_Expand_Lookup`  
  - `TestAccDataRecordDatasource_Validate_Expand_Query`  
  - Other tests that use comma after the last item in HCL array/list/map entries.

## Code Issue

For example:
```go
expand = [
    {
        navigation_property = "primarycontactid"
        select              = ["contactid", "firstname", "lastname"],
    },
    {
        navigation_property = "owningbusinessunit"
        select              = ["createdon", "name"],
    }
]
```
and several similar HCL configuration string literals.

## Fix

Remove the trailing commas after the last element in the inline list/map in HCL string literals. The correct string should look like:

```go
expand = [
    {
        navigation_property = "primarycontactid"
        select              = ["contactid", "firstname", "lastname"]
    },
    {
        navigation_property = "owningbusinessunit"
        select              = ["createdon", "name"]
    }
]
```
And similarly remove other incorrect trailing commas from generated HCL in string blocks. This will ensure that the configuration is correctly parsed and interpreted by Terraform and the test frameworks.
