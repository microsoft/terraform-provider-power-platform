# Duplicate Assignment to `model.EnterprisePolicies` in `convertEnterprisePolicyModelFromDto`

##

/workspaces/terraform-provider-power-platform/internal/services/environment/models.go

## Problem

In the function `convertEnterprisePolicyModelFromDto`, if both `Vnets` and `CustomerManagedKeys` are present in `environmentDto.Properties.EnterprisePolicies`, only one will be retained in `model.EnterprisePolicies` due to duplicate assignment. Each assignment overwrites the previous, so only the last one (`CustomerManagedKeys` if present) will persist.

## Impact

- **Severity:** High
- Information about multiple enterprise policies (such as both Vnets and CustomerManagedKeys) will not be reflected, leading to missing or incomplete policy data.
- Can cause data inconsistency and unexpected behavior.

## Location

```go
if environmentDto.Properties.EnterprisePolicies != nil {
	if environmentDto.Properties.EnterprisePolicies.Vnets != nil {
		model.EnterprisePolicies = types.SetValueMust(enterprisePolicyAttrType, []attr.Value{ ... })
	}
	if environmentDto.Properties.EnterprisePolicies.CustomerManagedKeys != nil {
		model.EnterprisePolicies = types.SetValueMust(enterprisePolicyAttrType, []attr.Value{ ... })
	}
}
```

## Code Issue

```go
if environmentDto.Properties.EnterprisePolicies != nil {
	if environmentDto.Properties.EnterprisePolicies.Vnets != nil {
		model.EnterprisePolicies = types.SetValueMust(enterprisePolicyAttrType, []attr.Value{
			types.ObjectValueMust(
				...
			),
		})
	}
	if environmentDto.Properties.EnterprisePolicies.CustomerManagedKeys != nil {
		model.EnterprisePolicies = types.SetValueMust(enterprisePolicyAttrType, []attr.Value{
			types.ObjectValueMust(
				...
			),
		})
	}
}
```

## Fix

Accumulate all enterprise policy values in a temporary slice, and then set `model.EnterprisePolicies` once, preserving all collected information:

```go
if environmentDto.Properties.EnterprisePolicies != nil {
	var enterprisePolicyValues []attr.Value
	if environmentDto.Properties.EnterprisePolicies.Vnets != nil {
		enterprisePolicyValues = append(enterprisePolicyValues, types.ObjectValueMust(
			...
		))
	}
	if environmentDto.Properties.EnterprisePolicies.CustomerManagedKeys != nil {
		enterprisePolicyValues = append(enterprisePolicyValues, types.ObjectValueMust(
			...
		))
	}
	model.EnterprisePolicies = types.SetValueMust(enterprisePolicyAttrType, enterprisePolicyValues)
} else {
	model.EnterprisePolicies = types.SetValueMust(enterprisePolicyAttrType, []attr.Value{})
}
```
