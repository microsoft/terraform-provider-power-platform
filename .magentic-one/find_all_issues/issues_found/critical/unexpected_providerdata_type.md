# Title
Unexpected ProviderData Type

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/datasource_dlp_policy.go

## Problem
The `Configure` function does not handle the scenario where `ProviderData` is not of type `*api.ProviderClient` gracefully. Instead, it adds diagnostic errors, which might not be the best approach for handling type mismatches.

## Impact
If a type mismatch occurs, the function fails abruptly and assumes it’s an error in the provider. This might confuse users or developers troubleshooting the issue. Severity: critical.

## Location
Function: Configure

## Code Issue
```go
if req.ProviderData == nil {
    // ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
    return
}
client, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
    resp.Diagnostics.AddError(
        "Unexpected ProviderData Type",
        fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )
    return
}

```

## Fix
Use an approach that logs a warning and avoids stopping execution entirely when a type mismatch occurs. Provide a default behavior or fallback mechanism.
```go
if req.ProviderData == nil {
    // ProviderData will be null when Configure is called from ValidateConfig, logging warning.
    tflog.Warn(ctx, "ProviderData is null, default configuration applied.")
    return
}
client, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
    tflog.Warn(ctx, fmt.Sprintf("Unexpected ProviderData type: %T. Proceeding with default configuration.", req.ProviderData))
    // Optionally implement fallback behavior here.
    return
}

```