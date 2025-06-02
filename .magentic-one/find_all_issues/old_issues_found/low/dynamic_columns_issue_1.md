# Title

Misleading or Incorrect Function Name: `DynamicColumns`

##

/workspaces/terraform-provider-power-platform/internal/services/data_record/dynamic_columns.go

## Problem

The function `DynamicColumns` is named as if it deals with dynamic columns, which may imply processing or handling related to dynamic structures or column data. However, its functionality does not align with the provided name in any explicit way. Instead, the function returns a configuration validator without explicitly showcasing the dynamic column processing.

## Impact

This name can create confusion among developers trying to understand or maintain the code. Misaligned naming conventions make the codebase harder to navigate and can result in improper use or expectations of the function.

Severity: **low**

## Location

File Path: `/workspaces/terraform-provider-power-platform/internal/services/data_record/dynamic_columns.go`

## Code Issue

```go
// DynamicColumns returns a ConfigValidator that ensures that the given expression
// many-to-one relationships are using set collections.
func DynamicColumns(expression path.Expression) resource.ConfigValidator {
	return &DynamicsColumnsValidator{
		PathExpression: expression,
	}
}
```

## Fix

To resolve this issue, rename the function to something that more accurately reflects its purpose, such as `CreateConfigValidator`. This makes it clear that its responsibility is to create and return a configuration validator.

```go
// CreateConfigValidator returns a ConfigValidator that ensures that the given expression
// many-to-one relationships are using set collections.
func CreateConfigValidator(expression path.Expression) resource.ConfigValidator {
	return &DynamicsColumnsValidator{
		PathExpression: expression,
	}
}
```

This change increases clarity and improves code maintainability.
