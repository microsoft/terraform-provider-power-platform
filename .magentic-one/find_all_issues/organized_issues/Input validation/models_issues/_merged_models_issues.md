# Models Issues - Input Validation

This document contains all models-related input validation issues found in the terraform-provider-power-platform codebase.


## ISSUE 1

# Missing Error Handling in Data Conversion

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/models.go

## Problem

The function `ConvertFromPowerAppDto` directly sets all fields using values, but does not handle possible missing/null/nil scenarios (e.g., if `Properties.Environment` or `Properties` is nil, a panic will occur). Robust error handling or validation is missing.

## Impact

Lack of error handling can lead to runtime panics and crashes if the incoming DTO is incomplete or malformatted (severity: high).

## Location

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
	return EnvironmentPowerAppsDataSourceModel{
		EnvironmentId: types.StringValue(powerAppDto.Properties.Environment.Name),
		DisplayName:   types.StringValue(powerAppDto.Properties.DisplayName),
		Name:          types.StringValue(powerAppDto.Name),
		CreatedTime:   types.StringValue(powerAppDto.Properties.CreatedTime),
	}
```

## Fix

Add validation or nil-checks before dereferencing nested fields. For example:

```go
func ConvertFromPowerAppDto(dto powerAppBapiDto) EnvironmentPowerAppsDataSourceModel {
	envID := ""
	displayName := ""
	createdTime := ""
	if dto.Properties != nil {
		if dto.Properties.Environment != nil {
			envID = dto.Properties.Environment.Name
		}
		displayName = dto.Properties.DisplayName
		createdTime = dto.Properties.CreatedTime
	}
	return EnvironmentPowerAppsDataSourceModel{
		EnvironmentId: types.StringValue(envID),
		DisplayName:   types.StringValue(displayName),
		Name:          types.StringValue(dto.Name),
		CreatedTime:   types.StringValue(createdTime),
	}
}
```


## ISSUE 2

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


## ISSUE 3

# Type Safety: Absence of Validation for DTO Fields May Cause Inconsistent State

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/models.go

## Problem

In functions like `convertFromDto` and `convertAllowedTenantsFromDto`, values from external DTOs are translated directly into model objects without any explicit validation for constraints, such as the format or presence of required fields (other than a check for empty tenant IDs in one case). This approach risks type safety and data consistency, especially since DTOs might originate from untrusted sources or may change in structure over time.

## Impact

Although Go is statically typed, this design opens the risk of silent data inconsistencies entering the model layer, which could go undetected until much later (such as during resource apply/update/create cycles). 
Severity: **medium** — there is potential for hard-to-debug issues or even downstream panics if invalid or partial data is converted without proper checks.

## Location

```go
func convertFromDto(ctx context.Context, dto *TenantIsolationPolicyDto) (TenantIsolationPolicyResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	if dto == nil {
		return TenantIsolationPolicyResourceModel{}, diags
	}

	// Set defaults in case Properties is nil
	tenantId := ""
	var isDisabled *bool
	var allowedTenants []AllowedTenantDto

	if dto.Properties.TenantId != "" {
		tenantId = dto.Properties.TenantId
	}
	if dto.Properties.IsDisabled != nil {
		isDisabled = dto.Properties.IsDisabled
	}
	if dto.Properties.AllowedTenants != nil {
		allowedTenants = dto.Properties.AllowedTenants
	}
	// ...
}
```
And:

```go
func convertAllowedTenantsFromDto(dtoTenants []AllowedTenantDto) []AllowedTenantModel {
	if dtoTenants == nil {
		return []AllowedTenantModel{}
	}

	modelTenants := make([]AllowedTenantModel, 0, len(dtoTenants))
	for _, dtoTenant := range dtoTenants {
		inbound := false
		outbound := false

		if dtoTenant.Direction.Inbound != nil {
			inbound = *dtoTenant.Direction.Inbound
		}
		if dtoTenant.Direction.Outbound != nil {
			outbound = *dtoTenant.Direction.Outbound
		}

		// Skip tenants with empty IDs
		if dtoTenant.TenantId == "" {
			continue
		}

		// Create a consistent model from the DTO with all fields explicitly set
		modelTenants = append(modelTenants, AllowedTenantModel{
			TenantId: types.StringValue(dtoTenant.TenantId),
			Inbound:  types.BoolValue(inbound),
			Outbound: types.BoolValue(outbound),
		})
	}
	return modelTenants
}
```

## Fix

Add input validation checks to ensure data consistency and robustness. For example:

```go
// In convertFromDto, add checks like:
if tenantId == "" {
    diags.AddError("Missing tenant ID", "DTO must provide a tenant ID")
    return TenantIsolationPolicyResourceModel{}, diags
}
// Add similar checks for other required fields as needed.
```

And, in `convertAllowedTenantsFromDto`, consider adding more checks for the `Direction` subobject, e.g.:
```go
if dtoTenant.Direction == nil {
    // Either skip or flag as error
    continue
}
```

Additionally, consider using helper validation functions for centralizing and reusing validation logic.


## ISSUE 4

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


## ISSUE 5

# Potential Type Safety Issue with Slices of Unvalidated Models

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/models.go

## Problem

The `SecurityRoles` field in `SecurityRolesListDataSourceModel` is a slice of `SecurityRoleDataSourceModel`, but there is no indication of validation or input sanitation for its elements. This could allow invalid or incomplete `SecurityRoleDataSourceModel` entries, as the struct is entirely public and its members are not protected.

## Impact

Inadequate validation can introduce subtle bugs, allow propagation of invalid state throughout the system, and potentially trigger failures downstream. Severity: **medium**.

## Location

```go
type SecurityRolesListDataSourceModel struct {
	Timeouts       timeouts.Value                `tfsdk:"timeouts"`
	EnvironmentId  types.String                  `tfsdk:"environment_id"`
	BusinessUnitId types.String                  `tfsdk:"business_unit_id"`
	SecurityRoles  []SecurityRoleDataSourceModel `tfsdk:"security_roles"`
}
```

## Code Issue

```go
SecurityRoles  []SecurityRoleDataSourceModel `tfsdk:"security_roles"`
```

## Fix

Introduce input validation either during assignment or by providing a constructor or validation method for `SecurityRolesListDataSourceModel` that ensures all required fields in its `SecurityRoleDataSourceModel` elements are set and valid.

```go
func (s *SecurityRolesListDataSourceModel) Validate() error {
    for i, role := range s.SecurityRoles {
        if role.RoleId.IsNull() || role.Name.IsNull() {
            return fmt.Errorf("Security role at index %d is missing required fields", i)
        }
    }
    return nil
}
```



## ISSUE 6

# Issue 1: Lack of Input Validation in DTO Creation

##

/workspaces/terraform-provider-power-platform/internal/services/copilot_studio_application_insights/models.go

## Problem

The function `createAppInsightsConfigDtoFromSourceModel` directly converts values from the `ResourceModel` to the DTO without any validation for required fields (such as `EnvironmentId`, `BotId`, or `AppInsightsConnectionString`). If any of these fields are empty or invalid, the DTO may be created with incomplete or invalid data, potentially causing runtime errors further along in the workflow.

## Impact

Severity: **High**

Unvalidated input may lead to the creation of DTOs with missing or malformed data, which can cause downstream API errors, unexpected behavior, or subtle bugs that are difficult to trace.

## Location

```go
func createAppInsightsConfigDtoFromSourceModel(appInsightsConfigSource ResourceModel) (*CopilotStudioAppInsightsDto, error) {
	appInsightsConfigDto := &CopilotStudioAppInsightsDto{
		EnvironmentId:               appInsightsConfigSource.EnvironmentId.ValueString(),
		BotId:                       appInsightsConfigSource.BotId.ValueString(),
		AppInsightsConnectionString: appInsightsConfigSource.AppInsightsConnectionString.ValueString(),
		IncludeSensitiveInformation: appInsightsConfigSource.IncludeSensitiveInformation.ValueBool(),
		IncludeActivities:           appInsightsConfigSource.IncludeActivities.ValueBool(),
		IncludeActions:              appInsightsConfigSource.IncludeActions.ValueBool(),
		NetworkIsolation:            "PublicNetwork",
	}

	return appInsightsConfigDto, nil
}
```

## Fix

Add validation to ensure all required fields are present and valid before creating the DTO. Return an error if validation fails.

```go
func createAppInsightsConfigDtoFromSourceModel(appInsightsConfigSource ResourceModel) (*CopilotStudioAppInsightsDto, error) {
	envId := appInsightsConfigSource.EnvironmentId.ValueString()
	botId := appInsightsConfigSource.BotId.ValueString()
	connStr := appInsightsConfigSource.AppInsightsConnectionString.ValueString()

	if envId == "" {
		return nil, fmt.Errorf("EnvironmentId cannot be empty")
	}
	if botId == "" {
		return nil, fmt.Errorf("BotId cannot be empty")
	}
	if connStr == "" {
		return nil, fmt.Errorf("ApplicationInsightsConnectionString cannot be empty")
	}

	appInsightsConfigDto := &CopilotStudioAppInsightsDto{
		EnvironmentId:               envId,
		BotId:                       botId,
		AppInsightsConnectionString: connStr,
		IncludeSensitiveInformation: appInsightsConfigSource.IncludeSensitiveInformation.ValueBool(),
		IncludeActivities:           appInsightsConfigSource.IncludeActivities.ValueBool(),
		IncludeActions:              appInsightsConfigSource.IncludeActions.ValueBool(),
		NetworkIsolation:            "PublicNetwork",
	}

	return appInsightsConfigDto, nil
}
```


# To finish the task you have to 
1. Run linter and fix any issues 
2. Run UnitTest and fix any of failing ones
3. Generate docs 
4. Run Changie

# Changie Instructions
Create only one change log entry. Do not run the changie tool multiple times.

```bash
changie new --kind <kind_key> --body "<description>" --custom Issue=<issue_number>
```
Where:
- `<kind_key>` is one of: breaking, changed, deprecated, removed, fixed, security, documentation
- `<description>` is a clear explanation of what was fixed/changed search for 'copilot-commit-message-instructions.md' how to write description.
- `<issue_number>` pick the issue number or PR number
