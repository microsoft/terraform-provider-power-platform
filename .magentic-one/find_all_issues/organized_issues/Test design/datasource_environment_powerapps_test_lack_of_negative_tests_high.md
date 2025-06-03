# Lack of Negative or Error Case Testing

##

/workspaces/terraform-provider-power-platform/internal/services/powerapps/datasource_environment_powerapps_test.go

## Problem

The current tests cover only successful or "happy path" responses but there is no validation of expected error handling, such as API failures, empty/malformed responses, or authorization errors. Comprehensive testing should also include negative scenarios to ensure robust error handling and clear error messages for users.

## Impact

- **Test Coverage**: High – Omitting error paths introduces risk for undetected bugs in error handling, potentially leading to runtime failures that are hard to debug.
- **Resilience**: Application quality may suffer if negative/test failures aren't considered.

## Location

File-wide: No negative test cases covering API errors or invalid input.

## Code Issue

No code for negative/error test cases (not present in the file).

## Fix

Add explicit test steps or functions to cover negative scenarios, e.g., simulate API failures:

```go
func TestUnitEnvironmentPowerAppsDataSource_ErrorHandling(t *testing.T) {
    httpmock.Activate()
    defer httpmock.DeactivateAndReset()

    // Simulate an API error/response
    httpmock.RegisterResponder("GET", "<URL>",
        httpmock.NewStringResponder(http.StatusInternalServerError, "Internal Server Error"))

    resource.Test(t, resource.TestCase{
        IsUnitTest: true,
        ProtoV6ProviderFactories: mocks.TestUnitTestProtoV6ProviderFactories,
        Steps: []resource.TestStep{
            {
                Config: `<valid_CONFIG>`,
                ExpectError: regexp.MustCompile("Internal Server Error"),
            },
        },
    })
}
```
