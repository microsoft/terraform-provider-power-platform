# Typo in Struct and Variable Naming

##

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/api_enterprise_policy.go

## Problem

The struct and variable names `linkEnterprosePolicyDto` are incorrectly spelled. The correct spelling should be `linkEnterprisePolicyDto`. This typo appears at variable declarations and literal assignments in both `LinkEnterprisePolicy` and `UnLinkEnterprisePolicy` functions.

## Impact

This impacts both code readability and maintainability. Developers unfamiliar with the code might be confused, and searches for `enterprise` will not find these entries. Severity: **low** (does not directly affect functionality but affects clarity and can be error-prone in future maintenance).

## Location

Lines where the DTO struct and variables are named.

## Code Issue

```go
linkEnterprosePolicyDto := linkEnterprosePolicyDto{
	SystemId: systemId,
}
```

## Fix

Rename all instances of `linkEnterprosePolicyDto` to `linkEnterprisePolicyDto` (struct and variable names). Confirm the type definition is also correctly named elsewhere.

```go
linkEnterprisePolicyDto := linkEnterprisePolicyDto{
	SystemId: systemId,
}
```
