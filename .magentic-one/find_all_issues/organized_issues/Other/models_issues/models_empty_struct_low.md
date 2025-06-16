# Issue: Empty Struct Definition Redundancy

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/models.go

## Problem

The `PowerPagesSettings` struct is defined as an empty struct:

```go
type PowerPagesSettings struct {
}
```

If this struct serves as a placeholder for future expansion, that's acceptable, but if not, it represents unnecessary code and may confuse maintainers.

## Impact

Including unused or redundant empty structs can reduce code clarity, adding unneeded noise to the codebase. Severity: **low**.

## Location

- Around line 94

## Code Issue

```go
type PowerPagesSettings struct {
}
```

## Fix

- If you do not plan to add fields, remove the struct.
- If it is genuinely required as a placeholder, add a clarifying comment.

```go
// Placeholder for future PowerPages settings, intentionally left empty
type PowerPagesSettings struct{}
```
