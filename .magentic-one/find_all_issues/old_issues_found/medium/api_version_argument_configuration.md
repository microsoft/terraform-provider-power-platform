# Title

Hardcoded Values for API Version

##

/workspaces/terraform-provider-power-platform/internal/services/environment_templates/api_environment_templates.go

## Problem

The API version is hardcoded (`"2023-06-01"`) in the code, making future updates and modifications more error-prone. Any change in the API version would require modifying the source code repeatedly, which violates best practices for maintainability and scalability.

## Impact

1. Restricts flexibility to change API versions dynamically.
2. Increases maintenance effort when the API version changes.
3. Introduces potential risk of inconsistent handling across different parts of the codebase where the value may be hardcoded.

Severity: Medium

## Location

Within the `GetEnvironmentTemplatesByLocation` function:

```go
	apiUrl.RawQuery = url.Values{
		"api-version": []string{"2023-06-01"},
	}.Encode()
```

## Code Issue

Location of the hardcoded value:

```go
	apiUrl.RawQuery = url.Values{
		"api-version": []string{"2023-06-01"},
	}.Encode()
```

## Fix

Use a constant or configuration value for the API version to improve maintainability and scalability:

```go
	// Define a global or configurable constant for API version
	const defaultApiVersion = "2023-06-01"

	apiUrl.RawQuery = url.Values{
		"api-version": []string{defaultApiVersion},
	}.Encode()
```

Explanation: By utilizing a constant, you can ensure the API version is consistent across the codebase and easily updated in one place when needed.