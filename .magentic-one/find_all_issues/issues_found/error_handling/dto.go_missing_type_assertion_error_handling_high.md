# Title

Missing Error Handling When Using Type Assertion (`.(basetypes.ObjectValue)`)

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/dto.go

## Problem

Multiple conversion functions in the file perform type assertions like `object.(basetypes.ObjectValue)` without checking for assertion failures. If the value does not implement `basetypes.ObjectValue`, a runtime panic will occur.

## Impact

Failure to check the result of the type assertion may cause a panic that crashes the provider at runtime. This is critical for a Terraform provider, as user input or upstream changes may invalidate the assumed interface, causing unexplained provider failures and loss of state. Severity: high.

## Location

- Functions like `convertTeamsIntegrationModel`, `convertPowerAppsModel`, `convertPowerAutomateModel`, etc.

## Code Issue

```go
var teamsIntegrationSettings TeamsIntegrationSettings
teamIntegrationObject.(basetypes.ObjectValue).As(ctx, &teamsIntegrationSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
```

## Fix

Use a type assertion with `ok` check and handle the error gracefully.

```go
objectValue, ok := teamIntegrationObject.(basetypes.ObjectValue)
if !ok {
    // handle the error: return, or log as appropriate
    return // or possibly log or add diagnostics
}
objectValue.As(ctx, &teamsIntegrationSettings, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
```

Apply similar checks in all conversion functions where type assertions are used.

