# Title

Potential Error Wrapping Issue Missing in API Calls

##

/workspaces/terraform-provider-power-platform/internal/services/connection/api_connection.go

## Problem

Not all API errors are wrapped with contextual information. While some errors are wrapped using `customerrors.WrapIntoProviderError`, others are directly returned. Without consistent error wrapping, debugging and tracing issues becomes more difficult.

## Impact

This inconsistency in error handling can make it harder to identify the source and context of errors, impacting debuggability and support. Severity: Medium.

## Location

In methods such as:

```go
_, err := client.Api.Execute(ctx, nil, "PUT", apiUrl.String(), nil, connectionToCreate, []int{http.StatusCreated}, &connection)
if err != nil {
    return nil, err
}
```

## Code Issue

```go
if err != nil {
    return nil, err
}
```

## Fix

Consistently wrap errors returned from API calls with additional context using the existing error wrapping strategy or custom error types.

```go
if err != nil {
    return nil, fmt.Errorf("failed to create connection: %w", err)
}
```

Or using the project’s custom error approach if desired:

```go
if err != nil {
    return nil, customerrors.WrapIntoProviderError(err, customerrors.ERROR_OBJECT_CREATION_FAILED, "Failed to create connection")
}
```

---

*This file should be saved in:*
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/api_connection_error_wrapping_medium.md
