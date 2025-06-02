# Title

Hardcoded error messages with no localization support

##

`/workspaces/terraform-provider-power-platform/internal/services/authorization/resource_user.go`

## Problem

The error messages in the file are hardcoded strings, such as `"Unexpected ProviderData Type"` and `"Client error when creating"`. These do not support localization and impede scalability for multi-language applications. Moreover, hardcoded strings make maintenance harder and reduce flexibility for future enhancements.

## Impact

This issue is **medium severity**. While it does not cause immediate functional failure, it inhibits global application usability and can make error handling less robust. Errors will be harder to maintain and customize later.

## Location

Example of hardcoded error messages:
```go
resp.Diagnostics.AddError(
    "Unexpected ProviderData Type",
    fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
)

resp.Diagnostics.AddError(
    "Client error when creating",
    err.Error(),
)
```

## Code Issue

### Example

```go
resp.Diagnostics.AddError(
    "Unexpected ProviderData Type",
    fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
)

resp.Diagnostics.AddError(
    "Client error when creating a resource",
    err.Error(),
)
```

## Fix

Leverage constants or a structured format for error messages. Optionally integrate formatting or localization libraries to support multi-language functionality.

### Corrected Code

```go
const ErrUnexpectedProviderDataType = "Unexpected ProviderData Type"
const ErrClientCreatingResource = "Client error when creating a resource"

resp.Diagnostics.AddError(
    ErrUnexpectedProviderDataType,
    fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
)

resp.Diagnostics.AddError(
    ErrClientCreatingResource,
    err.Error(),
)
```

### Explanation

Centralizing error messages makes it easier to update them, expand support for multiple languages, and ensure consistent formatting across the application.
