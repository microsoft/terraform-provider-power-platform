# DTO Naming Conventions - Use Consistent PascalCase for Acronyms

##

/workspaces/terraform-provider-power-platform/internal/services/locations/dto.go

## Problem

The type name for the main DTO is suggested as `LocationDTO`, but acronyms in Go should usually use `DTO` (not `Dto`). Consistency across the codebase is important (alternatively, some teams prefer `Dto` for readability). The current names (`locationDto`, `locationsArrayDto`) do not follow this convention.

## Impact

Severity: **Low**

This is a stylistic issue, but can hurt code readability, hinder onboarding, and result in confusion if other files use a different convention.

## Location

Type definitions for DTOs.

## Code Issue

```go
type locationDto struct {
	Value []locationsArrayDto `json:"value"`
}

type locationsArrayDto struct {
	// ...
}
```

## Fix

Decide on a standard for acronym casing in type names (common in Go is PascalCase, e.g., `LocationDTO`). Refactor all DTO types accordingly:

```go
type LocationDTO struct {
	Value []LocationsArrayDTO `json:"value"`
}

type LocationsArrayDTO struct {
	// ...
}
```
