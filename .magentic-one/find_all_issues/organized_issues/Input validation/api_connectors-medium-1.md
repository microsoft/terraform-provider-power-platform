# Inefficient Double For-Loop for Matching Connectors

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/api_connectors.go

## Problem

The following code block uses a double for-loop to find and update connectors with matching IDs for the `Unblockable` property:

```go
for inx, connector := range connectorArray.Value {
	for _, unblockableConnector := range unblockableConnectorArray {
		if connector.Id == unblockableConnector.Id {
			connectorArray.Value[inx].Properties.Unblockable = unblockableConnector.Metadata.Unblockable
		}
	}
}
```

Since both slices could be sizable, this approach is O(n*m) and is inefficient for large input sizes.

## Impact

For large lists of connectors, this significantly slows the execution, impacting performance. The severity is **medium** as it doesn't break functionality, but can degrade user experience or increase resource utilization.

## Location

Lines inside the `GetConnectors` method, after fetching both `connectorArray` and `unblockableConnectorArray` (first for-loop assignment to `inx` and `connector`).

## Code Issue

```go
for inx, connector := range connectorArray.Value {
	for _, unblockableConnector := range unblockableConnectorArray {
		if connector.Id == unblockableConnector.Id {
			connectorArray.Value[inx].Properties.Unblockable = unblockableConnector.Metadata.Unblockable
		}
	}
}
```

## Fix

Use a map to reduce lookup time for `unblockableConnector.Id`:

```go
// Build a map for fast lookup
unblockableMap := make(map[string]bool)
for _, uc := range unblockableConnectorArray {
	unblockableMap[uc.Id] = uc.Metadata.Unblockable
}

for inx, connector := range connectorArray.Value {
	if unblockable, ok := unblockableMap[connector.Id]; ok {
		connectorArray.Value[inx].Properties.Unblockable = unblockable
	}
}
```

This fix reduces the time complexity to O(n+m).
