# Issue 2: Repeated String Joining for Single Concatenations

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/odata.go

## Problem

Several `buildOData*Part` functions use `strings.Join([]string{a, b}, "")` for concatenation, which is inefficient for only two strings. Simple string concatenation (e.g., `a + b`) is more readable and performant in Go.

## Impact

Minor performance inefficiency and reduced readability. Severity is **low**.

## Location

Examples at lines: 100, 108, 116, 124, 132, 140, 148, 156

## Code Issue

```go
resultQuery = strings.Join([]string{"savedQuery=", url.QueryEscape(*savedQuery)}, "")
// Repeated elsewhere: "userQuery=", "$filter=", "$orderby=", "$top=", "$apply="
```

## Fix

Replace with simple concatenation.

```go
resultQuery = "savedQuery=" + url.QueryEscape(*savedQuery)
// And similarly for the other places.
```
