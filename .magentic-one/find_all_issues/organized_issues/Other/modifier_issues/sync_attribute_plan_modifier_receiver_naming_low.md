# Title

Non-idiomatic Receiver Name (`d`) in struct methods

##

/workspaces/terraform-provider-power-platform/internal/modifiers/sync_attribute_plan_modifier.go

## Problem

The receiver for the `syncAttributePlanModifier` struct's methods is named `d`, which is non-descriptive and typically reserved for types like “decoder.” Best practice is to use a short form of the struct type for clarity and maintainability (e.g., `sam` or `m`).

## Impact

Reduces code readability and maintainability, especially in larger codebases or for new contributors. Severity: **low**.

## Location

Throughout struct method receivers, e.g.:
```go
func (d *syncAttributePlanModifier) Description(ctx context.Context) string { ... }
```

## Code Issue

```go
func (d *syncAttributePlanModifier) Description(ctx context.Context) string { ... }
```

## Fix

Change receiver name to something more descriptive, such as `sam`:

```go
func (sam *syncAttributePlanModifier) Description(ctx context.Context) string { ... }
```
