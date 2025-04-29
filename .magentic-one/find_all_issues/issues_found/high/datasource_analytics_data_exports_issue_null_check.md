# Title

Null Check for `analyticsDataExport` Leaves Other Potential Context Unaddressed

##

`/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/datasource_analytics_data_exports.go`

## Problem

Within the `Read` method, the code validates if the `analyticsDataExport` variable is `nil`. If it is, it throws a diagnostic error indicating that such analytics data exports are not found. However, this fails to consider scenarios beyond just the `nil` condition—such as cases where the response structure from the API is invalid or incomplete.

## Impact

- **Severity:** High  
- Failure to account for alternative failure scenarios could lead to misleading error messages. Users may think the problem is that the data export is "not found," whereas the actual issue might be an API malfunction or data corruption.
- Detecting potential alternative causes would simplify debugging, ensuring the provider remains dependable for users.

## Location

`AnalyticsExportDataSource.Read`

## Code Issue

```go
	analyticsDataExport, err := d.analyticsExportClient.GetAnalyticsDataExport(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching analytics data export",
			fmt.Sprintf("Unable to fetch analytics data export: %s", err),
		)
		return
	}
	if analyticsDataExport == nil {
		resp.Diagnostics.AddError(
			"Analytics data export not found",
			"Unable to find analytics data export with the specified ID",
		)
		return
	}
```

## Fix

Instead of narrowly checking only for `nil`, expand validations to confirm the structural integrity of the `analyticsDataExport` response:

```go
	analyticsDataExport, err := d.analyticsExportClient.GetAnalyticsDataExport(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching analytics data export",
			fmt.Sprintf("Unable to fetch analytics data export: %s", err),
		)
		return
	}

	if analyticsDataExport == nil || len(*analyticsDataExport) == 0 {
		resp.Diagnostics.AddError(
			"Analytics data export not found",
			fmt.Sprintf("Analytics data export with the specified ID returned no results. Please check for existence or API response issues."),
		)
		return
	}

	// Add checks for malformed data exports
	for _, export := range *analyticsDataExport {
		if export.SomeRequiredField == "" {
			resp.Diagnostics.AddError(
				"Analytics data export contains invalid entries",
				"One or more entries in the analytics data export response are incomplete or invalid. Verify the API structure or its configuration.",
			)
			return
		}
	}
```