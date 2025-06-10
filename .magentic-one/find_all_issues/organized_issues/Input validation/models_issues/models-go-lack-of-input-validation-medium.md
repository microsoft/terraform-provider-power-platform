# Title

Lack of Input Validation When Mapping DTO Fields

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/models.go

## Problem

The function `convertDtoToModel` assumes all DTO string fields and slices contain valid and expected values (e.g., essentials like `dto.ID`, `dto.Source`, `dto.Environments`, etc.). There is no validation or sanity checking of incoming data for expected format, non-empty or valid values—particularly problematic when values are marshalled directly into resource state.

This could result in accidental propagation of invalid or unexpected values through the system, causing issues downstream.

## Impact

Severity: **Medium**

A lack of validation carries the risk of introducing invalid or inconsistent state into Terraform resources, which could propagate through to infrastructure deployments.

## Location

Within this mapping (and throughout the conversion function):

```go
	ID:           types.StringValue(dto.ID),
	Source:       types.StringValue(dto.Source),
	Environments: environments,
	Status:       status,
	...
```

## Fix

Sanitize and validate input where appropriate before converting to Terraform-compatible types. For example, check for required fields being empty and set a `types.StringNull()` or log a warning, as appropriate.

```go
	id := types.StringNull()
	if dto.ID != "" {
		id = types.StringValue(dto.ID)
	}
	// ...repeat for required fields

	return &AnalyticsDataModel{
		ID:           id,
		// ...
	}
```

A similar approach can be applied for other required fields to ensure data consistency.
