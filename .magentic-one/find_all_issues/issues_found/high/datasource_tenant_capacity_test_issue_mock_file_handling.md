# Title

Improper Mock File Handling in HTTP Responder

##

/workspaces/terraform-provider-power-platform/internal/services/capacity/datasource_tenant_capacity_test.go

## Problem

The `httpmock.RegisterResponder` function uses `httpmock.File` to retrieve a JSON file for constructing a mock HTTP response. However, there is no explicit error handling for cases where the file might not exist or cannot be read. This can lead to runtime panics and unreliable test execution.

## Impact

If the `tests/datasource/Validate_Read/get_tenant_capacity.json` file is missing, inaccessible, or malformed, the test will crash unexpectedly during execution, impacting developer productivity and the reliability of test cases. The severity of this issue is **high**, considering its criticality during test automation.

## Location

Within the function `TestUnitTenantCapacityDataSource_Validate_Read`:

```go
httpmock.RegisterResponder(
	"GET",
	`https://licensing.powerplatform.microsoft.com/v0.1-alpha/tenants/00000000-0000-0000-0000-000000000001/TenantCapacity`,
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_Read/get_tenant_capacity.json").String()), nil
	});
```

## Code Issue

The code block does not include checks for file existence or readability before using the `httpmock.File` API:

```go
httpmock.File("tests/datasource/Validate_Read/get_tenant_capacity.json").String()
```

## Fix

Introduce a check to validate the existence and readability of the JSON file before using it to generate a response. Additionally, log meaningful error messages for debugging purposes.

```go
func(req *http.Request) (*http.Response, error) {
	filePath := "tests/datasource/Validate_Read/get_tenant_capacity.json"
	mockFileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		// Log a meaningful error or fail gracefully
		fmt.Printf("Failed to load mock file %s: %v\n", filePath, err)
		return httpmock.NewStringResponse(http.StatusInternalServerError, ""), nil
	}
	return httpmock.NewStringResponse(http.StatusOK, string(mockFileContent)), nil
});
```

- This fix ensures that tests fail gracefully if the mock file is missing or inaccessible, rather than panicking.
- The severity remains high due to its potential impact on test reliability.
