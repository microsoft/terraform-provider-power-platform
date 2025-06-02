# Title

Incorrect Declaration of `DynamicsColumnsValidator` Type

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/dynamic_columns.go

## Problem

The function `DynamicColumns` is returning a new configuration validator of type `DynamicsColumnsValidator`. However, the declaration of the `DynamicsColumnsValidator` type is missing entirely in the provided code. Without this type existing or its definition being properly included, the function implementation is incomplete and will result in a runtime error.

## Impact

This mistake causes the application to fail at compile-time if the `DynamicsColumnsValidator` type is undefined, effectively breaking the functionality that depends on this function. This is a **critical issue** because it prevents the application from running, and such code cannot be deployed or executed in any environment.

Severity: **critical**

## Location

File Path: `/workspaces/terraform-provider-power-platform/internal/services/data_record/dynamic_columns.go`

## Code Issue

```go
func DynamicColumns(expression path.Expression) resource.ConfigValidator {
	return &DynamicsColumnsValidator{ // Problem: DynamicsColumnsValidator type is not defined anywhere.
		PathExpression: expression,
	}
}
```

## Fix

Define the missing `DynamicsColumnsValidator` type in the appropriate location within the code. For example:

```go
// DynamicsColumnsValidator is a custom implementation of resource.ConfigValidator.
type DynamicsColumnsValidator struct {
	PathExpression path.Expression
}

// ValidateConfig performs the validation logic for many-to-one relationships.
func (v *DynamicsColumnsValidator) ValidateConfig() resource.ValidationResponse {
	// Implement specific validation logic here.
	return resource.ValidationResponse{}
}
```

This definition resolves the missing type issue, enabling proper compilation and usage of the `DynamicColumns` function.
