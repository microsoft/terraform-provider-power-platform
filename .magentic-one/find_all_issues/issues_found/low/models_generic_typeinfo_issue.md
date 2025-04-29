# Title

Generic TypeInfo Without Clear Purpose

## 

/workspaces/terraform-provider-power-platform/internal/services/tenant/models.go

## Problem

The `helpers.TypeInfo` embedded in `DataSource` is ambiguous and lacks clear documentation or purpose within this context. Without comments, it is difficult to understand why it is present or what role it plays in this struct.

## Impact

- This reduces code maintainability and readability.
- It can cause confusion for developers who are unfamiliar with the struct's purpose.

**Severity: Low**

## Location

```go
type DataSource struct {
	helpers.TypeInfo
	TenantClient Client
}
```

## Fix

Add well-defined comments or documentation indicating the purpose of `helpers.TypeInfo` in this struct. If `TypeInfo` is not necessary, consider removing it.

```go
type DataSource struct {
	// TypeInfo provides meta-level information required for tenant data sources.
	helpers.TypeInfo
	// TenantClient is responsible for handling interactions with the tenant API.
	TenantClient Client
}
```