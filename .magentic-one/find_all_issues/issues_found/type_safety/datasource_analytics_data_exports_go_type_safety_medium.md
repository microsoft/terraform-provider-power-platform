# Type Safety: Lack of Validation on Downstream API Data

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/datasource_analytics_data_exports.go

## Problem

After calling the downstream API to retrieve analytics data exports, there is no validation of the content, shape, or completeness of the data before mapping it to the model. If the downstream API changes (e.g., omits required fields, returns nulls/unexpected values), this could lead to panics, zero-values, or subtle bugs communicated to end-users.

## Impact

- Could cause runtime panics if downstream returns malformed data
- Silent loss of information if fields become missing
- Reduces robustness, and user-facing errors become harder to diagnose

**Severity:** Medium

## Location

```go
analyticsDataExport, err := d.analyticsExportClient.GetAnalyticsDataExport(ctx)
// ...
for _, export := range *analyticsDataExport {
    if model := convertDtoToModel(&export); model != nil {
        exports = append(exports, *model)
    }
}
```

## Code Issue

No explicit type or value checks after obtaining data from downstream API.

## Fix

Add validation or sanity checks before mapping data (inside or prior to `convertDtoToModel`). Example:

```go
for _, export := range *analyticsDataExport {
    // Validate required fields are present and sane
    if export.ID == "" || export.Sink == nil {
        resp.Diagnostics.AddWarning("Incomplete data", fmt.Sprintf("Found analytics data export with nil/empty ID or Sink: %+v", export))
        continue
    }
    if model := convertDtoToModel(&export); model != nil {
        exports = append(exports, *model)
    }
}
```

Consider making `convertDtoToModel` return errors for invalid or inconsistent structures, or introducing a validation helper inline.
