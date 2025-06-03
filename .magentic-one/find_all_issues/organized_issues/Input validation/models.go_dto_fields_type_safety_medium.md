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
