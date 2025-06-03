# Variable Shadowing and Hidden Field Assigment

##

/workspaces/terraform-provider-power-platform/internal/services/environment_wave/resource_environment_wave.go

## Problem

In the `Metadata` method, the assignment `r.ProviderTypeName = req.ProviderTypeName` references a promoted field from `helpers.TypeInfo`, but this is not clearly visible to the reader unless they know the content of the embedded struct. This may reduce maintainability and make code harder to follow for those unfamiliar with the project’s struct embedding patterns.

## Impact

Low severity: code is correct, but the readability and maintainability for the codebase are affected, especially for those new to the codebase who may be surprised by embedded fields being used and set.

## Location

```go
func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	r.ProviderTypeName = req.ProviderTypeName
	...
}
```

## Code Issue

```go
r.ProviderTypeName = req.ProviderTypeName
```

## Fix

Add a line of documentation to clarify that `ProviderTypeName` is from the embedded `TypeInfo`, or explicitly reference the embedded struct as `r.TypeInfo.ProviderTypeName` for increased clarity:

```go
// r.ProviderTypeName is from embedded helpers.TypeInfo
r.ProviderTypeName = req.ProviderTypeName
```

Or:

```go
r.TypeInfo.ProviderTypeName = req.ProviderTypeName
```

**Explanation:**
- This increases codebase comprehensibility and helps future maintainers understand where fields are coming from, and avoids subtle bugs due to shadowing or accidental extra fields.