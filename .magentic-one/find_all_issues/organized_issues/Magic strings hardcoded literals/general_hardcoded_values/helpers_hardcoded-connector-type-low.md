# Title

Hardcoded string values for connector type

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/helpers.go

## Problem

The connector `Type` is hardcoded as `"Microsoft.PowerApps/apis"` in multiple places in the code rather than being extracted from a constant or from configuration. Magic strings in code can lead to spelling inconsistencies and make it harder to update or maintain the namespace/type in the future.

## Impact

Low severity. This is a maintainability and reliability problem—if the connector type ever needs to change, every occurrence has to be manually updated, increasing the risk of mistakes.

## Location

There are multiple locations, for example:

Lines 112–121 (`getConnectorGroup`):

```go
for _, connector := range connectors {
	connectorGroup.Connectors = append(connectorGroup.Connectors, dlpConnectorModelDto{
		Id:                        connector.Id.ValueString(),
		Type:                      "Microsoft.PowerApps/apis", // <--- hardcoded
		DefaultActionRuleBehavior: connector.DefaultActionRuleBehavior.ValueString(),
		ActionRules:               convertToDlpActionRule(connector),
		EndpointRules:             convertToDlpEndpointRule(connector),
	})
}
```

Lines 137 and 147 (`convertToDlpConnectorGroup`):

```go
connectorGroup.Connectors = append(connectorGroup.Connectors, dlpConnectorModelDto{
	Id:   connector.Id.ValueString(),
	Type: "Microsoft.PowerApps/apis", // <--- hardcoded
	DefaultActionRuleBehavior: defaultAction,
	ActionRules:               convertToDlpActionRule(connector),
	EndpointRules:             convertToDlpEndpointRule(connector),
})
```

## Code Issue

```go
Type: "Microsoft.PowerApps/apis",
```

## Fix

Extract the type string into a package-level constant, and use that constant throughout the codebase. For example:

```go
const connectorTypePowerApps = "Microsoft.PowerApps/apis"
```

And then in all code locations:

```go
Type: connectorTypePowerApps,
```
