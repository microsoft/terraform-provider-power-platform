# Error Handling on analyticsDataExport `nil` Return

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/datasource_analytics_data_exports.go

## Problem

When `analyticsDataExport` is `nil` after calling `GetAnalyticsDataExport`, a generic error is reported:

```
if analyticsDataExport == nil {
    resp.Diagnostics.AddError(
        "Analytics data export not found",
        "Unable to find analytics data export with the specified ID",
    )
    return
}
```

There is a lack of context to help diagnose why `nil` was returned (i.e. no ID in request, filtering logic, or an API-level explanation), and it's not clear if returning `nil` means a logical 'not found' or an internal error occurred. The error message gives a potentially misleading impression there is an ID-based lookup (none is present in code shown; possibly a copy-paste mistake in the message).

## Impact

This could confuse users and maintainers, making debugging more difficult and user feedback inadequate.

**Severity:** High.

## Location

```go
if analyticsDataExport == nil {
    resp.Diagnostics.AddError(
        "Analytics data export not found",
        "Unable to find analytics data export with the specified ID",
    )
    return
}
```

## Code Issue

```go
if analyticsDataExport == nil {
    resp.Diagnostics.AddError(
        "Analytics data export not found",
        "Unable to find analytics data export with the specified ID",
    )
    return
}
```

## Fix

Return a more helpful error explaining the likely reasons and removing reference to an ID (if not present). If another layer returns context, surface it here. If `nil` is legal and indicates a non-error empty state, consider returning early or sending a warning instead.

```go
if analyticsDataExport == nil {
    resp.Diagnostics.AddError(
        "Analytics data export not found",
        "No analytics data exports were returned from the downstream API. This may indicate a configuration issue or that no exports exist for this tenant.",
    )
    return
}
```
