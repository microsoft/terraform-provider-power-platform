# Issue 5: Inefficiency in Result Query Building in Filters

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/odata.go

## Problem

Code in `buildExpandFilterQueryPart` repeatedly uses `strings.Join([]string{...}, "")` and string concatenations in multiple if-blocks, decreasing readability and maintainability.

## Impact

Makes the logic less straightforward for maintainers, affecting readability (severity: **low**).

## Location

Lines 57-79

## Code Issue

```go
resultQuery = strings.Join([]string{resultQuery, *selectString}, "")
// and similar for filterString, orderByString, and so on.
```

## Fix

Just use direct concatenation and compound assignments.

```go
resultQuery += *selectString
```
