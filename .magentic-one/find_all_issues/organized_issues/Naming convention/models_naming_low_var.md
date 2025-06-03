# Title

Misleading or Inconsistent Variable Naming in Set and ObjectType Variables

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/models.go

## Problem

Variable names like `customConnectorPatternSetObjectType`, `connectorSetObjectType`, `endpointRuleListObjectType`, and `actionRuleListObjectType` are inconsistent in their naming conventions (sometimes referring to "Set", sometimes "List"), which can be misleading. Furthermore, it is not clear if the type relates directly to a set or list in the schema, or to the internal model, by way of naming.

## Impact

Severity: Low

Poorly or inconsistently named variables reduce code readability and can lead to confusion, especially for new contributors or maintainers.

## Location

```go
var customConnectorPatternSetObjectType = types.ObjectType{ ... }
var endpointRuleListObjectType = types.ObjectType{ ... }
var actionRuleListObjectType = types.ObjectType{ ... }
```

## Code Issue

```go
var endpointRuleListObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"order":    types.Int64Type,
		"behavior": types.StringType,
		"endpoint": types.StringType,
	},
}
```

## Fix

Name these variables consistently according to their use:
```go
var customConnectorPatternObjectType = types.ObjectType{ ... }
var connectorObjectType = types.ObjectType{ ... }
var endpointRuleObjectType = types.ObjectType{ ... }
var actionRuleObjectType = types.ObjectType{ ... }
```
Add comments if the distinction between "Set" and "List" is important to clarify their use in the schema.
