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
