# Inconsistent Naming: `Resource` Type Name is Too Generic

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/models.go

## Problem

The type `Resource` defined in this package is too generic and can be confusing. It does not reflect the context of the tenant isolation policy domain in which it is used. In Go, it is best practice to use more descriptive struct names, especially for exported types, to avoid ambiguity and improve maintainability. 

## Impact

A generic name like `Resource` introduces the risk of confusion when the package is imported elsewhere or when searching through code, especially as many plugins or providers use similar concepts (resources, clients, services, etc.). Severity: **medium** due to potential for maintenance problems in larger codebases or in code reviews.

## Location

```go
// Resource represents the tenant isolation policy resource.
type Resource struct {
	helpers.TypeInfo
	Client Client
}
```

## Fix

Rename the struct to something more specific, e.g., `TenantIsolationPolicyResource`, and update all usages accordingly.

```go
// TenantIsolationPolicyResource represents the tenant isolation policy resource.
type TenantIsolationPolicyResource struct {
	helpers.TypeInfo
	Client Client
}
```

This change improves clarity and discoverability of the type throughout your codebase.
