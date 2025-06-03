# Title

Missing Test Coverage for Edge Cases and Error Handling

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/api_analytics_data_exports.go

## Problem

There are no unit tests or test files present for the APIs and utility functions in this file. Especially for areas dealing with varied error handling, region lookups, and building URLs, there should be explicit test coverage to verify both expected and edge-case behaviors.

## Impact

Lack of automated testing increases the risk of regressions and undetected bugs, particularly as the codebase scales or is modified by others. This impacts long-term maintainability and code reliability. Severity: Medium.

## Location

All exported functions and important helpers, specifically:
- `GetGatewayCluster`
- `GetAnalyticsDataExport`
- `getAnalyticsUrl`

## Code Issue

_No specific code snippet as this is a test/QA gap for the whole file._

## Fix

Add a corresponding `_test.go` file with table-driven tests for:
- Successful invocation of each exported function.
- Error path testing (e.g., region missing, tenant client/API call errors).
- Validation that error wrapping provides useful context.

Example for `getAnalyticsUrl`:

```go
func TestGetAnalyticsUrl(t *testing.T) {
	tests := []struct {
		region string
		want   string
		wantErr bool
	}{
		{"US", "https://na.csanalytics.powerplatform.microsoft.com/", false},
		{"INVALID", "", true},
	}

	for _, tt := range tests {
		got, err := getAnalyticsUrl(tt.region)
		if (err != nil) != tt.wantErr {
			t.Errorf("getAnalyticsUrl(%q) error = %v, wantErr %v", tt.region, err, tt.wantErr)
			continue
		}
		if got != tt.want {
			t.Errorf("getAnalyticsUrl(%q) = %v, want %v", tt.region, got, tt.want)
		}
	}
}
```

---

This file will be saved to:

```
/workspaces/terraform-provider-power-platform/.magentic-one/find_all_issues/issues_found/testing/api_analytics_data_exports_testing_medium.md
```
