# Naming: Function Type Not Idiomatic

##

/workspaces/terraform-provider-power-platform/internal/services/languages/api_languages.go

## Problem

The function `newLanguagesClient` should be named `NewLanguagesClient` if it is intended to be part of the package’s public API.

## Impact

It makes the API confusing if an intended public API is unexported. Severity: **low**.

## Location

```go
func newLanguagesClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}
```

## Code Issue

```go
func newLanguagesClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}
```

## Fix

Rename the function for export:

```go
func NewLanguagesClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}
```
