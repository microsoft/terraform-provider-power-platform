# Title

Skipping DTO Conversion for Empty IDs in `convertAllowedTenantsFromDto`

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/models.go`

## Problem

The code in the `convertAllowedTenantsFromDto` function skips DTOs if their `TenantId` is empty (`if dtoTenant.TenantId == "" { continue }`). While this prevents adding invalid entries, it silently overlooks these entries without logging or warning.

## Impact

Silently skipping entries during conversions can lead to inconsistencies and difficult-to-debug errors if expected data is silently discarded. Severity: **Medium**

## Location

Line 138 in `convertAllowedTenantsFromDto`

## Code Issue

```go
if dtoTenant.TenantId == "" {
	continue
}
```

## Fix

Include logging or diagnostics to indicate why entries are skipped:

```go
if dtoTenant.TenantId == "" {
	fmt.Printf("Skipping tenant conversion due to empty TenantId\n")
	continue
}
```