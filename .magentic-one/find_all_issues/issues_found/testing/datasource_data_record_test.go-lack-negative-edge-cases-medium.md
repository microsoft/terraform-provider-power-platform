# Title

Lack of Negative and Edge Case Test Coverage

##

internal/services/data_record/datasource_data_record_test.go

## Problem

All current test cases focus on positive/expected usage scenarios (happy paths). There are no tests that:
- Check how the provider responds to invalid configuration (e.g., missing required fields, invalid values)
- Simulate API failures or partial/unexpected responses from external services
- Test edge cases such as empty results, invalid expand/nesting, or very large datasets

## Impact

Gaps in test coverage reduce confidence that the provider handles errors correctly and will behave robustly in production. It may allow subtle bugs to slip through to customers, especially around error handling, unexpected server responses, or fully/partially missing data. The severity is medium to high for projects that rely on stability in production environments.

## Location

Absence of negative/edge test cases throughout the test file. All tests are structured to simulate only valid Terraform config and valid/empty API results.

## Code Issue

```go
// There are no test cases like:
- TestAccDataRecordDatasource_InvalidConfig
- TestAccDataRecordDatasource_ApiError
- TestAccDataRecordDatasource_EmptyExpandArray
// etc.
```

## Fix

Add negative and edge case test functions, such as:

```go
func TestAccDataRecordDatasource_InvalidConfig(t *testing.T) {
    resource.Test(t, resource.TestCase{
        ProtoV6ProviderFactories: mocks.TestAccProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: `data "powerplatform_data_records" "invalid" { }`, // minimal/broken config
                ExpectError: regexp.MustCompile("required argument ... is missing"),
            },
        },
    })
}

// Similarly, simulate API returning an error
func TestUnitDataRecordDatasource_ApiError(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    httpmock.RegisterResponder("GET", "...", func(req *http.Request) (*http.Response, error) {
        return httpmock.NewStringResponse(http.StatusInternalServerError, `{"error": ...}`), nil
    })
    // ...
    // Expect TestError or fatal, etc.
}
```
Adding such tests will improve robustness and ensure that error handling and validation are not inadvertently broken by future code changes.

Save as a testing and quality assurance issue.
