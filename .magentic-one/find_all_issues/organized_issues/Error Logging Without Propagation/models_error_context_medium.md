# Error Handling: Potential Loss of Context in `convertCreateEnvironmentDtoFromSourceModel`

##

/workspaces/terraform-provider-power-platform/internal/services/environment/models.go

## Problem

In `convertCreateEnvironmentDtoFromSourceModel`, if retrieving the tenant ID from `r.EnvironmentClient.tenantClient.GetTenant(ctx)` fails, the error is returned directly without additional context or logging, making it harder to diagnose the original call's context and origin during debugging.

## Impact

- **Severity:** Medium
- Makes tracing errors difficult, especially in complex workflows or remote API calls.
- May frustrate operators/users since log output may be unclear.
- May prevent proper root-cause diagnosis if similar errors occur in related library functions.

## Location

```go
if !environmentSource.OwnerId.IsNull() && !environmentSource.OwnerId.IsUnknown() {
	tenantId, err := r.EnvironmentClient.tenantClient.GetTenant(ctx)
	if err != nil {
		return nil, err
	}
	environmentDto.Properties.UsedBy = &UsedByDto{
		Id:       environmentSource.OwnerId.ValueString(),
		Type:     "1",
		TenantId: tenantId.TenantId,
	}
}
```

## Code Issue

```go
tenantId, err := r.EnvironmentClient.tenantClient.GetTenant(ctx)
if err != nil {
	return nil, err
}
```

## Fix

Wrap errors with context for better traceability, using `fmt.Errorf` or `%w`:

```go
tenantId, err := r.EnvironmentClient.tenantClient.GetTenant(ctx)
if err != nil {
	return nil, fmt.Errorf("failed to retrieve tenant for OwnerId %s: %w", environmentSource.OwnerId.ValueString(), err)
}
```

This will help operators quickly understand what failed if this surface-level call errors.

---

**This markdown should be saved as:**  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/models_error_context_medium.md`
