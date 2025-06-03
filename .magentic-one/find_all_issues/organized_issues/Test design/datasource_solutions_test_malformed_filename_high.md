# Title

Malformed Filename in httpmock.File Path in TestUnitSolutionsDataSource_Validate_No_Dataverse

##

/workspaces/terraform-provider-power-platform/internal/services/solution/datasource_solutions_test.go

## Problem

In the `TestUnitSolutionsDataSource_Validate_No_Dataverse` test, the JSON file path `"tests/datasource/Validate_No_Dataverse(/get_environments.json"` is missing a closing parenthesis and should instead be `"tests/datasource/Validate_No_Dataverse/get_environments.json"`. This typo in the file path could lead to test failures due to the file not being found.

## Impact

This will cause the test to fail because the file path is incorrect, resulting in an inability to load the required mock data. Severity: high.

## Location

Line inside the function `TestUnitSolutionsDataSource_Validate_No_Dataverse`:

```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?api-version=2023-06-01`,
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_No_Dataverse(/get_environments.json").String()), nil
	})
```

## Code Issue

```go
return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_No_Dataverse(/get_environments.json").String()), nil
```

## Fix

Replace the incorrect file path with the correct one:

```go
return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/Validate_No_Dataverse/get_environments.json").String()), nil
```

