# Title

Function `newLocationsClient` does not follow Go's naming conventions

##

/workspaces/terraform-provider-power-platform/internal/services/locations/api_locations.go

## Problem

Go conventionally uses camelCase for unexported functions. However, the current name, while technically correct, could be confused with a constructor for a type named `LocationsClient`, which does not exist. It would be better named as `newClient`, reflecting the actual type constructed.

## Impact

Low severity. This can cause minor confusion during code navigation and reduce readability and maintainability.

## Location

Top of the file, function definition:

## Code Issue

```go
func newLocationsClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}
```

## Fix

Rename the function to `newClient` to clearly associate with the struct being constructed.

```go
func newClient(apiClient *api.Client) client {
	return client{
		Api: apiClient,
	}
}
```
