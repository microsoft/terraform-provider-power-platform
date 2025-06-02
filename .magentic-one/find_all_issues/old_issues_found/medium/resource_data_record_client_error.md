# Title

Generalized Client Error Handling in Multiple Functions

## Path

`/workspaces/terraform-provider-power-platform/internal/services/data_record/resource_data_record.go`

## Problem

In several functions like `Create`, `Update`, and `Delete`, client error handling is overly generalized without specifying the exact reason or types of client failures:

```go
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
    return
}
```

This lack of granularity in error reporting can make it difficult to accurately debug issues or comprehend the underlying client behavior.

## Impact

The lack of specific error handling reduces debugging efficiency and diagnostic clarity, making it harder to pinpoint the failure scenario in production environments. Severity: **Medium**

## Location

Functions affected:
1. Create
2. Update
3. Delete

File: `/internal/services/data_record/resource_data_record.go`

## Code Issue

```go
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
    return
}
```

## Fix

Introduce error categorization and detailed client error reporting. For example:

```go
if err != nil {
    if clientError, ok := err.(ClientSpecificErrorType); ok {
        resp.Diagnostics.AddError(
            fmt.Sprintf("Specific client error in method %s", r.FullTypeName()),
            fmt.Sprintf("Client-specific issue: %s", clientError.Message()),
        )
    } else {
        resp.Diagnostics.AddError(
            fmt.Sprintf("Unknown client error in %s", r.FullTypeName()),
            fmt.Sprintf("Error details: %s", err.Error()),
        )
    }
    return
}
```