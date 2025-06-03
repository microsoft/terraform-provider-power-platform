# Issue: Failure to Log and Propagate Detailed Errors in Resource Methods

##

/workspaces/terraform-provider-power-platform/internal/services/enterprise_policy/resource_enterprise_policy.go

## Problem

In several resource CRUD methods (`Create`, `Read`, `Delete`), when an error is returned from the API client, the response is set by calling `resp.Diagnostics.AddError` with a basic error message and `err.Error()`. This practice can result in loss of valuable context for debugging, as the original error is stringified and additional context from error wrapping or stack traces is lost. There is also insufficient logging of errors via `tflog.Error`, which can aid debugging in production systems.

## Impact

- Loss of error context, making debugging harder.
- Less useful diagnostics for Terraform users.
- Missed opportunity for richer structured logging.
- **Severity:** Medium

## Location

```go
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
	return
}
...
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading environment in %s", r.FullTypeName()), err.Error())
	return
}
...
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s", r.FullTypeName()), err.Error())
	return
}
```

## Code Issue

```go
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err.Error())
	return
}
```

## Fix

- Log structured error with `tflog.Error`.
- Consider returning wrapped errors or more context-rich diagnostics.
- Example:

```go
if err != nil {
	tflog.Error(ctx, fmt.Sprintf("API client error during create for %s: %v", r.FullTypeName(), err))
	resp.Diagnostics.AddError(
		fmt.Sprintf("API client error when creating %s", r.FullTypeName()),
		fmt.Sprintf("Could not create resource.\nError: %v", err),
	)
	return
}
```
Do similar for `Read` and `Delete`.

---

This markdown will be saved under:

`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/resource_enterprise_policy_error_handling_medium_logging.md`
