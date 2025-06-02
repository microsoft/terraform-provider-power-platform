# Title

Improper Handling of Null ProviderData

##

/workspaces/terraform-provider-power-platform/internal/services/rest/datasource_rest_query.go

## Problem

In the `Configure` function, there is insufficient clarity in handling the scenario where `ProviderData` is null. Currently, the code skips processing in such cases without providing any diagnostic information to the user or contributor.

## Impact

This can lead to confusion during debugging or usage, especially if null `ProviderData` is caused by an unexpected configuration issue. Severity is **medium**, as it may result in difficult-to-trace runtime behavior.

## Location

The issue is present in the `Configure` function, specifically in the handling of a null `ProviderData` value.

## Code Issue

```go
if req.ProviderData == nil {
    // ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
    return
}
```

## Fix

Add diagnostic information in cases where `ProviderData` is null, making it easier to trace issues and avoid silent failures.

```go
if req.ProviderData == nil {
    resp.Diagnostics.AddWarning(
        "Null ProviderData",
        "ProviderData is null. This may be expected during ValidateConfig, but if encountered in other contexts, please check the configuration.",
    )
    return
}
```

This improvement ensures that users and contributors are informed about the status of `ProviderData` under all circumstances.