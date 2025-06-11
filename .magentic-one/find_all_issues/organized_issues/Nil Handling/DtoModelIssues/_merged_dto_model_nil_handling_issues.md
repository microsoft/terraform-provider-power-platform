# DTO and Model Nil Handling Issues

This document contains all identified nil handling issues related to DTO (Data Transfer Object) and model components in the Terraform Power Platform provider codebase.

## ISSUE 1

<!-- Source: dto.go-try_get_pointer_for_loop-high.md -->

# Type Safety Missed in `tryGetRuleValueFromDto` and Related Usage

##

/workspaces/terraform-provider-power-platform/internal/services/environment_group_rule_set/dto.go

## Problem

The `tryGetRuleValueFromDto` function returns a pointer to an element in the supplied slice:

```go
func tryGetRuleValueFromDto(values []environmentGroupRuleSetValueDto, valueId string) *environmentGroupRuleSetValueDto {
    for _, value := range values {
        if value.Id == valueId {
            return &value
        }
    }
    return nil
}
```

This returns the address of the iteration variable `value`. In Go, the variable is reused across loop iterations, so you end up returning a pointer to the last value in the slice, not the intended matching element.

## Impact

- **Severity:** High
- Subtle but critical bug: causes all returned pointers to reference the wrong memory location, leading to difficult-to-diagnose errors and unpredictable behavior.
- Data corruption or incorrect data returned from routines using this method.

## Location

```go
func tryGetRuleValueFromDto(values []environmentGroupRuleSetValueDto, valueId string) *environmentGroupRuleSetValueDto {
    for _, value := range values {
        if value.Id == valueId {
            return &value  // Pointer to for-loop variable - always the same address!
        }
    }
    return nil
}
```

## Fix

Return a pointer to the actual slice element using its index:

```go
func tryGetRuleValueFromDto(values []environmentGroupRuleSetValueDto, valueId string) *environmentGroupRuleSetValueDto {
    for i := range values {
        if values[i].Id == valueId {
            return &values[i]
        }
    }
    return nil
}
```

This way, the address returned is stable and correct for that element.

## ISSUE 2

<!-- Source: dto.go-unused-struct-tags-and-type-safety-improvement-high.md -->

# Unused Struct Tags and Type Safety Improvement

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/dto.go

## Problem

The code utilizes pointer types for fields like `IsDisabled`, `Inbound`, and `Outbound`, presumably to allow for omitempty or distinguish between unset and explicitly set false values. However, all public struct fields should ideally be documented for exported DTOs, and the pointer usage should be evaluated for clarity regarding their optional nature and to avoid potential nil dereferences without checks in higher layers. Moreover, struct tags should use full JSON tag strings (omitempty for pointers if intended).

## Impact

If these fields are used without proper nil checking, it can lead to runtime panics (high severity). The use of pointer types also hinders code readability and maintainability unless explicitly needed for omitting fields in JSON serialization.

## Location

- TenantIsolationPolicyPropertiesDto
- AllowedTenantDto
- DirectionDto

## Code Issue

```go
type TenantIsolationPolicyPropertiesDto struct {
 TenantId       string             `json:"tenantId"`
 IsDisabled     *bool              `json:"isDisabled,omitempty"`
 AllowedTenants []AllowedTenantDto `json:"allowedTenants"`
}

type AllowedTenantDto struct {
 TenantId  string       `json:"tenantId"`
 Direction DirectionDto `json:"direction"`
}

type DirectionDto struct {
 Inbound  *bool `json:"inbound"`
 Outbound *bool `json:"outbound"`
}
```

## Fix

Where pointer fields are kept, always perform nil checking before dereferencing. If omitempty is not truly required for IsDisabled, Inbound, or Outbound, switch to plain bool for improved safety. Otherwise, ensure usage patterns always account for nil cases.

```go
type TenantIsolationPolicyPropertiesDto struct {
 TenantId       string             `json:"tenantId"`
 IsDisabled     *bool              `json:"isDisabled,omitempty"`
 AllowedTenants []AllowedTenantDto `json:"allowedTenants"`
}

// When using
if props.IsDisabled != nil && *props.IsDisabled {
    // ...
}

// For non-pointer usage (if omitempty not essential):
type TenantIsolationPolicyPropertiesDto struct {
 TenantId       string             `json:"tenantId"`
 IsDisabled     bool               `json:"isDisabled"`
 AllowedTenants []AllowedTenantDto `json:"allowedTenants"`
}
```

## ISSUE 3

<!-- Source: dto_go_missing_pointer_omitempties_medium.md -->

# Missing `omitempty` and Pointer Usage for Structs

## /workspaces/terraform-provider-power-platform/internal/services/environment/dto.go

## Problem

Several structs (such as `environmentCreatePropertiesDto` and others) have fields that are other structs (not pointers) tagged with `omitempty` for JSON marshaling. However, non-pointer structs are always included and cannot be omitted when using `omitempty`; only a pointer or an interface can be omitted. For example, fields such as `BillingPolicy BillingPolicyDto`json:"billingPolicy,omitempty"`` cannot actually be omitted, as the zero value of a struct (not nil) is always present.

## Impact

This can lead to the unintended inclusion of empty objects in the marshaled JSON output, misleading API consumers and violating expectations set by the field tags. This is a **medium** severity issue as it affects API compatibility and data transfer correctness.

## Location

- `environmentCreatePropertiesDto`, field `BillingPolicy` (and potentially similar patterns elsewhere)

## Code Issue

```go
BillingPolicy BillingPolicyDto `json:"billingPolicy,omitempty"`
```

## Fix

Change the field to a pointer so it can properly be omitted if not set:

```go
BillingPolicy *BillingPolicyDto `json:"billingPolicy,omitempty"`
```

Apply this pattern for all struct-type fields that are intended to be omitted via `omitempty`.

## ISSUE 4

<!-- Source: models.go_lack_of_error_handling_for_element_conversion_medium.md -->

# Lack of Error Handling for Element Conversion in `convertToDto`

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/models.go

## Problem

In the `convertToDto` function, the error diagnostics from the call to `model.AllowedTenants.ElementsAs` are collected and appended to `diags`. There is a check immediately afterwards (`if diags.HasError() { return nil, diags }`), which is an appropriate pattern.

However, in the subsequent for-loop, there is no check for potentially malformed or incomplete data in `tenantsModel`. If `ElementsAs` partially fails, or if `AllowedTenants` contains `nil` or partial entries, this may lead to runtime issues when accessing values via methods like `allowedTenant.Inbound.ValueBool()`, but this will not be caught early. The handling for partial errors is insufficient if underlying data inconsistencies exist.

## Impact

If the tenant data is malformed, this could lead to panics or silent errors further down the line. This is a medium severity issue because it may not always surface as a problem, but could potentially break the resource unexpectedly and hinder debugging.

## Location

```go
func convertToDto(ctx context.Context, tenantId string, model *TenantIsolationPolicyResourceModel) (*TenantIsolationPolicyDto, diag.Diagnostics) {
 var diags diag.Diagnostics
 var tenantsModel []AllowedTenantModel
 diags.Append(model.AllowedTenants.ElementsAs(ctx, &tenantsModel, false)...)
 if diags.HasError() {
  return nil, diags
 }
 // ...
 for _, allowedTenant := range tenantsModel {
  inbound := allowedTenant.Inbound.ValueBool()
  outbound := allowedTenant.Outbound.ValueBool()
  dtoTenants = append(dtoTenants, AllowedTenantDto{
   TenantId: allowedTenant.TenantId.ValueString(),
   // ...
  })
 }
 // ...
}
```

## Fix

Add additional validation after conversion to ensure that the elements within `tenantsModel` are valid and non-`nil`. Defensive checks can be included before value extraction, and potential errors should be appended to diagnostics for the calling code to handle:

```go
 // After ensuring diags.HasError() is false:
 for i, allowedTenant := range tenantsModel {
  if allowedTenant.TenantId.IsNull() || allowedTenant.Inbound.IsNull() || allowedTenant.Outbound.IsNull() {
   diags.AddError(
    "Invalid AllowedTenantModel",
    fmt.Sprintf("Allowed tenant at index %d has missing required values", i),
   )
   continue // or return nil, diags if you prefer hard failure
  }
  inbound := allowedTenant.Inbound.ValueBool()
  outbound := allowedTenant.Outbound.ValueBool()
  dtoTenants = append(dtoTenants, AllowedTenantDto{
   TenantId: allowedTenant.TenantId.ValueString(),
   Direction: DirectionDto{
    Inbound:  &inbound,
    Outbound: &outbound,
   },
  })
 }
 if diags.HasError() {
  return nil, diags
 }
```

## ISSUE 5

<!-- Source: models_error_handling_medium.md -->

# Missing Error Handling in `convertFromConnectorDto` Function

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/models.go

## Problem

The function `convertFromConnectorDto` assumes that all fields within `connectorDto` and its nested `Properties` are present and initialized. If any of these are missing, nil, or otherwise invalid, this could result in a runtime panic or default zero-value assignment, which may not be desired behavior.

## Impact

Lack of error handling makes the system susceptible to panics or undetected data inconsistencies, especially when data is received from external sources. This can cause subtle bugs or service crashes. Severity is **medium**, due to risk of runtime errors when code is exposed to incomplete or malformed input.

## Location

- `convertFromConnectorDto(connectorDto connectorDto) DataSourceModel`

## Code Issue

```go
func convertFromConnectorDto(connectorDto connectorDto) DataSourceModel {
 return DataSourceModel{
  Id:          types.StringValue(connectorDto.Id),
  Name:        types.StringValue(connectorDto.Name),
  Type:        types.StringValue(connectorDto.Type),
  Description: types.StringValue(connectorDto.Properties.Description),
  DisplayName: types.StringValue(connectorDto.Properties.DisplayName),
  Tier:        types.StringValue(connectorDto.Properties.Tier),
  Publisher:   types.StringValue(connectorDto.Properties.Publisher),
  Unblockable: types.BoolValue(connectorDto.Properties.Unblockable),
 }
}
```

## Fix

Add validation for presence of necessary fields and gracefully handle missing or nil sub-structs:

```go
func convertFromConnectorDto(connectorDto connectorDto) DataSourceModel {
 var description, displayName, tier, publisher string
 var unblockable bool

 if connectorDto.Properties != nil {
  description = connectorDto.Properties.Description
  displayName = connectorDto.Properties.DisplayName
  tier = connectorDto.Properties.Tier
  publisher = connectorDto.Properties.Publisher
  unblockable = connectorDto.Properties.Unblockable
 }

 return DataSourceModel{
  Id:          types.StringValue(connectorDto.Id),
  Name:        types.StringValue(connectorDto.Name),
  Type:        types.StringValue(connectorDto.Type),
  Description: types.StringValue(description),
  DisplayName: types.StringValue(displayName),
  Tier:        types.StringValue(tier),
  Publisher:   types.StringValue(publisher),
  Unblockable: types.BoolValue(unblockable),
 }
}
```

This avoids potential panics and increases the robustness of the code when handling untrusted input.

## ISSUE 6

<!-- Source: models_nil_dereference_critical.md -->

# Title

Possible nil dereference when using pointers in `convertFromEnvironmentSettingsDto`

##

/workspaces/terraform-provider-power-platform/internal/services/environment_settings/models.go

## Problem

The function `convertFromEnvironmentSettingsDto` uses various pointer values with the `*` operator, directly dereferencing them without checking if the pointer is non-nil. For example, `*environmentSettingsDto.IsAuditEnabled` is used without verifying if `environmentSettingsDto.IsAuditEnabled != nil`.

## Impact

If the dto fields are nil (i.e., not set), this will cause a runtime panic due to a nil pointer dereference. This is a critical severity issue, as it leads to possible application crashes.

## Location

Example:

```go
attrValuesAuditSettingsProperties := map[string]attr.Value{
 "is_audit_enabled":             types.BoolValue(*environmentSettingsDto.IsAuditEnabled),
 // ... and others.
}
```

## Code Issue

```go
attrValuesAuditSettingsProperties := map[string]attr.Value{
 "is_audit_enabled":             types.BoolValue(*environmentSettingsDto.IsAuditEnabled),
 "is_user_access_audit_enabled": types.BoolValue(*environmentSettingsDto.IsUserAccessAuditEnabled),
 "is_read_audit_enabled":        types.BoolValue(*environmentSettingsDto.IsReadAuditEnabled),
 "log_retention_period_in_days": logRetentionPeriodTypeValue,
}
```

## Fix

Before dereferencing each pointer, check for non-nil and provide a default value if it is nil, e.g.:

```go
isAuditEnabled := false
if environmentSettingsDto.IsAuditEnabled != nil {
 isAuditEnabled = *environmentSettingsDto.IsAuditEnabled
}

attrValuesAuditSettingsProperties := map[string]attr.Value{
 "is_audit_enabled": types.BoolValue(isAuditEnabled),
 // repeat for others
}
```

This will prevent panic and ensure safe usage.

---

# To finish the task you have to

1. Run linter and fix any issues
2. Run UnitTest and fix any of failing ones
3. Generate docs
4. Run Changie

# Changie Instructions

Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```

Where:

- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number

Apply this fix to the whole codebase
