# Title

Typo in variable name: `linkEnterprosePolicyDto`

##

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/api_enterprise_policy.go

## Problem

The variable `linkEnterprosePolicyDto` contains a typographical error in its name. The word "enterprise" is misspelled as "enterprose."

## Impact

The typo in the variable name negatively affects code readability and maintainability. Developers might struggle with understanding or may introduce errors while referencing it in the code elsewhere. The severity of this issue is **low**, as it does not cause runtime errors but does hinder effective collaboration.

## Location

- `LinkEnterprisePolicy` function at multiple places
- `UnLinkEnterprisePolicy` function at multiple places

## Code Issue

```go
linkEnterprosePolicyDto := linkEnterprosePolicyDto{
	SystemId: systemId,
}
```

## Fix

Correct the variable name to `linkEnterprisePolicyDto` for better readability and consistency. Also, update the type wherever it is declared in the codebase.

```go
linkEnterprisePolicyDto := linkEnterprisePolicyDto{
	SystemId: systemId,
}
```

This fix improves the code readability and aligns variable naming conventions with the intended meaning. Update all occurrences in both `LinkEnterprisePolicy` and `UnLinkEnterprisePolicy` functions to ensure consistency.
