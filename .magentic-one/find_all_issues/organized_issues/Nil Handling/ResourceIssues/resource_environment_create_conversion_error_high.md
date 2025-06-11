# Title

Error Handling: Error Return Not Checked After Conversion Call in `Create`

##

/workspaces/terraform-provider-power-platform/internal/services/environment/resource_environment.go

## Problem

In the `Create` function, after calling `convertCreateEnvironmentDtoFromSourceModel`, the returned `err` is not followed by an immediate return when not `nil`. The error is logged to diagnostics, but execution continues, leading to potential use of an incomplete or nil `envToCreate` object downstream (for `LocationValidator`), which could result in a panic or unintended behavior.

## Impact

- **Severity**: High  
- **Explanation**: Could cause unexpected panics or access to nil pointers, downstream errors, or undefined code execution paths.

## Location

```go
envToCreate, err := convertCreateEnvironmentDtoFromSourceModel(ctx, plan, r)

if err != nil {
	resp.Diagnostics.AddError("Error when converting source model to create environment dto", err.Error())
}

// No return after AddError - code continues and envToCreate can be nil/garbage for following code
err = r.EnvironmentClient.LocationValidator(ctx, envToCreate.Location, envToCreate.Properties.AzureRegion)
```

## Code Issue

```go
envToCreate, err := convertCreateEnvironmentDtoFromSourceModel(ctx, plan, r)

if err != nil {
	resp.Diagnostics.AddError("Error when converting source model to create environment dto", err.Error())
}

// CONTINUES EVEN IF ERROR, POSSIBLE PANIC ON NEXT LINE.
err = r.EnvironmentClient.LocationValidator(ctx, envToCreate.Location, envToCreate.Properties.AzureRegion)
```

## Fix

Add an explicit `return` after logging the error to halt execution if conversion failed.

```go
envToCreate, err := convertCreateEnvironmentDtoFromSourceModel(ctx, plan, r)
if err != nil {
	resp.Diagnostics.AddError("Error when converting source model to create environment dto", err.Error())
	return // <- Add this to prevent further execution with invalid data
}
```

---
