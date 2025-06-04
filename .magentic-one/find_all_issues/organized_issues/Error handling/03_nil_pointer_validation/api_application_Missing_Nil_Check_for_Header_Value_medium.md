# Issue 3: API Client: Missing Nil Check for Header Value

##

/workspaces/terraform-provider-power-platform/internal/services/application/api_application.go

## Problem

There’s no check if `operationLocationHeader` is empty before usage. If the response header is missing, it could lead to confusing behavior or errors.

## Impact

This issue has a **medium** severity. Skipping this check could lead to downstream HTTP requests with an empty URL or crash from nil/empty string dereferences.

## Location

In InstallApplicationInEnvironment after retrieving `operationLocationHeader`:

## Code Issue

```go
operationLocationHeader := response.GetHeader(constants.HEADER_OPERATION_LOCATION)
tflog.Debug(ctx, "Operation Location Header: "+operationLocationHeader)

_, err = url.Parse(operationLocationHeader)
if err != nil {
    tflog.Error(ctx, "Error parsing location header: "+err.Error())
    return "", err
}
```

## Fix

Add a conditional check for the header before using it:

```go
operationLocationHeader := response.GetHeader(constants.HEADER_OPERATION_LOCATION)
if operationLocationHeader == "" {
    tflog.Error(ctx, "Missing operation location header in response")
    return "", errors.New("missing operation location header in response")
}
tflog.Debug(ctx, "Operation Location Header: "+operationLocationHeader)

_, err = url.Parse(operationLocationHeader)
if err != nil {
    tflog.Error(ctx, "Error parsing location header: "+err.Error())
    return "", err
}
```
