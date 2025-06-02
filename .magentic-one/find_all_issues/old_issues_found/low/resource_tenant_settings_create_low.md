# Title

Unnecessary Redundant `Get` Call in `Create` Function

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go`

## Problem

In the `Create` function, the `req.Plan.Get(ctx, &plan)` is invoked twice—once before saving the original tenant settings and again after. The second call validates the plan and retrieves data, but it is redundant since no modifications to `req.Plan` occur between the calls.

## Impact

Redundant operations can lead to performance degradation, especially in large-scale applications. Maintaining unnecessary calls also increases cognitive load on developers, making the code harder to optimize and debug. Severity: **Low**.

## Location

Lines 312 and 323: Inside the `Create` function.

## Code Issue

```go
// First occurrence
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
    return
}

// Second occurrence
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
    return
}
```

## Fix

Refactor the function to remove the redundant `Get` call, ensuring the `plan` is only retrieved once and reused thereafter.

```go
// Retrieve the plan once
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
    return
}

// Save the original tenant settings (no need to retrieve plan again)
originalSettings, erro := r.TenantSettingClient.GetTenantSettings(ctx)
if erro != nil {
    resp.Diagnostics.AddError(
        "Error reading tenant settings", fmt.Sprintf("Error reading tenant settings: %s", erro.Error()),
    )
    return
}
```