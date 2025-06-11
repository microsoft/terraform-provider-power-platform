# Non-idiomatic Resource Struct Naming

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy.go

## Problem

The main resource struct is called `Resource` instead of something more specific such as `TenantIsolationPolicyResource`. Generic naming can decrease the readability and maintainability of the code and can also introduce type collisions in future code development (or test code).

## Impact

- **Severity: Low/Medium**  
  While it does not break the code, it makes it harder to understand and work with, especially as the codebase grows and more resources are introduced.

## Location

```go
var _ resource.Resource = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}
var _ resource.ResourceWithValidateConfig = &Resource{}
// ...
func NewTenantIsolationPolicyResource() resource.Resource {
	return &Resource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "tenant_isolation_policy",
		},
	}
}
```

## Code Issue

```go
var _ resource.Resource = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}
var _ resource.ResourceWithValidateConfig = &Resource{}

func NewTenantIsolationPolicyResource() resource.Resource {
	return &Resource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "tenant_isolation_policy",
		},
	}
}
// ...
type Resource struct { ... }
```

## Fix

Rename the struct to something more specific like `TenantIsolationPolicyResource`, and update all references for clarity and maintainability.

```go
var _ resource.Resource = &TenantIsolationPolicyResource{}
var _ resource.ResourceWithImportState = &TenantIsolationPolicyResource{}
var _ resource.ResourceWithValidateConfig = &TenantIsolationPolicyResource{}

func NewTenantIsolationPolicyResource() resource.Resource {
	return &TenantIsolationPolicyResource{
		TypeInfo: helpers.TypeInfo{
			TypeName: "tenant_isolation_policy",
		},
	}
}
// ...
type TenantIsolationPolicyResource struct {
	// add your fields here
}
```
