# Title

Function `convertToAttrValueConnectorsGroup` prematurely returns from loop

##

/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/helpers.go

## Problem

In the function `convertToAttrValueConnectorsGroup`, the function immediately returns a value upon finding the first connectors group with a matching classification. If there are multiple connector groups with the same classification, only the first one is included and the rest are ignored. This could lead to missed data if the input slice contains more than one group of a given classification.

## Impact

Medium severity. This might cause incomplete data to be returned if more than one connectors group with the same classification is present, and thus can lead to data loss or unexpected behavior.

## Location

Lines 80-87:

```go
func convertToAttrValueConnectorsGroup(classification string, connectorsGroup []dlpConnectorGroupsModelDto) basetypes.SetValue {
	var connectorValues []attr.Value
	for _, conn := range connectorsGroup {
		if conn.Classification == classification {
			return types.SetValueMust(connectorSetObjectType, convertToAttrValueConnectors(conn, connectorValues))
		}
	}
	return types.SetValueMust(connectorSetObjectType, []attr.Value{})
}
```

## Code Issue

```go
for _, conn := range connectorsGroup {
	if conn.Classification == classification {
		return types.SetValueMust(connectorSetObjectType, convertToAttrValueConnectors(conn, connectorValues))
	}
}
```

## Fix

Accumulate all matching connector groups and return them together. For example:

```go
func convertToAttrValueConnectorsGroup(classification string, connectorsGroup []dlpConnectorGroupsModelDto) basetypes.SetValue {
	var connectorValues []attr.Value
	for _, conn := range connectorsGroup {
		if conn.Classification == classification {
			connectorValues = append(connectorValues, convertToAttrValueConnectors(conn, []attr.Value{})...)
		}
	}
	return types.SetValueMust(connectorSetObjectType, connectorValues)
}
```
