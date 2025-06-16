# Code Structure: `allowedTenantsObjectType(ctx context.Context)` Should Not Require Context

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/models.go

## Problem

The function `allowedTenantsObjectType(ctx context.Context)` takes a context parameter, but does not use it at all. The function simply returns a statically defined object type. Having a context parameter where none is needed can mislead readers and creates unnecessary complexity in function signatures.

## Impact

This is a **low severity** maintainability issue. It does not affect runtime, but can lead to confusion and poor API ergonomics for maintainers or contributors. It may also encourage the unnecessary passing of context where it is unneeded.

## Location

```go
// Define the object type for AllowedTenants set.
func allowedTenantsObjectType(ctx context.Context) types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"tenant_id": types.StringType,
			"inbound":   types.BoolType,
			"outbound":  types.BoolType,
		},
	}
}
```

## Fix

Remove the `ctx context.Context` parameter from the function (and update all usages). The signature becomes:

```go
// Define the object type for AllowedTenants set.
func allowedTenantsObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"tenant_id": types.StringType,
			"inbound":   types.BoolType,
			"outbound":  types.BoolType,
		},
	}
}
```
And update all call sites to omit the context argument.
