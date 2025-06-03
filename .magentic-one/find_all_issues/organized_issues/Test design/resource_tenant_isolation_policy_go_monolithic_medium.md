# Excessive Logic in Single File (Structure/Maintainability)

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy.go

## Problem

This file implements all resource behaviors and custom conversion/validation logic in a single file, which leads to a monolithic structure. This violates separation of concerns, making the file more difficult to read, maintain, and test. Examples include custom conversion functions (`convertToDto`, `convertFromDto`), complex diagnostics/error handling, and multiple responsibilities intermixed (schema definition, CRUD, validation, etc.).

## Impact

- **Severity: Medium**
- Hinders readability and collaborative development.
- Increases risk of merge conflicts.
- Makes future extension, troubleshooting, and testability more difficult.

## Location

```go
// Creation, validation, CRUD, diagnostic error branches, custom validators, conversion helpers, etc. are all present in this single file.
```

## Code Issue

```go
// Functions like convertToDto, convertFromDto, custom diagnostics, error reporting,
// schema construction, and all CRUD implementation appear in the same file.
```

## Fix

Refactor by extracting:

- DTO/model conversion logic (e.g., `convertToDto`, `convertFromDto`) to a separate file, such as `conversion.go` or `model.go`.
- Validation helpers out of the main resource file.
- Consider breaking up CRUD logic (Create, Update, Delete, etc.) to their own files if they grow or to group related logic together.
- Keep the core resource and schema/metadata in the main file, with helpers/factories elsewhere.

```go
// resource_tenant_isolation_policy.go:
//
// package tenant_isolation_policy
//
// ...core resource definition, schema, only minimal logic...
//
// conversion.go:
//
// package tenant_isolation_policy
//
// func convertToDto(ctx context.Context, id string, plan *TenantIsolationPolicyResourceModel) (TenantIsolationPolicyDto, diag.Diagnostics) { ... }
// func convertFromDto(ctx context.Context, dto TenantIsolationPolicyDto) (TenantIsolationPolicyResourceModel, diag.Diagnostics) { ... }
//
// validator.go:
//
// package tenant_isolation_policy
//
// func validateAllowedTenants(tenants []AllowedTenantModel) error { ... }
```
