# Title

Error propagation in `Read` method isn’t comprehensive

## 

`/workspaces/terraform-provider-power-platform/internal/services/licensing/datasource_billing_policies.go`

## Problem

The `Read` method only provides diagnostic messages when client errors occur, but does not add contextual information about the operation that triggered the error.

## Impact

Error diagnostics are not informative enough, which can hinder debugging for operations with multiple sub-steps. Severity: **Medium**

## Location

`Read` method error handling block.

## Code Issue

```go
if err != nil {
    resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
    return
}
```

## Fix

Enhance the error message to include more context, such as which step of the read operation caused the issue.

```go
if err != nil {
    resp.Diagnostics.AddError(
        fmt.Sprintf("Failed to retrieve billing policies for %s", d.FullTypeName()),
        fmt.Sprintf("An error occurred while fetching billing policies from the LicensingClient: %s. "+
                    "Verify that your configuration and API credentials are correct.", err.Error()),
    )
    return
}
```
