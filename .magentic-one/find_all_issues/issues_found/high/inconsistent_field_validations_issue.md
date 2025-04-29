# Title

Inconsistent Field Validations in DTO Conversions

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/models.go`

## Problem

The `convertAllowedTenantsFromDto` function assumes DTO contains valid fields but does not validate `Inbound` and `Outbound` directions fully. The same issue extends to `convertToDto`.

## Impact

Assuming non-validated inputs can cause runtime errors or unexpected behavior. Severity: **High**

## Location

Lines 134-150 in `convertAllowedTenantsFromDto` and 33-50 in `convertToDto`.

## Code Issue

```go
if dtoTenant.Direction.Inbound != nil {
	inbound = *dtoTenant.Direction.Inbound
}
if dtoTenant.Direction.Outbound != nil {
	outbound = *dtoTenant.Direction.Outbound
}
```

```go
inbound := allowedTenant.Inbound.ValueBool()
outbound := allowedTenant.Outbound.ValueBool()
```

## Fix

Introduce explicit validation for all fields before attempting to convert:

```go
if dtoTenant.Direction.Inbound == nil || dtoTenant.Direction.Outbound == nil {
	fmt.Printf("Skipping tenant due to invalid direction fields\n")
	continue
}
```