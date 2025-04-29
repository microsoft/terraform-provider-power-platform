# Title

Invalid HTTP Mock Responder Responses

##

/workspaces/terraform-provider-power-platform/internal/services/authorization/datasource_securityroles_test.go

## Problem

The HTTP mock responders return responses that include hardcoded JSON file paths.

## Impact

Mock responses could lead to fragile tests that fail when file paths change. This violates maintainability and introduces external dependencies in the unit test. Severity: High.

## Location

TestUnitSecurityDataSource_Validate_Read

## Code Issue

```go
httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/roles",
	func(req *http.Request) (*http.Response, error) {
		return httpmock.NewStringResponse(http.StatusOK, httpmock.File("tests/datasource/security_roles/Validate_Read/get_security_roles.json").String()), nil
	})
```

## Fix

Instead of relying on external files, inline JSON strings can be defined directly within the mock responder, or mock files can be centralized to ensure robustness.

```go
httpmock.RegisterResponder("GET", "https://00000000-0000-0000-0000-000000000001.crm4.dynamics.com/api/data/v9.2/roles",
	func(req *http.Request) (*http.Response, error) {
		mockResponse := `{ "name": "Test Role", "role_id": "00000000-0000-0000-0000-000000000001" }
`
		return httpmock.NewStringResponse(http.StatusOK, mockResponse), nil
	})
```