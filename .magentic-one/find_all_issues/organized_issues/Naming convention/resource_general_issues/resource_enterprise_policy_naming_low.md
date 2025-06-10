# Issue: Function Name Typo - `NewEnterpisePolicyResource`

##

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/resource_enterprise_policy.go

## Problem

The factory function for this resource is named `NewEnterpisePolicyResource`, which has a typographical error ("Enterpise" should be "Enterprise"). This inconsistency in naming can lead to confusion, reduce code readability, and make it more challenging for other developers to use or search for this constructor function.

## Impact

- Reduces code maintainability and readability.
- Can create difficulty during code search or automated tooling.
- Could cause subtle bugs if the function is mistakenly called elsewhere.
- **Severity:** Low

## Location

Line where the function is defined:

```go
func NewEnterpisePolicyResource() resource.Resource {
	return &Resource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "enterprise_policy",
		},
	}
}
```

## Code Issue

```go
func NewEnterpisePolicyResource() resource.Resource {
	return &Resource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "enterprise_policy",
		},
	}
}
```

## Fix

Rename the function to `NewEnterprisePolicyResource`, making sure you also update all usage references.

```go
func NewEnterprisePolicyResource() resource.Resource {
	return &Resource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "enterprise_policy",
		},
	}
}
```

---

This markdown detailing the naming issue will be saved under:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/naming/resource_enterprise_policy_naming_low.md`
