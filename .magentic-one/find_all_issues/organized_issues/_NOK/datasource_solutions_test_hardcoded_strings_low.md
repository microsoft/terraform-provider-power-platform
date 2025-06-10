# Title

Hardcoded Strings for URLs and File Paths in Test Setup

##

/workspaces/terraform-provider-power-platform/internal/services/solution/datasource_solutions_test.go

## Problem

Throughout the test setup, service endpoints and mock data file paths are repeated as hardcoded strings. Changes in endpoint paths or data structure necessitate modifying several locations in the tests, which is error-prone and not maintainable should the endpoints or paths change.

## Impact

Less maintainable and harder to refactor if URLs or file paths change. This increases risk for inconsistencies. Severity: low.

## Location

All over each test function:

```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?api-version=2023-06-01`, ...)
```

## Code Issue

```go
httpmock.RegisterResponder("GET", `https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?api-version=2023-06-01`, ...)
...
httpmock.RegisterResponder("GET", `https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=%28isvisible+eq+true%29&%24orderby=createdon+desc`, ...)
...
httpmock.File("tests/datasource/Validate_Read/get_environments.json")
```

## Fix

At the top of the test file, define constants:

```go
const (
	envsURL = "https://api.bap.microsoft.com/providers/Microsoft.BusinessAppPlatform/scopes/admin/environments?api-version=2023-06-01"
	solutionsURL = "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/solutions?%24expand=publisherid&%24filter=%28isvisible+eq+true%29&%24orderby=createdon+desc"
	// ...other URLs and file paths as needed
	getEnvironmentsFile = "tests/datasource/Validate_Read/get_environments.json"
)
```

Then update your setup code:

```go
httpmock.RegisterResponder("GET", envsURL, ...)
...
httpmock.File(getEnvironmentsFile)
```
