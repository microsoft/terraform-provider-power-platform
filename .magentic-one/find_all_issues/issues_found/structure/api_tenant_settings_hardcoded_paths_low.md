# Hardcoded API URL Paths and Methods

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/api_tenant_settings.go

## Problem

Throughout the file, API endpoint paths, HTTP methods, and versions are hardcoded as string literals in multiple functions (e.g., `"GET"`, `"POST"`, and `"/providers/Microsoft.BusinessAppPlatform/tenant"`). This duplication increases the risk of errors, makes future updates harder, and decreases overall maintainability. It is more maintainable and readable to define such constants centrally.

## Impact

- **Severity: Low**
- Increases maintenance effort if endpoints or methods ever change.
- Inconsistent usage can introduce subtle bugs.
- Makes the codebase harder to audit for required API version upgrades.

## Location

```go
// Example (snippets from multiple functions)
Path:   "/providers/Microsoft.BusinessAppPlatform/tenant",
...
Path:   "/providers/Microsoft.BusinessAppPlatform/listTenantSettings",
...
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &tenant)
...
_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, nil, []int{http.StatusOK}, &tenantSettings)
```

## Code Issue

```go
Path:   "/providers/Microsoft.BusinessAppPlatform/tenant",
...
Path:   "/providers/Microsoft.BusinessAppPlatform/listTenantSettings",
...
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &tenant)
...
_, err := client.Api.Execute(ctx, nil, "POST", apiUrl.String(), nil, nil, []int{http.StatusOK}, &tenantSettings)
```

## Fix

Define centralized constants for endpoint paths, method verbs, and API versions at the top of the file or in a constants package:

```go
const (
    getTenantPath             = "/providers/Microsoft.BusinessAppPlatform/tenant"
    listTenantSettingsPath    = "/providers/Microsoft.BusinessAppPlatform/listTenantSettings"
    updateTenantSettingsPath  = "/providers/Microsoft.BusinessAppPlatform/scopes/admin/updateTenantSettings"
    methodGet                 = "GET"
    methodPost                = "POST"
    apiVersion20200801        = "2020-08-01"
    apiVersion20230601        = "2023-06-01"
)
```

Use these constants in your function implementations, which will improve readability and maintainability.
