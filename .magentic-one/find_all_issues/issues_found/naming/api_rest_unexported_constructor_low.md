# Inconsistent visibility for constructor and type

##

/workspaces/terraform-provider-power-platform/internal/services/rest/api_rest.go

## Problem

The constructor `newWebApiClient` is unexported (lowercase), yet if the `Client` type is exported (see previous naming suggestion), the idiomatic Go pattern is to provide an exported constructor. This consistency improves usability for packages importing this API.

## Impact

Severity: Low. This affects code ergonomics and external usability, but not runtime functionality.

## Location

At the top of the file:

## Code Issue

```go
func newWebApiClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}
```

## Fix

Export the constructor and ensure the returned type is also exported:

```go
func NewWebApiClient(apiClient *api.Client) Client {
	return Client{
		Api: apiClient,
	}
}
```
