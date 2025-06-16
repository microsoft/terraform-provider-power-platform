# Redundant initialization of `policies := make([]dlpPolicyModelDto, 0)`

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/api_dlp_policy.go

## Problem

You use `make([]Type, 0)`, while `var policies []dlpPolicyModelDto` suffices. Idiomatic Go is to use `var policies []Type`, or when the count is known, preallocate with capacity.

## Impact

Very minimal, but improves Go idiomatic usage and potential performance in tight code paths. (Severity: Low)

## Location

Lines 36

## Code Issue

```go
policies := make([]dlpPolicyModelDto, 0)
```

## Fix

```go
var policies []dlpPolicyModelDto
```

Or, if the length is known, preallocate to capacity:

```go
policies := make([]dlpPolicyModelDto, 0, len(policiesArray.Value))
```

