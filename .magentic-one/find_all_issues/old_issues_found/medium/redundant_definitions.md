# Title

Redundant ObjectType Definitions

##

`/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/models.go`

## Problem

In the file, the `ObjectType` definitions such as `customConnectorPatternSetObjectType`, `connectorSetObjectType`, `endpointRuleListObjectType`, and `actionRuleListObjectType` are redundant because these mappings are already available in the form of structs or could be handled directly by the framework.

For example, instead of creating an `ObjectType` like:

```go
var customConnectorPatternSetObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"order":            types.Int64Type,
		"host_url_pattern": types.StringType,
		"data_group":       types.StringType,
	},
}
```

You could rely on the already defined `dataLossPreventionPolicyResourceCustomConnectorPattern` struct.

## Impact
- **Code Duplication**: It causes unnecessary duplication of object definitions, making the code harder to maintain.
- **Maintenance Overhead**: Any change to the struct would require a change in the corresponding `ObjectType` definition, increasing the potential for errors.
- **Impact Level**: Medium.

## Location
### Examples of Redundant Definitions:
```go
var customConnectorPatternSetObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"order":            types.Int64Type,
		"host_url_pattern": types.StringType,
		"data_group":       types.StringType,
	},
}

var connectorSetObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"id":                           types.StringType,
		"default_action_rule_behavior": types.StringType,
		"action_rules":                 types.ListType{ElemType: actionRuleListObjectType},
		"endpoint_rules":               types.ListType{ElemType: endpointRuleListObjectType},
	},
}
```

## Fix
Use structs directly instead of defining redundant `ObjectType` entries:

### Before
```go
var customConnectorPatternSetObjectType = types.ObjectType{
	AttrTypes: map[string]attr.Type{
		"order":            types.Int64Type,
		"host_url_pattern": types.StringType,
		"data_group":       types.StringType,
	},
}
```

### After
Maintain and utilize the `dataLossPreventionPolicyResourceCustomConnectorPattern` struct instead. The single source of truth approach simplifies code management.

```go
type dataLossPreventionPolicyResourceCustomConnectorPattern struct {
	Order          types.Int64  `tfsdk:"order"`
	HostUrlPattern types.String `tfsdk:"host_url_pattern"`
	DataGroup      types.String `tfsdk:"data_group"`
}
```

## Action
Save the above issue markdown as a separate file under `/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/medium/models_redundant_definitions.md`.