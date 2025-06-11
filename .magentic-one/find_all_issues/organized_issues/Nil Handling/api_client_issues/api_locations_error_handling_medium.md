# Title

Missing nil check on parameter `apiClient` in `newLocationsClient`

##

/workspaces/terraform-provider-power-platform/internal/services/locations/api_locations.go

## Problem

The `newLocationsClient` function does not check if its parameter `apiClient` is nil. This could lead to runtime panics when the returned `client` is used.

## Impact

Medium severity. Could lead to runtime panics if an invalid `apiClient` pointer is passed.

## Location

`newLocationsClient` function

## Code Issue

```go
func newLocationsClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}
```

## Fix

Add a check and optionally return an error or panic with a clear message.

```go
func newClient(apiClient *api.Client) (client, error) {
	if apiClient == nil {
		return client{}, fmt.Errorf("apiClient cannot be nil")
	}
	return client{
		API: apiClient,
	}, nil
}
```
