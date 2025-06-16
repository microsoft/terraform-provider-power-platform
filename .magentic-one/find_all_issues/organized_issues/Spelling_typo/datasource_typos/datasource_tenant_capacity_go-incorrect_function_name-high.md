# Title

Incorrect function name in constructor

##

/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity.go

## Problem

The constructor function for the data source is named `NewTenantCapcityDataSource`, which appears to be a typo – it should likely be `NewTenantCapacityDataSource` (missing the "a" in "Capacity").

## Impact

Could cause confusion and reduce codebase maintainability. This is a high-severity naming issue because consumers of the function may not expect the typo and it could lead to errors or misunderstandings.

## Location

Line 19

## Code Issue

```go
func NewTenantCapcityDataSource() datasource.DataSource {
    return &DataSource{
        TypeInfo: helpers.TypeInfo{
            TypeName: "tenant_capacity",
        },
    }
}
```

## Fix

Rename the function to correct the typo:

```go
func NewTenantCapacityDataSource() datasource.DataSource {
    return &DataSource{
        TypeInfo: helpers.TypeInfo{
            TypeName: "tenant_capacity",
        },
    }
}
```

This helps ensure clarity and correct usage in the codebase.
