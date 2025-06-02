# API Client Interface Coupling

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/datasource_analytics_data_exports.go

## Problem

The code tightly couples the `AnalyticsExportDataSource` to a specific client constructor:

```go
d.analyticsExportClient = NewAnalyticsExportClient(client.Api, tenant.NewTenantClient(client.Api))
```

There is no interface abstraction or dependency injection, which makes unit testing difficult and couples the provider code directly to concrete implementations. In the future, swapping out the client (e.g., for testing or for alternate backends) would require editing multiple places in code.

## Impact

- Reduces code testability (harder to mock or stub API usage)
- Makes refactoring or future abstraction more difficult
- Makes it harder to reuse or extend this provider code

**Severity:** Medium

## Location

```go
d.analyticsExportClient = NewAnalyticsExportClient(client.Api, tenant.NewTenantClient(client.Api))
```

## Code Issue

```go
d.analyticsExportClient = NewAnalyticsExportClient(client.Api, tenant.NewTenantClient(client.Api))
```

## Fix

Define an interface for the analytics export client, accept that interface in the struct (or via a setter for testing), and configure with concrete implementation in production:

```go
type AnalyticsExportClient interface {
    GetAnalyticsDataExport(ctx context.Context) ([]AnalyticsDataDto, error)
    // ...other methods
}

// In the struct definition
analyticsExportClient AnalyticsExportClient

// In production:
d.analyticsExportClient = NewAnalyticsExportClient(client.Api, tenant.NewTenantClient(client.Api))

// For testing, inject a mock implementing AnalyticsExportClient.
```
