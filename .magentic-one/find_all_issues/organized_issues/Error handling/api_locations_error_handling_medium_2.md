# Title

Does not provide context when returning errors from `GetLocations`

##

/workspaces/terraform-provider-power-platform/internal/services/locations/api_locations.go

## Problem

The `GetLocations` method returns errors from the API client directly, without wrapping or enriching them with context about the operation that failed.

## Impact

Medium severity. Debugging can be more difficult if consumers of this function cannot disambiguate where an error occurred.

## Location

`GetLocations` method return statement.

## Code Issue

```go
_, err := client.Api.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &locations)
return locations, err
```

## Fix

Wrap the error to provide more context.

```go
_, err := client.API.Execute(ctx, nil, "GET", apiUrl.String(), nil, nil, []int{http.StatusOK}, &locations)
if err != nil {
	return locations, fmt.Errorf("failed to get locations: %w", err)
}
return locations, nil
```
