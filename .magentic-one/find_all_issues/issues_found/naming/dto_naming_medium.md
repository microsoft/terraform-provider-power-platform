# Naming convention not following Go standards

##

/workspaces/terraform-provider-power-platform/internal/services/currencies/dto.go

## Problem

Type names such as `currenciesDto`, `currenciesArrayDto`, and `currenciesPropertiesDto` do not follow Go's naming convention. In Go, exported types (those intended for use outside of the package or for clarity) should use PascalCase—e.g., `CurrenciesDTO`, `CurrenciesArrayDTO`, and `CurrenciesPropertiesDTO`. Additionally, abbreviations like `DTO` are often written in all caps for clarity.

## Impact

Not following Go naming conventions can lead to code that is difficult to read, inconsistent with other Go codebases, and harder to maintain. It also makes exported types inaccessible outside the package if they start with a lowercase letter. **Severity: Medium**

## Location

- Type definitions at lines 5-24

## Code Issue

```go
type currenciesDto struct {
	Value []currenciesArrayDto `json:"value"`
}

type currenciesArrayDto struct {
	Name       string                  `json:"name"`
	ID         string                  `json:"id"`
	Type       string                  `json:"type"`
	Properties currenciesPropertiesDto `json:"properties"`
}

type currenciesPropertiesDto struct {
	Code            string `json:"code"`
	Symbol          string `json:"symbol"`
	IsTenantDefault bool   `json:"isTenantDefault"`
}
```

## Fix

Change type names to follow Go PascalCase convention with proper abbreviation usage:

```go
type CurrenciesDTO struct {
	Value []CurrenciesArrayDTO `json:"value"`
}

type CurrenciesArrayDTO struct {
	Name       string                `json:"name"`
	ID         string                `json:"id"`
	Type       string                `json:"type"`
	Properties CurrenciesPropertiesDTO `json:"properties"`
}

type CurrenciesPropertiesDTO struct {
	Code            string `json:"code"`
	Symbol          string `json:"symbol"`
	IsTenantDefault bool   `json:"isTenantDefault"`
}
```
