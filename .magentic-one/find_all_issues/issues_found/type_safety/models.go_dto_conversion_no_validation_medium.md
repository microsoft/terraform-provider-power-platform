# Title

Lack of Type or Value Validation in DTO–Model Conversion

##

/workspaces/terraform-provider-power-platform/internal/services/solution/models.go

## Problem

The `convertFromSolutionDto` function creates a `DataSourceModel` directly from the `SolutionDto` fields using the `types.StringValue` and `types.BoolValue` constructors, but does not provide any validation. If any fields in `solutionDto` are missing, empty, or have incorrect formats, the model could end up with invalid data.

## Impact

Medium. Without validation, downstream code may receive and operate on invalid or inconsistent state, increasing risk of bugs and unexpected failures further along the stack.

## Location

- `convertFromSolutionDto(solutionDto SolutionDto) DataSourceModel`

## Code Issue

```go
func convertFromSolutionDto(solutionDto SolutionDto) DataSourceModel {
	return DataSourceModel{
		EnvironmentId: types.StringValue(solutionDto.EnvironmentId),
		DisplayName:   types.StringValue(solutionDto.DisplayName),
		Name:          types.StringValue(solutionDto.Name),
		CreatedTime:   types.StringValue(solutionDto.CreatedTime),
		Id:            types.StringValue(solutionDto.Id),
		ModifiedTime:  types.StringValue(solutionDto.ModifiedTime),
		InstallTime:   types.StringValue(solutionDto.InstallTime),
		Version:       types.StringValue(solutionDto.Version),
		IsManaged:     types.BoolValue(solutionDto.IsManaged),
	}
}
```

## Fix

Add validation before constructing the model. Consider returning an error if mandatory fields are missing or malformed.

```go
func convertFromSolutionDto(solutionDto SolutionDto) (DataSourceModel, error) {
	if solutionDto.EnvironmentId == "" || solutionDto.Id == "" {
		return DataSourceModel{}, fmt.Errorf("required fields are missing: EnvironmentId or Id")
	}
	// Add further validation as needed.
	return DataSourceModel{
		EnvironmentId: types.StringValue(solutionDto.EnvironmentId),
		DisplayName:   types.StringValue(solutionDto.DisplayName),
		Name:          types.StringValue(solutionDto.Name),
		CreatedTime:   types.StringValue(solutionDto.CreatedTime),
		Id:            types.StringValue(solutionDto.Id),
		ModifiedTime:  types.StringValue(solutionDto.ModifiedTime),
		InstallTime:   types.StringValue(solutionDto.InstallTime),
		Version:       types.StringValue(solutionDto.Version),
		IsManaged:     types.BoolValue(solutionDto.IsManaged),
	}, nil
}
```
