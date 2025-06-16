# Title

Possible schema and model divergence risk in mapping

##

/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity.go

## Problem

The mapping from API DTOs to the Terraform provider state model is done manually and field-by-field, which is standard. However, there is no explicit check or guarantee that the schema attributes, `DataSourceModel`, and the API DTO fields do not diverge over time. This has led to issues in similar projects where newly added fields in the API or schema can silently fail to propagate, causing hard-to-trace bugs.

## Impact

Medium severity. Code maintainability and reliability may suffer as APIs or provider schemas evolve. New required fields might be omitted in the state model, or deprecated fields might linger, both of which are subtle and hard to debug mistakes.

## Location

Function: `Read`, throughout state mapping

## Code Issue

```go
// Multiple field-by-field assignments from tenantCapacityDto to state,
// no defensive assertion or automated mapping style.
```

## Fix

Consider adding explicit comments, checks, or even a test to validate that the schema, the DTO, and the state model stay in sync. At a minimum, document the need to update all parts on change:

```go
// Ensure every change in schema or TenantCapacityDto is reflected in DataSourceModel/TenantCapacityDataSourceModel, and update this mapping accordingly.
```

Or implement code-generation/mapping tests, when possible.
