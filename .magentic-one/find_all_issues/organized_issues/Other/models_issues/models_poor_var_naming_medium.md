# Poor Variable Naming: `Resource` Type Is Too Generic

##

/workspaces/terraform-provider-power-platform/internal/services/environment/models.go

## Problem

The struct type `Resource` is named too generically, which is ambiguous and doesn't communicate what "resource" it represents in the context of Power Platform environments. In Go, general names like `Resource` can lead to confusion, make it harder to search/grep usage, and increase risk of name collision, especially in large codebases or when importing multiple packages.

## Impact

- **Severity:** Medium
- Reduces code readability and maintainability.
- Makes tracking resource-related issues difficult in a multi-resource provider context.
- May accidentally conflict with more specific resources (in future expansion).

## Location

```go
type Resource struct {
	helpers.TypeInfo
	EnvironmentClient Client
	LicensingClient   licensing.Client
}
```

## Code Issue

```go
type Resource struct {
	helpers.TypeInfo
	EnvironmentClient Client
	LicensingClient   licensing.Client
}
```

## Fix

Use a type name that describes its domain and responsibility, such as `EnvironmentResource`. This improves clarity, especially for teams handling multiple types of resources.

```go
type EnvironmentResource struct {
	helpers.TypeInfo
	EnvironmentClient Client
	LicensingClient   licensing.Client
}
```

Usage throughout the package should also be replaced with `EnvironmentResource`. This small change improves self-documentation and searchability.
