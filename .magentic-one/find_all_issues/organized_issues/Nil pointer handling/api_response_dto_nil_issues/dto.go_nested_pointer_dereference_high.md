# Title

Lack of Null Checks for Nested Pointers in DTO Access

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/dto.go

## Problem

When accessing nested pointer fields (e.g., `tenantSettingsDto.PowerPlatform.Search.DisableDocsSearch`), many places in the code presume non-nil pointers. If a parent pointer (`PowerPlatform` or `Search`) is nil, this will cause a panic due to dereferencing a nil pointer.

## Impact

Dereferencing nil pointers leads to panics, causing a provider crash and potential Terraform workflow interruption. This is a critical reliability and user experience problem. Severity: high.

## Location

Example in `convertSearchSettings` and similar conversion functions:

```go
if tenantSettingsDto.PowerPlatform == nil || tenantSettingsDto.PowerPlatform.Search == nil {
    return types.ObjectType{AttrTypes: attrTypesSearchProperties}, types.ObjectNull(attrTypesSearchProperties)
}
attrValuesSearchProperties := map[string]attr.Value{
    "disable_docs_search":       types.BoolPointerValue(tenantSettingsDto.PowerPlatform.Search.DisableDocsSearch),
    ...
}
```

But in some paths, the code does not check for all possible nesting levels before accessing a field, which creates potential race opportunities for panics:

- In DTO setters inside the conversion chains.
- In reading DTO values elsewhere, relying only on cursor null-checks.

## Code Issue

```go
tenantSettingsDto.PowerPlatform.Search.DisableDocsSearch
```

## Fix

Always check for nil at each pointer access level before dereferencing:

```go
if tenantSettingsDto.PowerPlatform != nil && tenantSettingsDto.PowerPlatform.Search != nil {
    // safe to access tenantSettingsDto.PowerPlatform.Search.DisableDocsSearch
}
```

Apply this check everywhere nested pointers are dereferenced or read.

