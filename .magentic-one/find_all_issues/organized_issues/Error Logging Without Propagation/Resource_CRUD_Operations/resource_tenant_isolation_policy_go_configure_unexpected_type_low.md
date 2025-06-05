# Error Handling for Unexpected Resource Type in Configure

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy.go

## Problem

Within the `Configure` function, if the provider data is not of the expected type (`*api.ProviderClient`), an error is appended to diagnostics, but no further control flow actions are taken (such as aborting or returning early). Although the code later returns after appending the error, it is best practice to always follow an error-adding diagnostic block with an explicit `return` to avoid potential future code additions after the error handling that could cause logic bugs.

## Impact

- **Severity: Low**  
  The current design is safe as there is a `return` immediately after, but it is best practice for future maintainability to make error returns explicit after diagnostic error appends that are meant to halt processing.

## Location

```go
if req.ProviderData == nil {
	// ProviderData will be null when Configure is called from ValidateConfig.  It's ok.
	return
}

client, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
	resp.Diagnostics.AddError(
		"Unexpected Resource Configure Type",
		fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
	)
	return
}
```

## Code Issue

```go
client, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
	resp.Diagnostics.AddError(
		"Unexpected Resource Configure Type",
		fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
	)
	return
}
```

## Fix

You should continue always placing explicit `return` after `AddError` in error conditions, as is currently done. No modification is strictly needed, but developers should be aware of this flow and avoid code additions after error blocks.

```go
client, ok := req.ProviderData.(*api.ProviderClient)
if !ok {
	resp.Diagnostics.AddError(
		"Unexpected Resource Configure Type",
		fmt.Sprintf("Expected *api.ProviderClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
	)
	return // Explicit, to avoid unintended code execution after error
}
```
