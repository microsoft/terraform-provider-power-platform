# Title

Missing Error Handling for `convertFromConnectorDto` Function

##

`/workspaces/terraform-provider-power-platform/internal/services/connectors/models.go`

## Problem

The `convertFromConnectorDto` function does not perform any error handling or validation on the input `connectorDto`. If the struct `connectorDto` contains `nil` or unexpected values, the function could propagate errors silently into the `DataSourceModel`.

## Impact

Without validation, invalid data may enter the system, potentially causing unexpected behavior or panics elsewhere in the code. This is a **medium** severity issue since it could lead to subtle bugs or runtime errors.

## Location

`convertFromConnectorDto` function, lines 31-45.

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

Add validation checks for `connectorDto` fields before attempting to convert them into `DataSourceModel`. For instance:

```go
func convertFromConnectorDto(connectorDto connectorDto) (DataSourceModel, error) {
	if connectorDto == nil {
		return DataSourceModel{}, fmt.Errorf("connectorDto is nil")
	}

	if connectorDto.Id == "" || connectorDto.Name == "" || connectorDto.Type == "" {
		return DataSourceModel{}, fmt.Errorf("connectorDto contains empty critical fields")
	}

	return DataSourceModel{
		Id:          types.StringValue(connectorDto.Id),
		Name:        types.StringValue(connectorDto.Name),
		Type:        types.StringValue(connectorDto.Type),
		Description: types.StringValue(connectorDto.Properties.Description),
		DisplayName: types.StringValue(connectorDto.Properties.DisplayName),
		Tier:        types.StringValue(connectorDto.Properties.Tier),
		Publisher:   types.StringValue(connectorDto.Properties.Publisher),
		Unblockable: types.BoolValue(connectorDto.Properties.Unblockable),
	}, nil
}
```

This ensures that the function behaves predictably when given invalid data.