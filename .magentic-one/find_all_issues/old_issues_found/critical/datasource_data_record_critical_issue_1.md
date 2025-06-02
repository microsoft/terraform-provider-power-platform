# Title

Incorrect Error Handling in the Read Function

##

`/workspaces/terraform-provider-power-platform/internal/services/data_record/datasource_data_record.go`

## Problem

In the `Read` function of the DataRecordDataSource struct, the variable initialization of query and headers is not verified for the presence of errors effectively before use. Furthermore, the error message in `resp.Diagnostics.AddError` lacks appropriate context and detail.

## Impact

Errors are handled inappropriately, leading to incomplete diagnostic messages when an unexpected error occurs during query generation. This impacts debuggability and maintainability of the code. Severity: **critical**

## Location

Located in the `Read` function of the `DataRecordDataSource` struct.

## Code Issue

```go
query, headers, err := BuildODataQueryFromModel(&config)
tflog.Debug(ctx, fmt.Sprintf("Query: %s", query))
tflog.Debug(ctx, fmt.Sprintf("Headers: %v", headers))
if err != nil {
    resp.Diagnostics.AddError("Failed to build OData query", err.Error())
}
tflog.Debug(ctx, fmt.Sprintf("Query: %s", query))
```

## Fix

Modify the error-handling logic to validate errors explicitly. Add more context to the error diagnostic for better debugging capabilities.

```go
query, headers, err := BuildODataQueryFromModel(&config)
if err != nil {
    resp.Diagnostics.AddError(
        "Failed to build OData query",
        fmt.Sprintf("An error occurred while building the query. Details: %s", err.Error()),
    )
    return
}

tflog.Debug(ctx, fmt.Sprintf("Query: %s", query))
tflog.Debug(ctx, fmt.Sprintf("Headers: %v", headers))
```