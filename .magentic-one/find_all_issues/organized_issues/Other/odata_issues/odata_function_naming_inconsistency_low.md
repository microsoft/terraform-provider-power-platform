# Issue 4: Unexported Function Naming Inconsistency

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/odata.go

## Problem

Some functions use a capitalized "Odata" (e.g., `buildOdataApplyPart`) while others use "OData" (e.g., `buildODataOrderByPart`). This inconsistency can be confusing.

## Impact

Reduces maintainability and clarity. Severity is **low**.

## Location

```go
func buildOdataApplyPart(apply *string) *string {
```

## Fix

Rename to `buildODataApplyPart` for consistency.

```go
func buildODataApplyPart(apply *string) *string {
```
