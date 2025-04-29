# Title

Incorrect Use of `Private.SetKey` in `Create` Function

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go`

## Problem

In the `Create` function, `resp.Private.SetKey` is used to store `original_settings` as a key, but it does not validate whether the key is successfully set in case of an error. This presents the risk of incomplete or incorrect internal state storage, especially in scenarios where `resp.Private` encounters issue during execution.

## Impact

This can lead to unreliability in resource state preservation and debugging inconsistency. If `Private.SetKey` silently fails, it can corrupt resource data handling during future operations (e.g., Delete). Severity: **High**.

## Location

Line 314: Inside the `Create` function, in the block handling private state storage.

## Code Issue

```go
resp.Private.SetKey(ctx, "original_settings", jsonSettings)
```

## Fix

Validate the response of `SetKey` and append relevant diagnostics when an error occurs. This ensures visibility into state saving failures and safeguards subsequent operations.

```go
diag := resp.Private.SetKey(ctx, "original_settings", jsonSettings)
if diag.HasError() {
    resp.Diagnostics.AddError(
        "Failed to Set Private Key",
        fmt.Sprintf("An error occurred while setting the private key: %v", diag),
    )
    return
}
```