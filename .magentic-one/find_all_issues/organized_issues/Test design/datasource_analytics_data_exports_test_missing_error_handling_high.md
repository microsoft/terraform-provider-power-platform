# Lack of Test Coverage for Error Handling and Negative Scenarios

##

/workspaces/terraform-provider-power-platform/internal/services/analytics_data_export/datasource_analytics_data_exports_test.go

## Problem

The test file only validates positive cases (happy paths) and the ability to read/export analytics data. There are no tests that verify the system's behavior upon encountering errors, such as API failures, invalid responses, empty data, permission errors, or malformatted payloads. Good testing practice should include validation for expected errors and edge cases.

## Impact

High. Lack of coverage for error handling could result in undetected bugs making it to production, especially when the real API behaves unexpectedly.

## Location

```go
// No negative test steps or error response registration
```

## Code Issue

```go
	// Register responder for analytics data exports API
	httpmock.RegisterResponder(
		"GET",
		"https://na.csanalytics.powerplatform.microsoft.com/api/v2/connections",
		httpmock.NewStringResponder(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_analytics_data_exports.json").String()))

	resource.UnitTest(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				data "powerplatform_analytics_data_exports" "test" {}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
                    // only success checks, no checks for errors
```

## Fix

Add test cases to simulate error conditions, server failures, or malformed data. Use httpmock to return error HTTP codes or invalid payloads, and verify that the underlying system correctly handles these scenarios—i.e., raises errors, fails with expected error messages, or recovers gracefully.

```go
// Example error responder registration:
httpmock.RegisterResponder(
    "GET",
    "https://na.csanalytics.powerplatform.microsoft.com/api/v2/connections",
    httpmock.NewStringResponder(http.StatusInternalServerError, `{"error":"Internal server error"}`),
)

// Add a new test step to validate error handling
{
    Config: `
        data "powerplatform_analytics_data_exports" "test" {}
    `,
    ExpectError: regexp.MustCompile("Internal server error"), // or the expected error message from the provider
},
```
