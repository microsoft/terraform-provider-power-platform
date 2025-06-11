# No Nil Check for LocationsClient Before Method Call

##

/workspaces/terraform-provider-power-platform/internal/services/locations/datasource_locations.go

## Problem

In the `Read` method, code calls `d.LocationsClient.GetLocations(ctx)` without ever verifying that `d.LocationsClient` is non-nil. If the `Configure` method was not called or failed to set up the client (or if it was set to nil by some external event or mistake), calling a method on a nil struct would result in a runtime panic.

## Impact

This oversight could cause panics during provider execution, leading to plugin crashes and poor user/developer experience. This is a **High** severity control flow and error handling issue.

## Location

In the `Read` method:
```go
locations, err := d.LocationsClient.GetLocations(ctx)
```

## Code Issue

```go
locations, err := d.LocationsClient.GetLocations(ctx)
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
	return
}
```

## Fix

Check `d.LocationsClient` for `nil` before usage, and add a meaningful diagnostic if it is not configured:

```go
if d.LocationsClient == nil {
	resp.Diagnostics.AddError(
		"Locations client not configured",
		"The locations client was not configured. This may indicate a problem with the provider setup or authentication.",
	)
	return
}

locations, err := d.LocationsClient.GetLocations(ctx)
if err != nil {
	resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s", d.FullTypeName()), err.Error())
	return
}
```
