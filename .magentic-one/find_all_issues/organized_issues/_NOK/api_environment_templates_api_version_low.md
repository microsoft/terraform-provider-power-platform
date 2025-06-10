# API Client: Hardcoded API Version String

##

/workspaces/terraform-provider-power-platform/internal/services/environment_templates/api_environment_templates.go

## Problem

The API version is hardcoded as `"2023-06-01"` in the `GetEnvironmentTemplatesByLocation` function. This makes upgrading/bumping API versions error-prone and inconsistent across the codebase.

## Impact

Hardcoding API versions can lead to inconsistencies and maintenance overhead if API versions change. Severity: **Low**.

## Location

```go
	apiUrl.RawQuery = url.Values{
		"api-version": []string{"2023-06-01"},
	}.Encode()
```

## Code Issue

```go
	apiUrl.RawQuery = url.Values{
		"api-version": []string{"2023-06-01"},
	}.Encode()
```

## Fix

Define the API version as a constant in a shared place, or at least at the top of the file and reuse it.

```go
const apiVersion = "2023-06-01"

// ...

	apiUrl.RawQuery = url.Values{
		"api-version": []string{apiVersion},
	}.Encode()
```
