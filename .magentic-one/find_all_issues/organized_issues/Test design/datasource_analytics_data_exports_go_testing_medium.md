# Testing/Quality Assurance - No Evidence of Method-Level Unit Tests

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/datasource_analytics_data_exports.go

## Problem

There is no evidence from this file or from comments that unit tests exist *for the methods or exported logic* in this data source implementation. The public API (`Read`, `Configure`, `Schema`, `Metadata`) has logic that should be tested, especially edge and error cases, but there are no corresponding test files or even a `//go:generate mockgen` directive or similar commenting.

## Impact

- Reduces confidence in changes and refactoring
- Increases risk of regressions (especially in error or edge case handling)
- Slows down onboarding and review for new contributors
- Negative impact on code quality over time

**Severity:** Medium

## Location

All exported methods in the file.

## Code Issue

Absence of testing, not direct code.

## Fix

Add a test file, e.g., `datasource_analytics_data_exports_test.go`, and use Go's `testing` package with suitable mocking for the analytics client.

```go
func TestAnalyticsExportDataSource_Read_errorHandling(t *testing.T) {
    // Use a mock AnalyticsExportClient to inject errors, nils, and valid data for coverage
}
```

Add tests for edge cases, error surfaces, and happy path.
