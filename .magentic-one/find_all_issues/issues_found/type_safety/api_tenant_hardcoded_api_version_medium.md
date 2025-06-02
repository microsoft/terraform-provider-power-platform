# Usage of Hardcoded API Version String

##

/workspaces/terraform-provider-power-platform/internal/services/tenant/api_tenant.go

## Problem

The function `GetTenant` contains a hardcoded API version string (`"2021-04-01"`). This can be error-prone when the API version changes in future, and makes it harder to update, maintain, and spot all usages.

## Impact

Hardcoded values decrease maintainability and risk subtle bugs if API versions need changing across the project. **Severity:** Medium

## Location

```go
values := url.Values{}
values.Add("api-version", "2021-04-01")
```

## Code Issue

```go
values := url.Values{}
values.Add("api-version", "2021-04-01")
```

## Fix

Define the API version as a constant, ideally in your `constants` package:

```go
// in constants.go
const TenantApiVersion = "2021-04-01"

// in api_tenant.go
values := url.Values{}
values.Add("api-version", constants.TenantApiVersion)
```
