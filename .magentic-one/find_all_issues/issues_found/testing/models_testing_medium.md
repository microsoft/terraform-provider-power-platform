# Lack of Unit Tests for Model Conversion Logic

##

/workspaces/terraform-provider-power-platform/internal/services/connectors/models.go

## Problem

There are no unit tests provided or indicated for the `convertFromConnectorDto` function, which is responsible for transforming external DTOs into the internal DataSourceModel. Conversion logic often harbors edge cases and should be independently tested to ensure correctness, especially regarding handling of nil values, missing fields, or unexpected input formats.

## Impact

Absence of unit tests can lead to undetected conversion bugs or regressions, particularly as requirements or DTOs evolve. This has a **medium severity**, as it impacts software reliability and future refactor safety.

## Location

- `convertFromConnectorDto`

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

Add targeted unit tests for this function (and similar helpers), covering typical, edge, and erroneous cases. Example in `_test.go` file:

```go
func TestConvertFromConnectorDto(t *testing.T) {
	dto := connectorDto{
		Id:   "id1",
		Name: "example",
		Type: "foo",
		Properties: &connectorProperties{
			Description: "desc",
			DisplayName: "display",
			Tier:        "basic",
			Publisher:   "publisher",
			Unblockable: true,
		},
	}
	result := convertFromConnectorDto(dto)
	if result.Id.ValueString() != "id1" {
		t.Errorf("expected 'id1', got %v", result.Id.ValueString())
	}
	// Additional assertions for other fields, including missing or nil Properties
}
```
