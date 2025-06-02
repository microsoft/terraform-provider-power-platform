# Struct Field Naming Not Idiomatic (Mixed Case)

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/api_tenant.go

## Problem

The `Client` struct defines the field `Api`, which should be `API` for idiomatic Go (initialisms should be all capitalized).

## Impact

Deviating from naming conventions makes the code appear unprofessional and can confuse future maintainers. Severity: Low.

## Location

Line 18

## Code Issue

```go
type Client struct {
	Api *api.Client
}
```

## Fix

Use proper Go naming for initialisms:

```go
type Client struct {
	API *api.Client
}
```

Update any usage of `Api` to `API` accordingly.
