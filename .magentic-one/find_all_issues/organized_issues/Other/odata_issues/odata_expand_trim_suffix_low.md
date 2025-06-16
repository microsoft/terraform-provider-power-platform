# Issue 3: Lack of Helper Function for Repeated OData Part Suffix Handling

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/odata.go

## Problem

`buildExpandODataQueryPart` includes code to join and trim the result of sub-expands, e.g., `strings.Join(expandQueryStrings, ",")` and `strings.TrimSuffix(result, ",")`, but this trimming is unnecessary as `Join` does not append a trailing separator.

## Impact

Unnecessary/no-op operation reduces clarity. Severity is **low**.

## Location

```go
result := strings.Join(expandQueryStrings, ",")
result = "$expand=" + strings.TrimSuffix(result, ",")
return &result
```

## Fix

Drop the `TrimSuffix` call.

```go
result := "$expand=" + strings.Join(expandQueryStrings, ",")
return &result
```
