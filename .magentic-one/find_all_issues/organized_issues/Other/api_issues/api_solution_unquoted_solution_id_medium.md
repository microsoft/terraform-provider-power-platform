# Title

Incorrect HTTP Status Codes When Fetching an Individual Solution by ID

##

internal/services/solution/api_solution.go

## Problem

In `GetSolutionById`, the function constructs a filter for solution ID, but the API filter string does not quote the string value for `solutionId`. According to OData and typical REST conventions, string values should be quoted, e.g., `solutionid eq 'xxx'` instead of `solutionid eq xxx`, otherwise the server may reject the query.

## Impact

Severity: **medium**. Passing the solutionId as an unquoted value can cause inconsistent or failed API requests and make the code fragile if OData or the backend provider becomes more strict in the future.

## Location

```go
values.Add("$filter", fmt.Sprintf("solutionid eq %s", solutionId))
```

## Code Issue

```go
values.Add("$filter", fmt.Sprintf("solutionid eq %s", solutionId))
```

## Fix

The correct filter should quote the ID as a string:

```go
values.Add("$filter", fmt.Sprintf("solutionid eq '%s'", solutionId))
```
