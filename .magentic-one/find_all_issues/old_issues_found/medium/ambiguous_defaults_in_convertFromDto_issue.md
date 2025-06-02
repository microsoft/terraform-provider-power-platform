# Title

Ambiguous Defaults in `convertFromDto` Function

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/models.go`

## Problem

The `convertFromDto` function initializes variables like `tenantId`, `isDisabled`, and `allowedTenants` with defaults. While this approach avoids `nil` issues, it does not explicitly handle cases where required properties might be missing in `dto.Properties`.

```go
tenantId := ""
var isDisabled *bool
var allowedTenants []AllowedTenantDto
```

## Impact

Ambiguous defaults make it hard for callers to detect and handle missing attributes, potentially causing unintended behavior. Severity: **Medium**

## Location

Line 68-71 in `convertFromDto`

## Code Issue

```go
tenantId := ""
var isDisabled *bool
var allowedTenants []AllowedTenantDto
```

## Fix

Add validation for missing or invalid properties in `dto.Properties` with explicit error reporting:
```go
if dto.Properties == nil {
	return TenantIsolationPolicyResourceModel{}, diag.NewErrorDiagnostic("Properties are nil")
}
```