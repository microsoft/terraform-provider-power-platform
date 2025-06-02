# Title

Error Handling Omits Logging and Telemetry

##
/workspaces/terraform-provider-power-platform/internal/services/dlp_policy/resource_dlp_policy.go

## Problem

Throughout error handling in the CRUD methods (`Read`, `Create`, `Update`, `Delete`), errors are added to `resp.Diagnostics` for user-facing reporting, but there is no consistent or explicit use of logging/telemetry hooks (e.g., via `tflog.Error` or similar) to capture these failures for operator or maintainer telemetry. This means provider developers and support engineers may miss error signals, delays in debugging or RCA, and lose important context for production issues.

## Impact

Severity: Medium

Medium supportability risk: the lack of back-end logs or provider telemetry makes diagnosing production or field issues slow and less efficient, particularly if dealing with cloud operational environments (i.e., managed Terraform Cloud or Enterprise).

## Location

```go
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err_client.Error())
	return
}
```

## Code Issue

```go
if err_client != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err_client.Error())
	return
}
```

(similar in Read, Update, Delete)

## Fix

Add explicit logging before adding diagnostics, e.g.

```go
if err_client != nil {
	tflog.Error(ctx, fmt.Sprintf("Create error: %s", err_client.Error()))
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when creating %s", r.FullTypeName()), err_client.Error())
	return
}
```
Or, standardize `tflog.Error` calls alongside diagnostic errors throughout the resource.

---
Save as:  
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/error_handling/resource_dlp_policy_missing_tflog_medium.md
