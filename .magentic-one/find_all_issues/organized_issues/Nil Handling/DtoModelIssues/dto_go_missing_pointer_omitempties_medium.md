# Missing `omitempty` and Pointer Usage for Structs

## /workspaces/terraform-provider-power-platform/internal/services/environment/dto.go

## Problem

Several structs (such as `environmentCreatePropertiesDto` and others) have fields that are other structs (not pointers) tagged with `omitempty` for JSON marshaling. However, non-pointer structs are always included and cannot be omitted when using `omitempty`; only a pointer or an interface can be omitted. For example, fields such as `BillingPolicy BillingPolicyDto  `json:"billingPolicy,omitempty"`` cannot actually be omitted, as the zero value of a struct (not nil) is always present.

## Impact

This can lead to the unintended inclusion of empty objects in the marshaled JSON output, misleading API consumers and violating expectations set by the field tags. This is a **medium** severity issue as it affects API compatibility and data transfer correctness.

## Location

- `environmentCreatePropertiesDto`, field `BillingPolicy` (and potentially similar patterns elsewhere)

## Code Issue

```go
BillingPolicy BillingPolicyDto `json:"billingPolicy,omitempty"`
```

## Fix

Change the field to a pointer so it can properly be omitted if not set:

```go
BillingPolicy *BillingPolicyDto `json:"billingPolicy,omitempty"`
```

Apply this pattern for all struct-type fields that are intended to be omitted via `omitempty`.
