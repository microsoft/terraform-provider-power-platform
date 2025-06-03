# Title

Lack of Unit Tests for Custom Plan Modifier

##

/workspaces/terraform-provider-power-platform/internal/modifiers/sync_attribute_plan_modifier.go

## Problem

There's no evidence of unit tests for the custom `syncAttributePlanModifier`. Custom plan modifiers are essential for ensuring correct Terraform behavior and preventing regression, and should be tested for the various code branches (null, unknown, error, valid hash).

## Impact

Reduces reliability, maintainability, and increases risk of regression on this core behavior. Severity: **medium**.

## Location

Not in source file; absence of relevant `_test.go` for `sync_attribute_plan_modifier.go`.

## Code Issue

_No test support found for:_
- returning null or unknown values
- error handling when `CalculateSHA256` fails
- valid SHA256 calculation

## Fix

Create a corresponding test file with cases covering the major branches of `PlanModifyString` and the construction function.

```go
// file: internal/modifiers/sync_attribute_plan_modifier_test.go

func TestPlanModifyString_NullOrUnknown(t *testing.T) { /* ... */ }
func TestPlanModifyString_SHA256Error(t *testing.T) { /* ... */ }
func TestPlanModifyString_ValidHash(t *testing.T) { /* ... */ }
```
