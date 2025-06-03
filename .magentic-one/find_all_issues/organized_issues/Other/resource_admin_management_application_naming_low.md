# Function Name Should Follow Go Naming Conventions

##

/workspaces/terraform-provider-power-platform/internal/services/admin_management_application/resource_admin_management_application.go

## Problem

The name `NewAdminManagementApplicationResource` follows a Go convention but could be shortened for clarity. However, in the context of a Terraform provider, this full name is acceptable, but reviewing whether all exported types and functions have clear and consistent names is important.

Additionally, the struct field `AdminManagementApplicationClient` is verbose but matches type intent; however, it's worth considering shortening (such as `Client` or `AdminAppClient`) for brevity if consistently used.

## Impact

Severity: **low**

Minor impact; mostly concerns development ergonomics, onboarding, and code readability. Existing names are not misleading but could be less verbose.

## Location

```go
func NewAdminManagementApplicationResource() resource.Resource { ... }
// and usage of AdminManagementApplicationClient within the resource struct.
```

## Fix

Consider more concise but still descriptive naming for both function and struct fields, for example:

```go
func NewAdminAppResource() resource.Resource { ... }

// In struct:
AdminAppClient *adminManagementApplicationClient
```

This is optional but may improve readability over time.
