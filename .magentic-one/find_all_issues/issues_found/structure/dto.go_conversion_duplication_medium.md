# Title

Duplication of Conversion Logic Across Similar Functions

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/dto.go

## Problem

Many conversion functions are highly similar, each translating a specific DTO segment, but repeating the same patterns of attribute type creation, null checks, and attribute value assignment. This duplication leads to a proliferation of boilerplate, increases the maintenance burden, and makes changes prone to error.

## Impact

- Harder to maintain: any fixes or enhancements must be applied in multiple places.
- Higher risk of inconsistency and subtle bugs, as some functions may drift apart over time.
- Difficult to refactor at a systemic level because of scattered duplications.
- Reduces clarity, especially for newcomers or reviewers. Severity: medium.

## Location

- Functions like: `convertUserManagementSettings`, `convertCatalogSettings`, `convertPowerAppsSettings`, `convertLicensingSettings`, etc.

## Code Issue

```go
func convertUserManagementSettings(tenantSettingsDto tenantSettingsDto) (basetypes.ObjectType, basetypes.ObjectValue) {
    attrTypesUserManagementSettings := map[string]attr.Type{
        "enable_delete_disabled_user_in_all_environments": types.BoolType,
    }

    if tenantSettingsDto.PowerPlatform == nil || tenantSettingsDto.PowerPlatform.UserManagementSettings == nil {
        return types.ObjectType{AttrTypes: attrTypesUserManagementSettings}, types.ObjectNull(attrTypesUserManagementSettings)
    }
    attrValuesUserManagementSettings := map[string]attr.Value{
        "enable_delete_disabled_user_in_all_environments": types.BoolPointerValue(tenantSettingsDto.PowerPlatform.UserManagementSettings.EnableDeleteDisabledUserinAllEnvironments),
    }
    return types.ObjectType{AttrTypes: attrTypesUserManagementSettings}, types.ObjectValueMust(attrTypesUserManagementSettings, attrValuesUserManagementSettings)
}
```
(similar for other resource-specific convertXSettings functions)

## Fix

Create a set of utility functions or higher-order abstractions that can factor out common conversion patterns, e.g.:
- Generic conversion for single-field structs/DTOs.
- Table-driven factories for repeated mappings.

Example:

```go
func convertSingleBoolField(dto interface{}, attrName string, fieldPtr *bool) (basetypes.ObjectType, basetypes.ObjectValue) {
    attrTypes := map[string]attr.Type{attrName: types.BoolType}
    if fieldPtr == nil {
        return types.ObjectType{AttrTypes: attrTypes}, types.ObjectNull(attrTypes)
    }
    attrValues := map[string]attr.Value{attrName: types.BoolPointerValue(fieldPtr)}
    return types.ObjectType{AttrTypes: attrTypes}, types.ObjectValueMust(attrTypes, attrValues)
}
```

Apply this principle to remove boilerplate and improve code quality throughout the conversion functions.

