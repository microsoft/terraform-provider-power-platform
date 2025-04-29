# Title

Error Handling in the Metadata Function

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/resources_environment_settings.go

## Problem

The `Metadata` function does not validate the `req.ProviderTypeName` or handle errors related to its absence or incorrect value.

## Impact

If the `req.ProviderTypeName` is missing or invalid, the function might behave unpredictably, leading to confusing debug logs or missing type settings. This could hinder troubleshooting and debugging efforts.

Severity: **High**

## Location

Line ~51: Function `Metadata`. 

## Code Issue

```go
func (r *EnvironmentSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    r.ProviderTypeName = req.ProviderTypeName

    ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
    defer exitContext()

    resp.TypeName = r.FullTypeName()
    tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}
```

## Fix

Add validation for `req.ProviderTypeName` to handle cases where it is missing or invalid. Explicit error handling should be implemented.

```go
func (r *EnvironmentSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    if req.ProviderTypeName == "" {
        resp.Diagnostics.AddError("Missing Provider Type Name", "ProviderTypeName from MetadataRequest is missing or empty.")
        return
    }

    r.ProviderTypeName = req.ProviderTypeName

    ctx, exitContext := helpers.EnterRequestContext(ctx, r.TypeInfo, req)
    defer exitContext()

    resp.TypeName = r.FullTypeName()
    tflog.Debug(ctx, fmt.Sprintf("METADATA: %s", resp.TypeName))
}
```
