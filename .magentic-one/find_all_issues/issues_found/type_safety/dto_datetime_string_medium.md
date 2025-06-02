# Title

Data Consistency: Use of String for Date Fields Instead of Time Type

##

/workspaces/terraform-provider-power-platform/internal/services/licensing/dto.go

## Problem

Fields such as `CreatedOn` and `LastModifiedOn` in `BillingPolicyDto` are defined as `string` instead of using Go's `time.Time` (with the standard `encoding/json` support via the `time` package). Using `string` for date/time values can lead to inconsistent date formats, lack of validation, and more error-prone handling of these fields.

## Impact

This can introduce bugs in date/time manipulation, cause inconsistency in how dates are serialized or deserialized, and make validation harder. Severity is **medium** as it affects data consistency and type safety, though it might also be API-driven.

## Location

`BillingPolicyDto` struct:

## Code Issue

```go
type BillingPolicyDto struct {
	// ...
	CreatedOn         string               `json:"createdOn"`
	CreatedBy         PrincipalDto         `json:"createdBy"`
	LastModifiedOn    string               `json:"lastModifiedOn"`
	LastModifiedBy    PrincipalDto         `json:"lastModifiedBy"`
}
```

## Fix

Use the `time.Time` type for date fields, and ensure that JSON marshaling/unmarshaling is handled correctly (e.g., via RFC3339 or whatever format your API expects).

```go
import "time"

type BillingPolicyDto struct {
	// ...
	CreatedOn         time.Time            `json:"createdOn"`
	CreatedBy         PrincipalDto         `json:"createdBy"`
	LastModifiedOn    time.Time            `json:"lastModifiedOn"`
	LastModifiedBy    PrincipalDto         `json:"lastModifiedBy"`
}
```

If custom formatting is required, provide the appropriate (un)marshal methods. Only use `string` if the upstream API dictates it *and* there is no ability to change serialization behavior.
