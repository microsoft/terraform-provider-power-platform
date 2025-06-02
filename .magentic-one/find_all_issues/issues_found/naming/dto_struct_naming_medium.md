# Title

Type Naming Inconsistency - Struct Name Should Be Exported

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/dto.go

## Problem

The struct `billingPolicyCreateDto` uses camelCase naming, making it unexported. In Go, if this struct is intended to be used outside its package (which is likely in DTOs), it should use PascalCase to be exported and consistent with other struct names in this file. Furthermore, its naming style is inconsistent with `BillingInstrumentDto`, `BillingPolicyDto`, and others in this file.

## Impact

This causes confusion for consumers of the package, prevents the struct from being exported (if required), and reduces code consistency and maintainability. Severity is **medium** because it may cause problems in package usage and code understanding.

## Location

Top of the file, first struct definition.

## Code Issue

```go
type billingPolicyCreateDto struct {
	Location          string               `json:"location"`
	Name              string               `json:"name"`
	Status            string               `json:"status"`
	BillingInstrument BillingInstrumentDto `json:"billingInstrument"`
}
```

## Fix

Rename the struct to `BillingPolicyCreateDto` for export and consistency.

```go
type BillingPolicyCreateDto struct {
	Location          string               `json:"location"`
	Name              string               `json:"name"`
	Status            string               `json:"status"`
	BillingInstrument BillingInstrumentDto `json:"billingInstrument"`
}
```
