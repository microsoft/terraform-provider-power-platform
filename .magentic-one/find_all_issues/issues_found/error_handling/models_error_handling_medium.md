# Missing Error Handling in `convertFromConnectorDto` Function

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/models.go

## Problem

The function `convertFromConnectorDto` assumes that all fields within `connectorDto` and its nested `Properties` are present and initialized. If any of these are missing, nil, or otherwise invalid, this could result in a runtime panic or default zero-value assignment, which may not be desired behavior.

## Impact

Lack of error handling makes the system susceptible to panics or undetected data inconsistencies, especially when data is received from external sources. This can cause subtle bugs or service crashes. Severity is **medium**, due to risk of runtime errors when code is exposed to incomplete or malformed input.

## Location

- `convertFromConnectorDto(connectorDto connectorDto) DataSourceModel`

## Code Issue

```go
func convertFromConnectorDto(connectorDto connectorDto) DataSourceModel {
	return DataSourceModel{
		Id:          types.StringValue(connectorDto.Id),
		Name:        types.StringValue(connectorDto.Name),
		Type:        types.StringValue(connectorDto.Type),
		Description: types.StringValue(connectorDto.Properties.Description),
		DisplayName: types.StringValue(connectorDto.Properties.DisplayName),
		Tier:        types.StringValue(connectorDto.Properties.Tier),
		Publisher:   types.StringValue(connectorDto.Properties.Publisher),
		Unblockable: types.BoolValue(connectorDto.Properties.Unblockable),
	}
}
```

## Fix

Add validation for presence of necessary fields and gracefully handle missing or nil sub-structs:

```go
func convertFromConnectorDto(connectorDto connectorDto) DataSourceModel {
	var description, displayName, tier, publisher string
	var unblockable bool

	if connectorDto.Properties != nil {
		description = connectorDto.Properties.Description
		displayName = connectorDto.Properties.DisplayName
		tier = connectorDto.Properties.Tier
		publisher = connectorDto.Properties.Publisher
		unblockable = connectorDto.Properties.Unblockable
	}

	return DataSourceModel{
		Id:          types.StringValue(connectorDto.Id),
		Name:        types.StringValue(connectorDto.Name),
		Type:        types.StringValue(connectorDto.Type),
		Description: types.StringValue(description),
		DisplayName: types.StringValue(displayName),
		Tier:        types.StringValue(tier),
		Publisher:   types.StringValue(publisher),
		Unblockable: types.BoolValue(unblockable),
	}
}
```

This avoids potential panics and increases the robustness of the code when handling untrusted input.
