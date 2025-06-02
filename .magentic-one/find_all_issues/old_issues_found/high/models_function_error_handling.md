# Title

Improper Error Handling in `ConvertFromPowerAppDto` function.

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/models.go

## Problem

The function `ConvertFromPowerAppDto` does not handle errors or nullability of fields within `powerAppDto`. If any of the properties within `powerAppDto` (such as `powerAppDto.Properties.Environment.Name`, `powerAppDto.Properties.DisplayName`, etc.) are null or invalid, this could lead to runtime errors or unexpected behavior.

## Impact

Lack of error handling could cause runtime panics when null or invalid field values are encountered in `powerAppDto`. This impacts the program's stability and robustness.

Severity: **High**.

## Location

File: `models.go`

```go
func ConvertFromPowerAppDto(powerAppDto powerAppBapiDto) EnvironmentPowerAppsDataSourceModel {
	return EnvironmentPowerAppsDataSourceModel{
		EnvironmentId: types.StringValue(powerAppDto.Properties.Environment.Name),
		DisplayName:   types.StringValue(powerAppDto.Properties.DisplayName),
		Name:          types.StringValue(powerAppDto.Name),
		CreatedTime:   types.StringValue(powerAppDto.Properties.CreatedTime),
	}
}
```

## Code Issue

```go
func ConvertFromPowerAppDto(powerAppDto powerAppBapiDto) EnvironmentPowerAppsDataSourceModel {
	return EnvironmentPowerAppsDataSourceModel{
		EnvironmentId: types.StringValue(powerAppDto.Properties.Environment.Name),
		DisplayName:   types.StringValue(powerAppDto.Properties.DisplayName),
		Name:          types.StringValue(powerAppDto.Name),
		CreatedTime:   types.StringValue(powerAppDto.Properties.CreatedTime),
	}
}
```

## Fix

Introduce null checks and error handling for each field within `powerAppDto`. This ensures that the conversion function does not fail due to a missing or invalid value.

```go
func ConvertFromPowerAppDto(powerAppDto powerAppBapiDto) (EnvironmentPowerAppsDataSourceModel, error) {
	if powerAppDto.Properties.Environment.Name == "" {
		return EnvironmentPowerAppsDataSourceModel{}, fmt.Errorf("Environment Name is missing in powerAppDto")
	}
	if powerAppDto.Properties.DisplayName == "" {
		return EnvironmentPowerAppsDataSourceModel{}, fmt.Errorf("Display Name is missing in powerAppDto")
	}
	if powerAppDto.Name == "" {
		return EnvironmentPowerAppsDataSourceModel{}, fmt.Errorf("Power App Name is missing in powerAppDto")
	}
	if powerAppDto.Properties.CreatedTime == "" {
		return EnvironmentPowerAppsDataSourceModel{}, fmt.Errorf("Power App Created Time is missing in powerAppDto")
	}
	
	return EnvironmentPowerAppsDataSourceModel{
		EnvironmentId: types.StringValue(powerAppDto.Properties.Environment.Name),
		DisplayName:   types.StringValue(powerAppDto.Properties.DisplayName),
		Name:          types.StringValue(powerAppDto.Name),
		CreatedTime:   types.StringValue(powerAppDto.Properties.CreatedTime),
	}, nil
}
```

Explanation:
This fix ensures that nullability and invalid field values in `powerAppDto` are caught early with appropriate error messages, preventing runtime crashes and improving code reliability.