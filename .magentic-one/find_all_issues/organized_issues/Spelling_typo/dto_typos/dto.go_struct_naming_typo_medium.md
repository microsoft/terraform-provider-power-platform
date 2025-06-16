# Issue: Struct Naming Typo

##

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/dto.go

## Problem

The struct is named `linkEnterprosePolicyDto`, but the word "Enterprose" is likely a typo for "Enterprise". Correcting this typo will improve code readability and maintainability.

## Impact

Incorrect naming can lead to confusion, reduced maintainability, and potential difficulty when using tooling or searching the codebase. Severity: Medium

## Location

Line defining the struct:

```go
type linkEnterprosePolicyDto struct {
    SystemId string `json:"systemId"`
}
```

## Code Issue

```go
type linkEnterprosePolicyDto struct {
    SystemId string `json:"systemId"`
}
```

## Fix

Rename the struct to `linkEnterprisePolicyDto` (and refactor any usage accordingly):

```go
type linkEnterprisePolicyDto struct {
    SystemId string `json:"systemId"`
}
```
