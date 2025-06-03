# Misspelled Variable Name: "virutualConnector"

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/api_connectors.go

## Problem

The variable `virutualConnector` in this section is misspelled; it should be `virtualConnector` for consistency and clarity.

## Impact

Misspelled variable names decrease code readability, can cause confusion during maintenance or code reviews, and may introduce bugs if similarly-named variables are later introduced. Severity: **low**.

## Location

Within the `GetConnectors` method, this line:

```go
for _, virutualConnector := range virtualConnectorArray {
	connectorArray.Value = append(connectorArray.Value, connectorDto{
		Id:   virutualConnector.Id,
		Name: virutualConnector.Metadata.Name,
		Type: virutualConnector.Metadata.Type,
		Properties: connectorPropertiesDto{
			DisplayName: virutualConnector.Metadata.DisplayName,
			Unblockable: false,
			Tier:        "Built-in",
			Publisher:   "Microsoft",
			Description: "",
		},
	})
}
```

## Code Issue

```go
for _, virutualConnector := range virtualConnectorArray {
	...
	Id:   virutualConnector.Id,
	Name: virutualConnector.Metadata.Name,
	Type: virutualConnector.Metadata.Type,
	Properties: connectorPropertiesDto{
		DisplayName: virutualConnector.Metadata.DisplayName,
		...
	}
}
```

## Fix

Rename the variable for correct spelling:

```go
for _, virtualConnector := range virtualConnectorArray {
	connectorArray.Value = append(connectorArray.Value, connectorDto{
		Id:   virtualConnector.Id,
		Name: virtualConnector.Metadata.Name,
		Type: virtualConnector.Metadata.Type,
		Properties: connectorPropertiesDto{
			DisplayName: virtualConnector.Metadata.DisplayName,
			Unblockable: false,
			Tier:        "Built-in",
			Publisher:   "Microsoft",
			Description: "",
		},
	})
}
```
