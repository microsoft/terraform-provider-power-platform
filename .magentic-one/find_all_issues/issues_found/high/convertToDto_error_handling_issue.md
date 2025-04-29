# Title

Inefficient Error Handling in `convertToDto` Function

##

`/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/models.go`

## Problem

In the `convertToDto` function, the error handling for the elements conversion (`model.AllowedTenants.ElementsAs`) may return errors. However, the code immediately exits without logging or handling errors in a meaningful manner.

```go
diags.Append(model.AllowedTenants.ElementsAs(ctx, &tenantsModel, false)...) 
if diags.HasError() {
	return nil, diags
}
```

This abrupt exit does not provide visibility into why the error occurred, and it is difficult to debug or investigate further.

## Impact

The lack of detailed error handling impacts debugging, traceability, and maintainability of the function. This makes it harder for developers and operations teams to identify the root cause of failures. Severity: **High**

## Location

Line 28-31 in `convertToDto`

## Code Issue

```go
diags.Append(model.AllowedTenants.ElementsAs(ctx, &tenantsModel, false)...) 
if diags.HasError() {
	return nil, diags
}
```

## Fix

Introduce logging or a more granular error-handling mechanism to capture useful diagnostics before exiting due to errors:
```go
diags.Append(model.AllowedTenants.ElementsAs(ctx, &tenantsModel, false)...) 
if diags.HasError() {
	// Log the error or append to diagnostic output
	for _, diagnostic := diags {
		fmt.Printf("Error converting AllowedTenants elements: %s\n", diagnostic)
	}
	return nil, diags
}
"