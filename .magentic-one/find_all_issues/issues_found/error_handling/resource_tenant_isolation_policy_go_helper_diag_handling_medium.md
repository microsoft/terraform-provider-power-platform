# Unhandled Errors in Helper Function Calls

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_isolation_policy/resource_tenant_isolation_policy.go

## Problem

Helper functions such as `helpers.EnterRequestContext` and DTO conversion helpers (e.g., `convertToDto`, `convertFromDto`) may return errors (directly or in diagnostic collections) that are appended to `resp.Diagnostics`. However, after appending, the code only checks for errors in some places, sometimes proceeding even after potential diagnostics have been added.

## Impact

- **Severity: Medium**
- Operations may continue after diagnostic errors, potentially leading to nil dereference, invalid state or further unintended errors. The code should ensure any diagnostic append is followed by an explicit error check.

## Location

E.g., in Create:

```go
resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
if resp.Diagnostics.HasError() {
	return
}
// ...
policyDto, diags := convertToDto(ctx, tenantInfo.TenantId, &plan)
resp.Diagnostics.Append(diags...)
if resp.Diagnostics.HasError() {
	return
}
```

But in other places, sometimes errors are appended but not always followed up by an immediate HasError() check.

## Code Issue

```go
// ... 
resp.Diagnostics.Append(helperResult...)
if resp.Diagnostics.HasError() {
    return
}
// (good)
```
However, be sure this is done **everywhere** a diagnostic can be appended, not just in resource CRUD but also in config validation and import. Review for missing checks.

## Fix

After **every** `resp.Diagnostics.Append`, immediately check `resp.Diagnostics.HasError()` and return if any error is present, to prevent subsequent code from executing with a bad state.

```go
resp.Diagnostics.Append(someDiagFunc()...)
if resp.Diagnostics.HasError() {
    return // abort processing further in this handler
}
```
