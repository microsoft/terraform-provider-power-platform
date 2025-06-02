# Diagnostics mishandling on Delete (cannot append Diagnostics to Diagnostics object)

##

/workspaces/terraform-provider-power-platform/internal/services/tenant_settings/resource_tenant_settings.go

## Problem

In the `Delete` method, the following code is used:

```go
previousBytes, diag := req.Private.GetKey(ctx, "original_settings")
if diag.HasError() {
    diag.Append(diag...)
    return
}
```
`diag` is a `tfsdk.Diagnostics` object. You are attempting to append diagnostics to itself (`diag.Append(diag...)`). However, `Append` is used to append a `Diagnostics` object to a different `Diagnostics` object. This does not write to `resp.Diagnostics` (the only diagnostics output that is considered by Terraform). As such, errors captured in `diag` are not presented to the end user in the Terraform output and may go unnoticed.

## Impact

Improper error handling and diagnostic output for users and maintainers. Errors may go unreported or make debugging more difficult. Severity: medium.

## Location

Function `Delete`; line with `diag.Append(diag...)`.

## Code Issue

```go
previousBytes, diag := req.Private.GetKey(ctx, "original_settings")
if diag.HasError() {
    diag.Append(diag...)
    return
}
```

## Fix

Append `diag` to `resp.Diagnostics`, not to itself. Only `resp.Diagnostics` is output. For example:

```go
previousBytes, diag := req.Private.GetKey(ctx, "original_settings")
if diag.HasError() {
    resp.Diagnostics.Append(diag...)
    return
}
```
Or, if `Append` is variadic:

```go
if diag.HasError() {
    resp.Diagnostics.Append(diag...)
    return
}
```

---

**This feedback will be saved under:**  
`/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/resource_tenant_settings.go-diagnostics_handling-medium.md`
