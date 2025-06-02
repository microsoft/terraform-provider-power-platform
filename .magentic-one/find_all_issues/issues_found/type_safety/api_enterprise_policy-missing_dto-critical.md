# Missing Definition of DTO Struct

##

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/api_enterprise_policy.go

## Problem

The struct `linkEnterprosePolicyDto` is referenced (used in variable and as a type) but is not defined anywhere in the visible code. This will cause a compile-time error.

## Impact

Code will not compile/run. Critical for correctness.

## Location

All locations where `linkEnterprosePolicyDto` is referenced, for example:

```go
linkEnterprosePolicyDto := linkEnterprosePolicyDto{
	SystemId: systemId,
}
```

## Code Issue

```go
linkEnterprosePolicyDto := linkEnterprosePolicyDto{
	SystemId: systemId,
}
```

## Fix

Define the struct (with correct spelling as well):

```go
type linkEnterprisePolicyDto struct {
	SystemId string `json:"systemId"`
}
```

Then use:

```go
linkEnterprisePolicyDto := linkEnterprisePolicyDto{
	SystemId: systemId,
}
```
