# Inconsistent receiver and exported type naming for DataSource struct

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/datasource_tenant.go

## Problem

The struct type `DataSource` is used without a package prefix, and the receiver is a plain single letter `d`. In the context of a provider, this can cause confusion and reduce readability, especially when browsing with an IDE across multiple packages where `DataSource` is generic and ambiguous.

## Impact

Low, but may affect code navigation and clarity for contributors. Potentially increases friction for maintainers.

## Location

Throughout the file, e.g.

## Code Issue

```go
type DataSource struct { ... }
func (d *DataSource) Read(...) {...}
```

## Fix

Use a more specific struct name and receiver, such as:

```go
type TenantDataSource struct { ... }
func (tds *TenantDataSource) Read(...) {...}
```

Or ensure struct is always referenced with package prefix in documentation/comments.
