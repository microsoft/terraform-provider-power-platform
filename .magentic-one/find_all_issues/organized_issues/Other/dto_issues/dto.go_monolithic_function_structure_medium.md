# Title

Large Monolithic Functions Without Decomposition

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/dto.go

## Problem

Key functions like `convertFromTenantSettingsModel` are excessively large and responsible for converting multiple sub-sections of the DTO. They manually call each individual conversion for platform sub-components, and directly assign to DTO fields, making them difficult to maintain, debug, or extend. Any small change to the structure requires editing multiple unrelated lines in a massive function.

## Impact

This practice leads to low maintainability, makes pinpointing logic errors more challenging, complicates unit testing, and increases the likelihood of merge conflicts. It adds cognitive load, and increases the risk of defects and technical debt. Severity: medium.

## Location

- `convertFromTenantSettingsModel`
- Similar conversion helpers that are highly repetitive and lengthy
- Functions with long repeated blocks like:

```go
if !tenantSettings.PowerPlatform.IsNull() && !tenantSettings.PowerPlatform.IsUnknown() {
    powerPlatformAttributes := tenantSettings.PowerPlatform.Attributes()
    err := convertSearchModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
    if err != nil {
        return tenantSettingsDto, err
    }
    convertTeamsIntegrationModel(ctx, powerPlatformAttributes, &tenantSettingsDto)
    // ...
}
```

## Code Issue

```go
func convertFromTenantSettingsModel(ctx context.Context, tenantSettings TenantSettingsResourceModel) (tenantSettingsDto, error) {
    tenantSettingsDto := tenantSettingsDto{}
    ...
    // many repeated blocks without decomposition
}
```

## Fix

Refactor large functions into a small set of responsibility-focused helpers. For example:

- Have a single function for each logical area (e.g., one to extract and assign `WalkMeOptOut`, one for `PowerPlatform`, etc.).
- Consider grouping repeated `if !IsNull() && !IsUnknown()` logic into a helper utility for clarity.
- Make subcomponent conversion chains use a table-driven approach or at least reduce boilerplate repetition.

Example (illustrative):

```go
func setBoolField(dest **bool, src basetypes.BoolValue) {
    if !src.IsNull() && !src.IsUnknown() {
        *dest = src.ValueBoolPointer()
    }
}

func convertFromTenantSettingsModel(...) ... {
    ...
    setBoolField(&tenantSettingsDto.WalkMeOptOut, tenantSettings.WalkMeOptOut)
    ...
}
```

This enhances code clarity, makes maintenance easier, and supports future unit testing at a finer granularity.

